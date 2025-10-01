/*
Copyright © 2024-2025 Thomas von Dein

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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tlinden/tablizer/cfg"
)

func TestMatchPattern(t *testing.T) {
	var input = []struct {
		name     string
		fuzzy    bool
		patterns []*cfg.Pattern
		line     string
	}{
		{
			name:     "normal",
			patterns: []*cfg.Pattern{{Pattern: "haus"}},
			line:     "hausparty",
		},
		{
			name:     "fuzzy",
			patterns: []*cfg.Pattern{{Pattern: "hpt"}},
			line:     "haus-party-termin",
			fuzzy:    true,
		},
	}

	for _, inputdata := range input {
		testname := fmt.Sprintf("match-pattern-%s", inputdata.name)

		t.Run(testname, func(t *testing.T) {
			conf := cfg.Config{}

			if inputdata.fuzzy {
				conf.UseFuzzySearch = true
			}

			err := conf.PreparePattern(inputdata.patterns)

			assert.NoError(t, err)

			res := matchPattern(conf, inputdata.line)
			assert.EqualValues(t, true, res)
		})
	}
}

func TestFilterByFields(t *testing.T) {
	data := Tabdata{
		headers: []string{
			"ONE", "TWO", "THREE",
		},
		entries: [][]string{
			{"asd", "igig", "cxxxncnc"},
			{"19191", "EDD 1", "x"},
			{"8d8", "AN 1", "y"},
		},
	}

	var input = []struct {
		name   string
		filter []string
		expect Tabdata
		invert bool
	}{
		{
			name:   "one-field",
			filter: []string{"one=19"},
			expect: Tabdata{
				headers: []string{
					"ONE", "TWO", "THREE",
				},
				entries: [][]string{
					{"19191", "EDD 1", "x"},
				},
			},
		},

		{
			name:   "one-field-negative",
			filter: []string{"one!=asd"},
			expect: Tabdata{
				headers: []string{
					"ONE", "TWO", "THREE",
				},
				entries: [][]string{
					{"19191", "EDD 1", "x"},
					{"8d8", "AN 1", "y"},
				},
			},
		},

		{
			name:   "one-field-inverted",
			filter: []string{"one=19"},
			invert: true,
			expect: Tabdata{
				headers: []string{
					"ONE", "TWO", "THREE",
				},
				entries: [][]string{
					{"asd", "igig", "cxxxncnc"},
					{"8d8", "AN 1", "y"},
				},
			},
		},

		{
			name:   "many-fields",
			filter: []string{"one=19", "two=DD"},
			expect: Tabdata{
				headers: []string{
					"ONE", "TWO", "THREE",
				},
				entries: [][]string{
					{"19191", "EDD 1", "x"},
				},
			},
		},

		{
			name:   "many-fields-inverted",
			filter: []string{"one=19", "two=DD"},
			invert: true,
			expect: Tabdata{
				headers: []string{
					"ONE", "TWO", "THREE",
				},
				entries: [][]string{
					{"asd", "igig", "cxxxncnc"},
					{"8d8", "AN 1", "y"},
				},
			},
		},
	}

	for _, inputdata := range input {
		testname := fmt.Sprintf("filter-by-fields-%s", inputdata.name)

		t.Run(testname, func(t *testing.T) {
			conf := cfg.Config{Rawfilters: inputdata.filter, InvertMatch: inputdata.invert}

			err := conf.PrepareFilters()

			assert.NoError(t, err)

			data, _, _ := FilterByFields(conf, &data)

			assert.EqualValues(t, inputdata.expect, *data)
		})
	}
}
