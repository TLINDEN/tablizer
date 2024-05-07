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
	"regexp"
	"strconv"

	"github.com/glycerine/zygomys/zygo"
)

func Splice2SexpList(list []string) zygo.Sexp {
	slist := []zygo.Sexp{}

	for _, item := range list {
		slist = append(slist, &zygo.SexpStr{S: item})
	}

	return zygo.MakeList(slist)
}

func StringReSplit(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	if len(args) < 2 {
		return zygo.SexpNull, errors.New("expecting 2 arguments")
	}

	var separator, input string

	switch t := args[0].(type) {
	case *zygo.SexpStr:
		input = t.S
	default:
		return zygo.SexpNull, errors.New("second argument must be a string")
	}

	switch t := args[1].(type) {
	case *zygo.SexpStr:
		separator = t.S
	default:
		return zygo.SexpNull, errors.New("first argument must be a string")
	}

	sep := regexp.MustCompile(separator)

	return Splice2SexpList(sep.Split(input, -1)), nil
}

func String2Int(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	var number int

	switch t := args[0].(type) {
	case *zygo.SexpStr:
		num, err := strconv.Atoi(t.S)

		if err != nil {
			return zygo.SexpNull, fmt.Errorf("failed to convert string to number: %w", err)
		}

		number = num

	default:
		return zygo.SexpNull, errors.New("argument must be a string")
	}

	return &zygo.SexpInt{Val: int64(number)}, nil
}

func RegisterLib(env *zygo.Zlisp) {
	env.AddFunction("resplit", StringReSplit)
	env.AddFunction("atoi", String2Int)
}
