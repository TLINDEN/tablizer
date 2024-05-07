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

	for _, in := range input {
		testname := fmt.Sprintf("parse-%s", in.name)
		t.Run(testname, func(t *testing.T) {
			readFd := strings.NewReader(strings.TrimSpace(in.text))
			c := cfg.Config{Separator: in.separator}
			gotdata, err := Parse(c, readFd)

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
		entries [][]string
		pattern string
		invert  bool
		want    bool
	}{
		{
			entries: [][]string{
				{"asd", "igig", "cxxxncnc"},
			},
			pattern: "ig",
			invert:  false,
		},
		{
			entries: [][]string{
				{"19191", "EDD 1", "X"},
			},
			pattern: "ig",
			invert:  true,
		},
	}

	for _, in := range input {
		for _, tt := range tests {
			testname := fmt.Sprintf("parse-%s-with-pattern-%s-inverted-%t",
				in.name, tt.pattern, tt.invert)
			t.Run(testname, func(t *testing.T) {
				conf := cfg.Config{InvertMatch: tt.invert, Pattern: tt.pattern,
					Separator: in.separator}

				_ = conf.PreparePattern(tt.pattern)

				readFd := strings.NewReader(strings.TrimSpace(in.text))
				gotdata, err := Parse(conf, readFd)

				if err != nil {
					if !tt.want {
						t.Errorf("Parser returned error: %s\nData processed so far: %+v",
							err, gotdata)
					}
				} else {
					if !reflect.DeepEqual(tt.entries, gotdata.entries) {
						t.Errorf("Parser returned invalid data (pattern: %s, invert: %t)\nExp: %+v\nGot: %+v\n",
							tt.pattern, tt.invert, tt.entries, gotdata.entries)
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
	c := cfg.Config{Separator: cfg.DefaultSeparator}
	gotdata, err := Parse(c, readFd)

	if err != nil {
		t.Errorf("Parser returned error: %s\nData processed so far: %+v", err, gotdata)
	}

	if !reflect.DeepEqual(data, gotdata) {
		t.Errorf("Parser returned invalid data, Regex: %s\nExp: %+v\nGot: %+v\n",
			c.Separator, data, gotdata)
	}
}
