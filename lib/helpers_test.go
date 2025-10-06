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
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tlinden/tablizer/cfg"
)

func TestContains(t *testing.T) {
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
			answer := slices.Contains(tt.list, tt.search)

			assert.EqualValues(t, tt.want, answer)
		})
	}
}

func TestPrepareColumns(t *testing.T) {
	data := Tabdata{
		maxwidthHeader: 5,
		columns:        3,
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
		{"T.", []int{2, 3}, false},
		{"T.,2,3", []int{2, 3}, false},
		{"[a-z,4,5", []int{4, 5}, true}, // invalid regexp
	}

	for _, testdata := range tests {
		testname := fmt.Sprintf("PrepareColumns-%s-%t",
			testdata.input, testdata.wanterror)
		t.Run(testname, func(t *testing.T) {
			conf := cfg.Config{Columns: testdata.input}
			err := PrepareColumns(&conf, &data)

			if testdata.wanterror {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, testdata.exp, conf.UseColumns)
			}
		})
	}
}

func TestPrepareTransposerColumns(t *testing.T) {
	data := Tabdata{
		maxwidthHeader: 5,
		columns:        3,
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
		transp    []string
		exp       int
		wanterror bool // expect error
	}{
		{
			"1",
			[]string{`/\d/x/`},
			1,
			false,
		},
		{
			"T.", // will match [T]WO and [T]HREE
			[]string{`/\d/x/`, `/.//`},
			2,
			false,
		},
		{
			"TH.,2",
			[]string{`/\d/x/`, `/.//`},
			2,
			false,
		},
		{
			"1",
			[]string{},
			1,
			true,
		},
		{
			"",
			[]string{`|.|N|`},
			0,
			true,
		},
		{
			"1",
			[]string{`|.|N|`},
			1,
			false,
		},
	}

	for _, testdata := range tests {
		testname := fmt.Sprintf("PrepareTransposerColumns-%s-%t", testdata.input, testdata.wanterror)
		t.Run(testname, func(t *testing.T) {
			conf := cfg.Config{TransposeColumns: testdata.input, Transposers: testdata.transp}
			err := PrepareTransposerColumns(&conf, &data)

			if testdata.wanterror {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, testdata.exp, len(conf.UseTransposeColumns))
				assert.EqualValues(t, len(conf.UseTransposeColumns), len(conf.Transposers))
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

	for _, testdata := range tests {
		testname := fmt.Sprintf("reduce-columns-by-%+v", testdata.columns)

		t.Run(testname, func(t *testing.T) {
			c := cfg.Config{Columns: "x", UseColumns: testdata.columns}
			data := Tabdata{entries: input}
			reduceColumns(c, &data)

			assert.EqualValues(t, testdata.expect, data.entries)
		})
	}
}

func TestNumberizeHeaders(t *testing.T) {
	data := Tabdata{
		headers: []string{"ONE", "TWO", "THREE"},
	}

	var tests = []struct {
		expect    []string
		columns   []int
		numberize bool
	}{
		{[]string{"ONE(1)", "TWO(2)", "THREE(3)"}, []int{1, 2, 3}, true},
		{[]string{"ONE(1)", "TWO(2)"}, []int{1, 2}, true},
		{[]string{"ONE", "TWO"}, []int{1, 2}, false},
	}

	for _, testdata := range tests {
		testname := fmt.Sprintf("numberize-headers-columns-%+v-nonum-%t",
			testdata.columns, testdata.numberize)

		t.Run(testname, func(t *testing.T) {
			conf := cfg.Config{Columns: "x", UseColumns: testdata.columns, Numbering: testdata.numberize}
			usedata := data
			numberizeAndReduceHeaders(conf, &usedata)

			assert.EqualValues(t, testdata.expect, usedata.headers)
		})
	}
}
