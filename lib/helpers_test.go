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

func TestArrayContains(t *testing.T) {
	var tests = []struct {
		list   []int
		search int
		want   bool
	}{
		{[]int{1, 2, 3}, 2, true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d,%d,%t", tt.list, tt.search, tt.want)
		t.Run(testname, func(t *testing.T) {
			answer := contains(tt.list, tt.search)
			if answer != tt.want {
				t.Errorf("got %t, want %t", answer, tt.want)
			}
		})
	}
}
