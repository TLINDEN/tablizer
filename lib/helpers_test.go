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
	"testing"
)

func Testcontains(t *testing.T) {
	var tests = []struct {
		list   []int
		search int
		want   bool
	}{
		{[]int{1, 2, 3}, 2, true},
		{[]int{2, 3, 4}, 5, false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("contains-%d,%d,%t", tt.list, tt.search, tt.want)
		t.Run(testname, func(t *testing.T) {
			answer := contains(tt.list, tt.search)
			if answer != tt.want {
				t.Errorf("got %t, want %t", answer, tt.want)
			}
		})
	}
}

func TestPrepareColumns(t *testing.T) {
	data := Tabdata{
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
			{
				"2", "3", "4",
			},
		},
	}
	var tests = []struct {
		input     string
		exp       []int
		wanterror bool // expect error
	}{
		{"1,2,3", []int{1, 2, 3}, false},
		{"1,2,", []int{}, true},
		{"T", []int{2, 3}, false},
		{"T,2,3", []int{2, 3}, false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("PrepareColumns-%s-%t", tt.input, tt.wanterror)
		t.Run(testname, func(t *testing.T) {
			Columns = tt.input
			err := PrepareColumns(&data)
			if err != nil {
				if !tt.wanterror {
					t.Errorf("got error: %v", err)
				}
			} else {
				if !reflect.DeepEqual(UseColumns, tt.exp) {
					t.Errorf("got: %v, expected: %v", UseColumns, tt.exp)
				}
			}
		})
	}
}

func TestReduceColumns(t *testing.T) {
	var tests = []struct {
		expect  [][]string
		columns []int
	}{
		{
			expect:  [][]string{{"a", "b"}},
			columns: []int{1, 2},
		},
		{
			expect:  [][]string{{"a", "c"}},
			columns: []int{1, 3},
		},
		{
			expect:  [][]string{{"a"}},
			columns: []int{1},
		},
		{
			expect:  [][]string{nil},
			columns: []int{4},
		},
	}

	input := [][]string{{"a", "b", "c"}}

	Columns = "y" // used as a flag with len(Columns)...

	for _, tt := range tests {
		testname := fmt.Sprintf("reduce-columns-by-%+v", tt.columns)
		t.Run(testname, func(t *testing.T) {
			UseColumns = tt.columns
			data := Tabdata{entries: input}
			reduceColumns(&data)
			if !reflect.DeepEqual(data.entries, tt.expect) {
				t.Errorf("reduceColumns returned invalid data:\ngot: %+v\nexp: %+v", data.entries, tt.expect)
			}
		})
	}

	Columns = "" // reset for other tests
	UseColumns = nil
}
