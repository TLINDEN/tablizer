/*
Copyright © 2023 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package lib

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/glycerine/zygomys/zygo"
	"github.com/tlinden/tablizer/cfg"
)

/*
needs to be global because we can't feed an cfg object to AddHook()
which is being called from user lisp code
*/
var Hooks map[string][]*zygo.SexpSymbol

/*
AddHook() (called addhook from lisp code)  can be used by the user to
add a function to one of the available hooks provided by tablizer.
*/
func AddHook(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	var hookname string

	if len(args) < 2 {
		return zygo.SexpNull, errors.New("argument of %add-hook should be: %hook-name %your-function")
	}

	switch sexptype := args[0].(type) {
	case *zygo.SexpSymbol:
		if !HookExists(sexptype.Name()) {
			return zygo.SexpNull, errors.New("Unknown hook " + sexptype.Name())
		}

		hookname = sexptype.Name()

	default:
		return zygo.SexpNull, errors.New("hook name must be a symbol ")
	}

	switch sexptype := args[1].(type) {
	case *zygo.SexpSymbol:
		_, exists := Hooks[hookname]
		if !exists {
			Hooks[hookname] = []*zygo.SexpSymbol{sexptype}
		} else {
			Hooks[hookname] = append(Hooks[hookname], sexptype)
		}

	default:
		return zygo.SexpNull, errors.New("hook function must be a symbol ")
	}

	return zygo.SexpNull, nil
}

/*
Check if a hook exists
*/
func HookExists(key string) bool {
	for _, hook := range cfg.ValidHooks {
		if hook == key {
			return true
		}
	}

	return false
}

/*
 * Basic sanity checks and load lisp file
 */
func LoadAndEvalFile(env *zygo.Zlisp, path string) error {
	if strings.HasSuffix(path, `.zy`) {
		code, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read lisp file %s: %w", path, err)
		}

		// FIXME: check what res (_ here) could be and mean
		_, err = env.EvalString(string(code))
		if err != nil {
			log.Fatal(env.GetStackTrace(err))
		}
	}

	return nil
}

/*
 * Setup lisp interpreter environment
 */
func SetupLisp(conf *cfg.Config) error {
	// iterate over load-path and evaluate all *.zy files there, if any
	// we ignore if load-path does not exist, which is the default anyway
	path, err := os.Stat(conf.LispLoadPath)

	if err != nil {
		if os.IsNotExist(err) {
			// ignore non-existent files
			return nil
		}

		return fmt.Errorf("failed to stat path: %w", err)
	}

	// init global hooks
	Hooks = make(map[string][]*zygo.SexpSymbol)

	// init sandbox
	env := zygo.NewZlispSandbox()
	env.AddFunction("addhook", AddHook)

	if !path.IsDir() {
		// load single lisp file
		err = LoadAndEvalFile(env, conf.LispLoadPath)
		if err != nil {
			return err
		}
	} else {
		// load all lisp file in load dir
		dir, err := os.ReadDir(conf.LispLoadPath)
		if err != nil {
			return fmt.Errorf("failed to read lisp dir %s: %w",
				conf.LispLoadPath, err)
		}

		for _, entry := range dir {
			if !entry.IsDir() {
				err := LoadAndEvalFile(env, conf.LispLoadPath+"/"+entry.Name())
				if err != nil {
					return err
				}
			}
		}
	}

	RegisterLib(env)

	conf.Lisp = env

	return nil
}

/*
Execute every user lisp function registered as filter hook.

Each function is given the current line as argument and is expected to
return a boolean.   True indicates to keep the line,  false to skip
it.

If there  are multiple such  functions registered, then the  first one
returning false wins,  that is if each function returns  true the line
will  be kept,  if at  least one  of them  returns false,  it will  be
skipped.
*/
func RunFilterHooks(conf cfg.Config, line string) (bool, error) {
	for _, hook := range Hooks["filter"] {
		var result bool

		conf.Lisp.Clear()

		res, err := conf.Lisp.EvalString(fmt.Sprintf("(%s `%s`)", hook.Name(), line))
		if err != nil {
			return false, fmt.Errorf("failed to evaluate hook loader: %w", err)
		}

		switch sexptype := res.(type) {
		case *zygo.SexpBool:
			result = sexptype.Val
		default:
			return false, fmt.Errorf("filter hook shall return bool")
		}

		if !result {
			// the first hook which returns false leads to complete false
			return result, nil
		}
	}

	// if no hook returned false, we succeed and accept the given line
	return true, nil
}

/*
These hooks get the data (Tabdata) readily processed by tablizer as
argument. They are expected to  return a SexpPair containing a boolean
denoting  if  the data  has  been  modified  and the  actual  modified
data. Columns must  be the same, rows may differ.  Cells may also have
been modified.

Replaces the  internal data structure  Tabdata with the  user supplied
version.

Only one process hook function is supported.

The somewhat complicated code is being  caused by the fact, that we
need to convert our internal structure  to a lisp variable and vice
versa afterwards.
*/
func RunProcessHooks(conf cfg.Config, data Tabdata) (Tabdata, bool, error) {
	var userdata Tabdata

	lisplist := []zygo.Sexp{}

	if len(Hooks["process"]) == 0 {
		return userdata, false, nil
	}

	if len(Hooks["process"]) > 1 {
		fmt.Println("Warning: only one process hook is allowed!")
	}

	// there are hook[s] installed, convert the go data structure 'data to lisp
	for _, row := range data.entries {
		var entry zygo.SexpHash

		for idx, cell := range row {
			err := entry.HashSet(&zygo.SexpStr{S: data.headers[idx]}, &zygo.SexpStr{S: cell})
			if err != nil {
				return userdata, false, fmt.Errorf("failed to convert to lisp data: %w", err)
			}
		}

		lisplist = append(lisplist, &entry)
	}

	// we need to add it to the env so that the function can use the struct directly
	conf.Lisp.AddGlobal("data", &zygo.SexpArray{Val: lisplist, Env: conf.Lisp})

	// execute the actual hook
	hook := Hooks["process"][0]

	conf.Lisp.Clear()

	var result bool

	res, err := conf.Lisp.EvalString(fmt.Sprintf("(%s data)", hook.Name()))
	if err != nil {
		return userdata, false, fmt.Errorf("failed to eval lisp loader: %w", err)
	}

	// we expect (bool, array(hash)) as return from the function
	switch sexptype := res.(type) {
	case *zygo.SexpPair:
		switch th := sexptype.Head.(type) {
		case *zygo.SexpBool:
			result = th.Val
		default:
			return userdata, false, errors.New("xpect (bool, array(hash)) as return value")
		}

		switch sexptailtype := sexptype.Tail.(type) {
		case *zygo.SexpArray:
			lisplist = sexptailtype.Val
		default:
			return userdata, false, errors.New("expect (bool, array(hash)) as return value ")
		}
	default:
		return userdata, false, errors.New("filter hook shall return array of hashes ")
	}

	if !result {
		// no further processing required
		return userdata, result, nil
	}

	// finally convert lispdata back to Tabdata
	for _, item := range lisplist {
		row := []string{}

		switch hash := item.(type) {
		case *zygo.SexpHash:
			for _, header := range data.headers {
				entry, err := hash.HashGetDefault(
					conf.Lisp,
					&zygo.SexpStr{S: header},
					&zygo.SexpStr{S: ""})
				if err != nil {
					return userdata, false, fmt.Errorf("failed to get lisp hash entry: %w", err)
				}

				switch sexptype := entry.(type) {
				case *zygo.SexpStr:
					row = append(row, sexptype.S)
				default:
					return userdata, false, errors.New("hsh values should be string ")
				}
			}
		default:
			return userdata, false, errors.New("rturned array should contain hashes ")
		}

		userdata.entries = append(userdata.entries, row)
	}

	userdata.headers = data.headers

	return userdata, result, nil
}
