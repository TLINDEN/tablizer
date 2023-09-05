/*
Copyright © 2022 Thomas von Dein

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
 * Setup lisp interpreter environment
 */
func SetupLisp(c *cfg.Config) error {
	Hooks = make(map[string][]*zygo.SexpSymbol)

	env := zygo.NewZlispSandbox()
	env.AddFunction("addhook", AddHook)

	// iterate over load-path and evaluate all *.ty files there, if any
	if _, err := os.Stat(c.LispLoadPath); !os.IsNotExist(err) {
		dir, err := os.ReadDir(c.LispLoadPath)
		if err != nil {
			return err
		}

		for _, entry := range dir {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), `.zy`) {
				code, err := os.ReadFile(c.LispLoadPath + "/" + entry.Name())
				if err != nil {
					return err
				}

				// FIXME: check what res (_ here) could be and mean
				_, err = env.EvalString(string(code))
				if err != nil {
					log.Fatalf(env.GetStackTrace(err))
				}
			}
		}
	}

	c.Lisp = env
	return nil
}

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
