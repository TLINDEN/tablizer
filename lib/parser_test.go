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

	"github.com/tlinden/tablizer/cfg"
)

var input = []struct {
	name      string
	text      string
	separator string
}{
	{
		name:      "tabular-data",
		separator: cfg.DefaultSeparator,
		text: `
ONE    TWO    THREE  
asd    igig   cxxxncnc  
19191  EDD 1  X`,
	},
	{
		name:      "csv-data",
		separator: ",",
		text: `
ONE,TWO,THREE
asd,igig,cxxxncnc
19191,"EDD 1",X`,
	},
}

func TestParser(t *testing.T) {
	data := Tabdata{
		maxwidthHeader: 5,
		columns:        3,
		headers: []string{
			"ONE", "TWO", "THREE",
		},
		entries: [][]string{
			{"asd", "igig", "cxxxncnc"},
			{"19191", "EDD 1", "X"},
		},
	}

	for _, testdata := range input {
		testname := fmt.Sprintf("parse-%s", testdata.name)
		t.Run(testname, func(t *testing.T) {
			readFd := strings.NewReader(strings.TrimSpace(testdata.text))
			conf := cfg.Config{Separator: testdata.separator}
			gotdata, err := Parse(conf, readFd)

			if err != nil {
				t.Errorf("Parser returned error: %s\nData processed so far: %+v", err, gotdata)
			}

			if !reflect.DeepEqual(data, gotdata) {
				t.Errorf("Parser returned invalid data\nExp: %+v\nGot: %+v\n",
					data, gotdata)
			}
		})
	}
}

func TestParserPatternmatching(t *testing.T) {
	var tests = []struct {
		name     string
		entries  [][]string
		patterns []*cfg.Pattern
		invert   bool
		want     bool
	}{
		{
			name: "match",
			entries: [][]string{
				{"asd", "igig", "cxxxncnc"},
			},
			patterns: []*cfg.Pattern{{Pattern: "ig"}},
			invert:   false,
		},
		{
			name: "invert",
			entries: [][]string{
				{"19191", "EDD 1", "X"},
			},
			patterns: []*cfg.Pattern{{Pattern: "ig"}},
			invert:   true,
		},
	}

	for _, inputdata := range input {
		for _, testdata := range tests {
			testname := fmt.Sprintf("parse-%s-with-pattern-%s-inverted-%t",
				inputdata.name, testdata.name, testdata.invert)
			t.Run(testname, func(t *testing.T) {
				conf := cfg.Config{
					InvertMatch: testdata.invert,
					Patterns:    testdata.patterns,
					Separator:   inputdata.separator,
				}

				_ = conf.PreparePattern(testdata.patterns)

				readFd := strings.NewReader(strings.TrimSpace(inputdata.text))
				gotdata, err := Parse(conf, readFd)

				if err != nil {
					if !testdata.want {
						t.Errorf("Parser returned error: %s\nData processed so far: %+v",
							err, gotdata)
					}
				} else {
					if !reflect.DeepEqual(testdata.entries, gotdata.entries) {
						t.Errorf("Parser returned invalid data (pattern: %s, invert: %t)\nExp: %+v\nGot: %+v\n",
							testdata.name, testdata.invert, testdata.entries, gotdata.entries)
					}
				}
			})
		}
	}
}

func TestParserIncompleteRows(t *testing.T) {
	data := Tabdata{
		maxwidthHeader: 5,
		columns:        3,
		headers: []string{
			"ONE", "TWO", "THREE",
		},
		entries: [][]string{
			{"asd", "igig", ""},
			{"19191", "EDD 1", "X"},
		},
	}

	table := `
ONE    TWO    THREE  
asd    igig
19191  EDD 1  X`

	readFd := strings.NewReader(strings.TrimSpace(table))
	conf := cfg.Config{Separator: cfg.DefaultSeparator}
	gotdata, err := Parse(conf, readFd)

	if err != nil {
		t.Errorf("Parser returned error: %s\nData processed so far: %+v", err, gotdata)
	}

	if !reflect.DeepEqual(data, gotdata) {
		t.Errorf("Parser returned invalid data, Regex: %s\nExp: %+v\nGot: %+v\n",
			conf.Separator, data, gotdata)
	}
}
