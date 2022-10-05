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
	"reflect"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	data := Tabdata{
		maxwidthHeader: 5,
		maxwidthPerCol: []int{
			5, 5, 8,
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

	table := `ONE    TWO    THREE  
asd    igig   cxxxncnc  
19191  EDD 1  X`

	readFd := strings.NewReader(table)
	gotdata, err := parseFile(readFd, "")
	Separator = DefaultSeparator

	if err != nil {
		t.Errorf("Parser returned error: %s\nData processed so far: %+v", err, gotdata)
	}

	if !reflect.DeepEqual(data, gotdata) {
		t.Errorf("Parser returned invalid data, Regex: %s\nExp: %+v\nGot: %+v\n", Separator, data, gotdata)
	}
}

func TestParserPatternmatching(t *testing.T) {
	var tests = []struct {
		entries [][]string
		pattern string
		invert  bool
	}{
		{
			entries: [][]string{
				[]string{
					"asd", "igig", "cxxxncnc",
				},
			},
			pattern: "ig",
			invert:  false,
		},
		{
			entries: [][]string{
				[]string{
					"19191", "EDD 1", "X",
				},
			},
			pattern: "ig",
			invert:  true,
		},
	}

	table := `ONE    TWO    THREE  
asd    igig   cxxxncnc  
19191  EDD 1  X`

	for _, tt := range tests {
		testname := fmt.Sprintf("parse-with-inverted-pattern-%t", tt.invert)
		t.Run(testname, func(t *testing.T) {
			InvertMatch = tt.invert

			readFd := strings.NewReader(table)
			gotdata, err := parseFile(readFd, tt.pattern)

			if err != nil {
				t.Errorf("Parser returned error: %s\nData processed so far: %+v", err, gotdata)
			}

			if !reflect.DeepEqual(tt.entries, gotdata.entries) {
				t.Errorf("Parser returned invalid data (pattern: %s, invert: %t)\nExp: %+v\nGot: %+v\n",
					tt.pattern, tt.invert, tt.entries, gotdata.entries)
			}
		})
	}
}
