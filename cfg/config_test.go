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

package cfg

import (
	"fmt"
	//	"reflect"
	"testing"
)

func TestPrepareModeFlags(t *testing.T) {
	var tests = []struct {
		flag   Modeflag
		expect int // output (constant enum)
	}{
		// short commandline flags like -M
		{Modeflag{X: true}, Extended},
		{Modeflag{S: true}, Shell},
		{Modeflag{O: true}, Orgtbl},
		{Modeflag{Y: true}, Yaml},
		{Modeflag{M: true}, Markdown},
		{Modeflag{}, Ascii},
	}

	// FIXME: use a map for easier printing
	for _, tt := range tests {
		testname := fmt.Sprintf("PrepareModeFlags-expect-%d", tt.expect)
		t.Run(testname, func(t *testing.T) {
			c := Config{}

			c.PrepareModeFlags(tt.flag)
			if c.OutputMode != tt.expect {
				t.Errorf("got: %d, expect: %d", c.OutputMode, tt.expect)
			}
		})
	}
}

func TestPrepareSortFlags(t *testing.T) {
	var tests = []struct {
		flag   Sortmode
		expect string // output
	}{
		// short commandline flags like -M
		{Sortmode{Numeric: true}, "numeric"},
		{Sortmode{Age: true}, "duration"},
		{Sortmode{Time: true}, "time"},
		{Sortmode{}, "string"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("PrepareSortFlags-expect-%s", tt.expect)
		t.Run(testname, func(t *testing.T) {
			c := Config{}

			c.PrepareSortFlags(tt.flag)

			if c.SortMode != tt.expect {
				t.Errorf("got: %s, expect: %s", c.SortMode, tt.expect)
			}
		})
	}
}

func TestPreparePattern(t *testing.T) {
	var tests = []struct {
		pattern string
		wanterr bool
	}{
		{"[A-Z]+", false},
		{"[a-z", true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("PreparePattern-pattern-%s-wanterr-%t", tt.pattern, tt.wanterr)
		t.Run(testname, func(t *testing.T) {
			c := Config{}

			err := c.PreparePattern(tt.pattern)

			if err != nil {
				if !tt.wanterr {
					t.Errorf("PreparePattern returned error: %s", err)
				}
			}
		})
	}
}
