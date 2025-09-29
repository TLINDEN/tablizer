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
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/tlinden/tablizer/cfg"
)

func newData() Tabdata {
	return Tabdata{
		maxwidthHeader: 8,
		columns:        4,
		headers: []string{
			"NAME",
			"DURATION",
			"COUNT",
			"WHEN",
		},
		entries: [][]string{
			{
				"beta",
				"1d10h5m1s",
				"33",
				"3/1/2014",
			},
			{
				"alpha",
				"4h35m",
				"170",
				"2013-Feb-03",
			},
			{
				"ceta",
				"33d12h",
				"9",
				"06/Jan/2008 15:04:05 -0700",
			},
		},
	}
}

var tests = []struct {
	name string // so we can identify which one fails, can be the same
	// for multiple tests, because flags will be appended to the name
	sortby    string // empty == default
	column    int    // sort by this column (numbers start by 1)
	desc      bool   // sort in descending order, default == ascending
	numberize bool   // add header numbering
	mode      int    // shell, orgtbl, etc. empty == default: ascii
	usecol    []int  // columns to display, empty == display all
	usecolstr string // for testname, must match usecol
	expect    string // rendered output we expect
}{
	// --------------------- Default settings mode tests ``
	{
		mode:      cfg.ASCII,
		numberize: true,
		name:      "default",
		expect: `
NAME(1)	DURATION(2)	COUNT(3) WHEN(4)                    
beta   	1d10h5m1s  	33       3/1/2014                   
alpha  	4h35m      	170      2013-Feb-03                
ceta   	33d12h     	9        06/Jan/2008 15:04:05 -0700`,
	},
	{
		mode:      cfg.CSV,
		numberize: false,
		name:      "csv",
		expect: `
NAME,DURATION,COUNT,WHEN
beta,1d10h5m1s,33,3/1/2014
alpha,4h35m,170,2013-Feb-03
ceta,33d12h,9,06/Jan/2008 15:04:05 -0700`,
	},
	{
		name:      "orgtbl",
		numberize: true,
		mode:      cfg.Orgtbl,
		expect: `
+---------+-------------+----------+----------------------------+
| NAME(1) | DURATION(2) | COUNT(3) | WHEN(4)                    |
+---------+-------------+----------+----------------------------+
| beta    | 1d10h5m1s   | 33       | 3/1/2014                   |
| alpha   | 4h35m       | 170      | 2013-Feb-03                |
| ceta    | 33d12h      | 9        | 06/Jan/2008 15:04:05 -0700 |
+---------+-------------+----------+----------------------------+`,
	},
	{
		name:      "markdown",
		mode:      cfg.Markdown,
		numberize: true,
		expect: `
| NAME(1) | DURATION(2) | COUNT(3) | WHEN(4)                    |
|---------|-------------|----------|----------------------------|
| beta    | 1d10h5m1s   | 33       | 3/1/2014                   |
| alpha   | 4h35m       | 170      | 2013-Feb-03                |
| ceta    | 33d12h      | 9        | 06/Jan/2008 15:04:05 -0700 |`,
	},
	{
		name:      "shell",
		mode:      cfg.Shell,
		numberize: false,
		expect: `
NAME="beta" DURATION="1d10h5m1s" COUNT="33" WHEN="3/1/2014"
NAME="alpha" DURATION="4h35m" COUNT="170" WHEN="2013-Feb-03"
NAME="ceta" DURATION="33d12h" COUNT="9" WHEN="06/Jan/2008 15:04:05 -0700"`,
	},
	{
		name:      "yaml",
		mode:      cfg.Yaml,
		numberize: false,
		expect: `
entries:
    - count: 33
      duration: "1d10h5m1s"
      name: "beta"
      when: "3/1/2014"
    - count: 170
      duration: "4h35m"
      name: "alpha"
      when: "2013-Feb-03"
    - count: 9
      duration: "33d12h"
      name: "ceta"
      when: "06/Jan/2008 15:04:05 -0700"`,
	},
	{
		name:      "extended",
		mode:      cfg.Extended,
		numberize: true,
		expect: `
    NAME(1): beta
DURATION(2): 1d10h5m1s
   COUNT(3): 33
    WHEN(4): 3/1/2014

    NAME(1): alpha
DURATION(2): 4h35m
   COUNT(3): 170
    WHEN(4): 2013-Feb-03

    NAME(1): ceta
DURATION(2): 33d12h
   COUNT(3): 9
    WHEN(4): 06/Jan/2008 15:04:05 -0700`,
	},

	//------------------------ SORT TESTS
	{
		name:      "sortbycolumn3",
		column:    3,
		sortby:    "numeric",
		numberize: true,
		desc:      false,
		expect: `
NAME(1) DURATION(2) COUNT(3) WHEN(4)                    
ceta    33d12h      9        06/Jan/2008 15:04:05 -0700
beta    1d10h5m1s   33       3/1/2014                   
alpha   4h35m       170      2013-Feb-03`,
	},
	{
		name:      "sortbycolumn4",
		column:    4,
		sortby:    "time",
		desc:      false,
		numberize: true,
		expect: `
NAME(1) DURATION(2) COUNT(3) WHEN(4)                    
ceta    33d12h      9        06/Jan/2008 15:04:05 -0700
alpha   4h35m       170      2013-Feb-03                
beta    1d10h5m1s   33       3/1/2014`,
	},
	{
		name:      "sortbycolumn2",
		column:    2,
		sortby:    "duration",
		numberize: true,
		desc:      false,
		expect: `
NAME(1) DURATION(2) COUNT(3) WHEN(4)                    
alpha   4h35m       170      2013-Feb-03                
beta    1d10h5m1s   33       3/1/2014                   
ceta    33d12h      9        06/Jan/2008 15:04:05 -0700`,
	},

	//  -----------------------  UseColumns Tests
	{
		name:      "usecolumns",
		usecol:    []int{1, 4},
		numberize: true,
		usecolstr: "1,4",
		expect: `
NAME(1) WHEN(4)                    
beta    3/1/2014                   
alpha   2013-Feb-03                
ceta    06/Jan/2008 15:04:05 -0700`,
	},
	{
		name:      "usecolumns",
		usecol:    []int{2},
		numberize: true,
		usecolstr: "2",
		expect: `
DURATION(2)
1d10h5m1s  
4h35m      
33d12h`,
	},
	{
		name:      "usecolumns",
		usecol:    []int{3},
		numberize: true,
		usecolstr: "3",
		expect: `
COUNT(3)
33      
170     
9`,
	},
	{
		name:      "usecolumns",
		column:    0,
		usecol:    []int{1, 3},
		numberize: true,
		usecolstr: "1,3",
		expect: `
NAME(1) COUNT(3)
beta    33       
alpha   170      
ceta    9`,
	},
	{
		name:      "usecolumns",
		usecol:    []int{2, 4},
		numberize: true,
		usecolstr: "2,4",
		expect: `
DURATION(2) WHEN(4)                    
1d10h5m1s   3/1/2014                   
4h35m       2013-Feb-03                
33d12h      06/Jan/2008 15:04:05 -0700`,
	},
}

func TestPrinter(t *testing.T) {
	for _, testdata := range tests {
		testname := fmt.Sprintf("print-%s-%d-desc-%t-sortby-%s-mode-%d-usecolumns-%s-numberize-%t",
			testdata.name, testdata.column, testdata.desc, testdata.sortby,
			testdata.mode, testdata.usecolstr, testdata.numberize)

		t.Run(testname, func(t *testing.T) {
			// replaces os.Stdout, but we ignore it
			var writer bytes.Buffer

			// cmd flags
			conf := cfg.Config{
				SortDescending: testdata.desc,
				SortMode:       testdata.sortby,
				OutputMode:     testdata.mode,
				Numbering:      testdata.numberize,
				UseColumns:     testdata.usecol,
				NoColor:        true,
			}

			if testdata.column > 0 {
				conf.UseSortByColumn = []int{testdata.column}
			}

			conf.ApplyDefaults()

			// the test checks the len!
			if len(testdata.usecol) > 0 {
				conf.Columns = "yes"
			} else {
				conf.Columns = ""
			}

			data := newData()
			exp := strings.TrimSpace(testdata.expect)

			printData(&writer, conf, &data)

			got := strings.TrimSpace(writer.String())

			if got != exp {
				t.Errorf("not rendered correctly:\n+++ got:\n%s\n+++ want:\n%s",
					got, exp)
			}
		})
	}
}
