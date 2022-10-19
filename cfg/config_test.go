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
		mode   string // input, if any
		expect string // output
		want   bool
	}{
		// short commandline flags like -M
		{Modeflag{X: true}, "", "extended", false},
		{Modeflag{S: true}, "", "shell", false},
		{Modeflag{O: true}, "", "orgtbl", false},
		{Modeflag{Y: true}, "", "yaml", false},
		{Modeflag{M: true}, "", "markdown", false},
		{Modeflag{}, "", "ascii", false},

		// long flags like -o yaml
		{Modeflag{}, "extended", "extended", false},
		{Modeflag{}, "shell", "shell", false},
		{Modeflag{}, "orgtbl", "orgtbl", false},
		{Modeflag{}, "yaml", "yaml", false},
		{Modeflag{}, "markdown", "markdown", false},

		// failing
		{Modeflag{}, "blah", "", true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("PrepareModeFlags-flags-mode-%s-expect-%s-want-%t",
			tt.mode, tt.expect, tt.want)
		t.Run(testname, func(t *testing.T) {
			c := Config{OutputMode: tt.mode}

			// check either flag or pre filled mode, whatever is defined in tt
			err := c.PrepareModeFlags(tt.flag, tt.mode)
			if err != nil {
				if !tt.want {
					// expect to fail
					t.Fatalf("PrepareModeFlags returned unexpected error: %s", err)
				}
			} else {
				if c.OutputMode != tt.expect {
					t.Errorf("got: %s, expect: %s", c.OutputMode, tt.expect)
				}
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
