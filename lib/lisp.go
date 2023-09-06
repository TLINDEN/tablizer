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

	switch t := args[0].(type) {
	case *zygo.SexpSymbol:
		if !HookExists(t.Name()) {
			return zygo.SexpNull, errors.New("Unknown hook " + t.Name())
		}
		hookname = t.Name()
	default:
		return zygo.SexpNull, errors.New("hook name must be a symbol!")
	}

	switch t := args[1].(type) {
	case *zygo.SexpSymbol:
		_, exists := Hooks[hookname]
		if !exists {
			Hooks[hookname] = []*zygo.SexpSymbol{t}
		} else {
			Hooks[hookname] = append(Hooks[hookname], t)
		}
	default:
		return zygo.SexpNull, errors.New("hook function must be a symbol!")
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
func LoadFile(env *zygo.Zlisp, path string) error {
	if strings.HasSuffix(path, `.zy`) {
		code, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// FIXME: check what res (_ here) could be and mean
		_, err = env.EvalString(string(code))
		if err != nil {
			log.Fatalf(env.GetStackTrace(err))
		}
	}

	return nil
}

/*
 * Setup lisp interpreter environment
 */
func SetupLisp(c *cfg.Config) error {
	Hooks = make(map[string][]*zygo.SexpSymbol)

	env := zygo.NewZlispSandbox()
	env.AddFunction("addhook", AddHook)

	// iterate over load-path and evaluate all *.zy files there, if any
	// we ignore if load-path does not exist, which is the default anyway
	if path, err := os.Stat(c.LispLoadPath); !os.IsNotExist(err) {
		if !path.IsDir() {
			err := LoadFile(env, c.LispLoadPath)
			if err != nil {
				return err
			}
		} else {
			dir, err := os.ReadDir(c.LispLoadPath)
			if err != nil {
				return err
			}

			for _, entry := range dir {
				if !entry.IsDir() {
					LoadFile(env, c.LispLoadPath+"/"+entry.Name())
				}
			}
		}
	}

	RegisterLib(env)

	c.Lisp = env
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
func RunFilterHooks(c cfg.Config, line string) (bool, error) {
	for _, hook := range Hooks["filter"] {
		var result bool
		c.Lisp.Clear()
		res, err := c.Lisp.EvalString(fmt.Sprintf("(%s `%s`)", hook.Name(), line))
		if err != nil {
			return false, err
		}

		switch t := res.(type) {
		case *zygo.SexpBool:
			result = t.Val
		default:
			return false, errors.New("filter hook shall return BOOL!")
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

The somewhat complicated code is being  caused by the fact, that we
need to convert our internal structure  to a lisp variable and vice
versa afterwards.
*/
func RunProcessHooks(c cfg.Config, data Tabdata) (Tabdata, bool, error) {
	var userdata Tabdata
	lisplist := []zygo.Sexp{}

	if len(Hooks["process"]) > 0 {
		// there are hook[s] installed, convert the go data structure 'data to lisp
		for _, row := range data.entries {
			var entry zygo.SexpHash

			for idx, cell := range row {
				err := entry.HashSet(&zygo.SexpStr{S: data.headers[idx]}, &zygo.SexpStr{S: cell})
				if err != nil {
					return userdata, false, err
				}
			}

			lisplist = append(lisplist, &entry)
		}

		// we need to add it to the env so that the function can use the struct directly
		c.Lisp.AddGlobal("data", &zygo.SexpArray{Val: lisplist, Env: c.Lisp})

		// execute the actual hooks
		for _, hook := range Hooks["process"] {
			var result bool

			c.Lisp.Clear()

			res, err := c.Lisp.EvalString(fmt.Sprintf("(%s data)", hook.Name()))
			if err != nil {
				return userdata, false, err
			}

			// we expect (bool, array(hash)) as return from the function
			switch t := res.(type) {
			case *zygo.SexpPair:
				switch th := t.Head.(type) {
				case *zygo.SexpBool:
					result = th.Val
				default:
					return userdata, false, errors.New("Expect (bool, array(hash)) as return value!")
				}

				switch tt := t.Tail.(type) {
				case *zygo.SexpArray:
					lisplist = tt.Val
				default:
					return userdata, false, errors.New("Expect (bool, array(hash)) as return value!")
				}
			default:
				return userdata, false, errors.New("filter hook shall return array of hashes!")
			}

			if !result {
				// the first hook which returns false leads to complete false
				return userdata, result, nil
			}

			// finally convert lispdata back to Tabdata
			for _, item := range lisplist {
				row := []string{}

				switch hash := item.(type) {
				case *zygo.SexpHash:
					for _, header := range data.headers {
						entry, err := hash.HashGetDefault(c.Lisp, &zygo.SexpStr{S: header}, &zygo.SexpStr{S: ""})
						if err != nil {
							return userdata, false, err
						}

						switch t := entry.(type) {
						case *zygo.SexpStr:
							row = append(row, t.S)
						default:
							return userdata, false, errors.New("Hash values should be string!")
						}
					}
				default:
					return userdata, false, errors.New("Returned array should contain hashes!")
				}

				userdata.entries = append(userdata.entries, row)
			}

			userdata.headers = data.headers

			return userdata, result, nil
		}
	}

	return userdata, false, nil
}
