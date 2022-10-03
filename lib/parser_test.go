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
	"reflect"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	data := Tabdata{
		maxwidthHeader: 5,
		maxwidthPerCol: []int{
			5,
			5,
			8,
		},
		columns: 3,
		headerIndices: []map[string]int{
			map[string]int{
				"beg": 0,
				"end": 6,
			},
			map[string]int{
				"end": 13,
				"beg": 7,
			},
			map[string]int{
				"beg": 14,
				"end": 0,
			},
		},
		headers: []string{
			"ONE",
			"TWO",
			"THREE",
		},
		entries: [][]string{
			[]string{
				"asd",
				"igig",
				"cxxxncnc",
			},
			[]string{
				"19191",
				"EDD 1",
				"X",
			},
		},
	}

	table := `ONE    TWO    THREE  
asd    igig   cxxxncnc  
19191  EDD 1  X`

	readFd := strings.NewReader(table)
	gotdata, err := parseFile(readFd, "")

	if err != nil {
		t.Errorf("Parser returned error: %s\nData processed so far: %+v", err, gotdata)
	}

	if !reflect.DeepEqual(data, gotdata) {
		t.Errorf("Parser returned invalid data\nExp: %+v\nGot: %+v\n", data, gotdata)
	}
}
