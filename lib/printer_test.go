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
	"fmt"
	"github.com/gookit/color"
	"os"
	"strings"
	"testing"
)

func TestPrinter(t *testing.T) {
	startdata := Tabdata{
		maxwidthHeader: 5,
		maxwidthPerCol: []int{
			5,
			5,
			8,
		},
		columns: 3,
		headers: []string{
			"ONE", "TWO", "THREE",
		},
		entries: [][]string{
			[]string{
				"asd", "igig", "cxxxncnc",
			},
			[]string{
				"19191", "EDD 1", "X",
			},
		},
	}

	expects := map[string]string{
		"ascii": `ONE(1)	TWO(2)	THREE(3) 
asd   	igig  	cxxxncnc	
19191 	EDD 1 	X`,
		"orgtbl": `|--------+--------+----------|
| ONE(1) | TWO(2) | THREE(3) |
|--------+--------+----------|
| asd    | igig   | cxxxncnc |
|  19191 | EDD 1  | X        |
|--------+--------+----------|`,
		"markdown": `| ONE(1) | TWO(2) | THREE(3) |
|--------|--------|----------|
| asd    | igig   | cxxxncnc |
|  19191 | EDD 1  | X        |`,
		"shell": `ONE="asd" TWO="igig" THREE="cxxxncnc"
ONE="19191" TWO="EDD 1" THREE="X"`,
	}

	NoColor = true

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	origStdout := os.Stdout
	os.Stdout = w

	// we need to tell the color mode the io.Writer, even if we don't usw colorization
	color.SetOutput(w)

	for mode, expect := range expects {
		testname := fmt.Sprintf("print-%s", mode)
		t.Run(testname, func(t *testing.T) {

			OutputMode = mode
			data := startdata // we need to reset our mock data, since it's being modified in printData()
			printData(&data)

			buf := make([]byte, 1024)
			n, err := r.Read(buf)
			if err != nil {
				t.Fatal(err)
			}
			buf = buf[:n]
			output := strings.TrimSpace(string(buf))

			if output != expect {
				t.Errorf("output mode: %s, got:\n%s\nwant:\n%s\n (%d <=> %d)", mode, output, expect, len(output), len(expect))
			}
		})
	}

	// Restore
	os.Stdout = origStdout

}
