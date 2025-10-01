/*
Copyright © 2022-2025 Thomas von Dein

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
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
			gotdata, err := wrapValidateParser(conf, readFd)

			assert.NoError(t, err)
			assert.EqualValues(t, data, gotdata)
		})
	}
}

func TestParserPatternmatching(t *testing.T) {
	var tests = []struct {
		name      string
		entries   [][]string
		patterns  []*cfg.Pattern
		invert    bool
		wanterror bool
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
				data, err := wrapValidateParser(conf, readFd)

				if testdata.wanterror {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.EqualValues(t, testdata.entries, data.entries)
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
	gotdata, err := wrapValidateParser(conf, readFd)

	assert.NoError(t, err)
	assert.EqualValues(t, data, gotdata)
}

func TestParserJSONInput(t *testing.T) {
	var tests = []struct {
		name      string
		input     string
		expect    Tabdata
		wanterror bool // true: expect fail, false: expect success
	}{
		{
			// too deep nesting
			name:      "invalidjson",
			wanterror: true,
			input: `[
  {
    "item": {
       "NAME": "postgres-operator-7f4c7c8485-ntlns",
       "READY": "1/1",
       "STATUS": "Running",
       "RESTARTS": "0",
       "AGE": "24h"
    }
  }
`,
			expect: Tabdata{},
		},

		{
			// one field missing + different order
			// but shall not fail
			name:      "kgpfail",
			wanterror: false,
			input: `[
  {
    "NAME": "postgres-operator-7f4c7c8485-ntlns",
    "READY": "1/1",
    "STATUS": "Running",
    "RESTARTS": "0",
    "AGE": "24h"
  },
  {
    "NAME": "wal-g-exporter-778dcd95f5-wcjzn",
    "RESTARTS": "0",
    "READY": "1/1",
    "AGE": "24h"
  }
]`,
			expect: Tabdata{
				columns: 5,
				headers: []string{"NAME", "READY", "STATUS", "RESTARTS", "AGE"},
				entries: [][]string{
					[]string{
						"postgres-operator-7f4c7c8485-ntlns",
						"1/1",
						"Running",
						"0",
						"24h",
					},
					[]string{
						"wal-g-exporter-778dcd95f5-wcjzn",
						"1/1",
						"",
						"0",
						"24h",
					},
				},
			},
		},

		{
			name:      "kgp",
			wanterror: false,
			input: `[
  {
    "NAME": "postgres-operator-7f4c7c8485-ntlns",
    "READY": "1/1",
    "STATUS": "Running",
    "RESTARTS": "0",
    "AGE": "24h"
  },
  {
    "NAME": "wal-g-exporter-778dcd95f5-wcjzn",
    "STATUS": "Running",
    "READY": "1/1",
    "RESTARTS": "0",
    "AGE": "24h"
  }
]`,
			expect: Tabdata{
				columns: 5,
				headers: []string{"NAME", "READY", "STATUS", "RESTARTS", "AGE"},
				entries: [][]string{
					[]string{
						"postgres-operator-7f4c7c8485-ntlns",
						"1/1",
						"Running",
						"0",
						"24h",
					},
					[]string{
						"wal-g-exporter-778dcd95f5-wcjzn",
						"1/1",
						"Running",
						"0",
						"24h",
					},
				},
			},
		},
	}

	for _, testdata := range tests {
		testname := fmt.Sprintf("parse-json-%s", testdata.name)
		t.Run(testname, func(t *testing.T) {
			conf := cfg.Config{InputJSON: true}

			readFd := strings.NewReader(strings.TrimSpace(testdata.input))
			data, err := wrapValidateParser(conf, readFd)

			if testdata.wanterror {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, testdata.expect, data)
			}
		})
	}
}

func wrapValidateParser(conf cfg.Config, input io.Reader) (Tabdata, error) {
	data, err := Parse(conf, input)

	if err != nil {
		return data, err
	}

	err = ValidateConsistency(&data)

	return data, err
}
