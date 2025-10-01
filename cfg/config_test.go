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

	"github.com/stretchr/testify/assert"
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
		{Modeflag{}, ASCII},
	}

	// FIXME: use a map for easier printing
	for _, testdata := range tests {
		testname := fmt.Sprintf("PrepareModeFlags-expect-%d", testdata.expect)
		t.Run(testname, func(t *testing.T) {
			conf := Config{}

			conf.PrepareModeFlags(testdata.flag)

			assert.EqualValues(t, testdata.expect, conf.OutputMode)
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

	for _, testdata := range tests {
		testname := fmt.Sprintf("PrepareSortFlags-expect-%s", testdata.expect)
		t.Run(testname, func(t *testing.T) {
			conf := Config{}

			conf.PrepareSortFlags(testdata.flag)

			assert.EqualValues(t, testdata.expect, conf.SortMode)
		})
	}
}

func TestPreparePattern(t *testing.T) {
	var tests = []struct {
		patterns  []*Pattern
		name      string
		wanterror bool
		wanticase bool
		wantneg   bool
	}{
		{
			[]*Pattern{{Pattern: "[A-Z]+"}},
			"simple",
			false,
			false,
			false,
		},
		{
			[]*Pattern{{Pattern: "[a-z"}},
			"regfail",
			true,
			false,
			false,
		},
		{
			[]*Pattern{{Pattern: "/[A-Z]+/i"}},
			"icase",
			false,
			true,
			false,
		},
		{
			[]*Pattern{{Pattern: "/[A-Z]+/!"}},
			"negate",
			false,
			false,
			true,
		},
		{
			[]*Pattern{{Pattern: "/[A-Z]+/!i"}},
			"negicase",
			false,
			true,
			true,
		},
	}

	for _, testdata := range tests {
		testname := fmt.Sprintf("PreparePattern-pattern-%s-wanterr-%t", testdata.name, testdata.wanterror)
		t.Run(testname, func(t *testing.T) {
			conf := Config{}

			err := conf.PreparePattern(testdata.patterns)

			if testdata.wanterror {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
