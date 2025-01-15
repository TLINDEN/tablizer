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

	"github.com/tlinden/tablizer/cfg"
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

	for _, testdata := range tests {
		testname := fmt.Sprintf("duration-%s", testdata.dur)
		t.Run(testname, func(t *testing.T) {
			seconds := duration2int(testdata.dur)
			if seconds != testdata.expect {
				t.Errorf("got %d, want %d", seconds, testdata.expect)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	var tests = []struct {
		mode string
		a    string
		b    string
		want int
		desc bool
	}{
		// ascending
		{"numeric", "10", "20", 0, false},
		{"duration", "2d4h5m", "45m", 1, false},
		{"time", "12/24/2022", "1/1/1970", 1, false},

		// descending
		{"numeric", "10", "20", 1, true},
		{"duration", "2d4h5m", "45m", 0, true},
		{"time", "12/24/2022", "1/1/1970", 0, true},
	}

	for _, testdata := range tests {
		testname := fmt.Sprintf("compare-mode-%s-a-%s-b-%s-desc-%t",
			testdata.mode, testdata.a, testdata.b, testdata.desc)

		t.Run(testname, func(t *testing.T) {
			c := cfg.Config{SortMode: testdata.mode, SortDescending: testdata.desc}
			got := compare(&c, testdata.a, testdata.b)
			if got != testdata.want {
				t.Errorf("got %d, want %d", got, testdata.want)
			}
		})
	}
}
