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

func stdout2pipe(t *testing.T) (*os.File, *os.File) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	origStdout := os.Stdout
	os.Stdout = writer

	// we need to tell the color mode the io.Writer, even if we don't usw colorization
	color.SetOutput(writer)

	return origStdout, reader
}

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

		"extended": `ONE(1): asd
  TWO(2): igig
THREE(3): cxxxncnc

  ONE(1): 19191
  TWO(2): EDD 1
THREE(3): X`,
	}

	NoColor = true
	SortByColumn = 0 // disable sorting

	origStdout, reader := stdout2pipe(t)

	for mode, expect := range expects {
		testname := fmt.Sprintf("print-%s", mode)
		t.Run(testname, func(t *testing.T) {

			OutputMode = mode
			data := startdata // we need to reset our mock data, since it's being modified in printData()
			printData(&data)

			buf := make([]byte, 1024)
			n, err := reader.Read(buf)
			if err != nil {
				t.Fatal(err)
			}
			buf = buf[:n]
			output := strings.TrimSpace(string(buf))

			if output != expect {
				t.Errorf("output mode: %s, got:\n%s\nwant:\n%s\n (%d <=> %d)",
					mode, output, expect, len(output), len(expect))
			}
		})
	}

	// Restore
	os.Stdout = origStdout

}

func TestSortPrinter(t *testing.T) {
	startdata := Tabdata{
		maxwidthHeader: 5,
		maxwidthPerCol: []int{
			3,
			3,
			2,
		},
		columns: 3,
		headers: []string{
			"ONE", "TWO", "THREE",
		},
		entries: [][]string{
			[]string{
				"abc", "345", "b1",
			},
			[]string{
				"bcd", "234", "a2",
			},
			[]string{
				"cde", "123", "c3",
			},
		},
	}

	var tests = []struct {
		data   Tabdata
		sortby int
		expect string
	}{
		{
			data:   startdata,
			sortby: 1,
			expect: `ONE(1)	TWO(2)	THREE(3) 
abc   	345   	b1      	
bcd   	234   	a2      	
cde   	123   	c3`,
		},

		{
			data:   startdata,
			sortby: 2,
			expect: `ONE(1)	TWO(2)	THREE(3) 
cde   	123   	c3      	
bcd   	234   	a2      	
abc   	345   	b1`,
		},

		{
			data:   startdata,
			sortby: 3,
			expect: `ONE(1)	TWO(2)	THREE(3) 
bcd   	234   	a2      	
abc   	345   	b1      	
cde   	123   	c3`,
		},
	}

	NoColor = true
	OutputMode = "ascii"
	origStdout, reader := stdout2pipe(t)

	for _, tt := range tests {
		testname := fmt.Sprintf("print-sorted-table-%d", tt.sortby)
		t.Run(testname, func(t *testing.T) {
			SortByColumn = tt.sortby

			printData(&tt.data)

			buf := make([]byte, 1024)
			n, err := reader.Read(buf)
			if err != nil {
				t.Fatal(err)
			}
			buf = buf[:n]
			output := strings.TrimSpace(string(buf))

			if output != tt.expect {
				t.Errorf("sort column: %d, got:\n%s\nwant:\n%s",
					tt.sortby, output, tt.expect)
			}
		})
	}

	// Restore
	os.Stdout = origStdout
}
