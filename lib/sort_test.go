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
	"testing"
)

func TestDuration2Seconds(t *testing.T) {
	var tests = []struct {
		dur    string
		expect int
	}{
		{"1d", 60 * 60 * 24},
		{"1h", 60 * 60},
		{"10m", 60 * 10},
		{"2h4m10s", (60 * 120) + (4 * 60) + 10},
		{"88u", 0},
		{"19t77X what?4s", 4},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("duration-%s", tt.dur)
		t.Run(testname, func(t *testing.T) {
			seconds := duration2int(tt.dur)
			if seconds != tt.expect {
				t.Errorf("got %d, want %d", seconds, tt.expect)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	var tests = []struct {
		mode string
		a    string
		b    string
		want bool
		desc bool
	}{
		// ascending
		{"numeric", "10", "20", true, false},
		{"duration", "2d4h5m", "45m", false, false},
		{"time", "12/24/2022", "1/1/1970", false, false},

		// descending
		{"numeric", "10", "20", false, true},
		{"duration", "2d4h5m", "45m", true, true},
		{"time", "12/24/2022", "1/1/1970", true, true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("compare-mode-%s-a-%s-b-%s-desc-%t", tt.mode, tt.a, tt.b, tt.desc)
		t.Run(testname, func(t *testing.T) {
			SortMode = tt.mode
			SortDescending = tt.desc
			got := compare(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("got %t, want %t", got, tt.want)
			}
		})
	}
}
