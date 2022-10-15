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
	"github.com/gookit/color"
	//"github.com/xo/terminfo"
)

var (
	// command line flags
	Debug            bool
	XtendedOut       bool
	NoNumbering      bool
	ShowVersion      bool
	Columns          string
	UseColumns       []int
	DefaultSeparator string = `(\s\s+|\t)`
	Separator        string = `(\s\s+|\t)`
	OutflagExtended  bool
	OutflagMarkdown  bool
	OutflagOrgtable  bool
	OutflagShell     bool
	OutputMode       string
	InvertMatch      bool
	Pattern          string

	/*
	 FIXME: make configurable somehow, config file or ENV
	 see https://github.com/gookit/color will be set by
	 io.ProcessFiles() according to currently supported
	 color mode.
	*/
	MatchFG string
	MatchBG string
	NoColor bool

	// colors to be used per supported color mode
	Colors = map[color.Level]map[string]string{
		color.Level16: {
			"bg": "green", "fg": "black",
		},
		color.Level256: {
			"bg": "lightGreen", "fg": "black",
		},
		color.LevelRgb: {
			// FIXME: maybe use something nicer
			"bg": "lightGreen", "fg": "black",
		},
	}

	// used for validation
	validOutputmodes = "(orgtbl|markdown|extended|ascii)"

	// main program version
	Version = "v1.0.9"

	// generated  version string, used  by -v contains  lib.Version on
	//  main  branch,   and  lib.Version-$branch-$lastcommit-$date  on
	// development branch
	VERSION string

	// sorting
	SortByColumn   int
	SortDescending bool
	SortNumeric    bool
	SortTime       bool
	SortAge        bool
)

// contains a whole parsed table
type Tabdata struct {
	maxwidthHeader int      // longest header
	maxwidthPerCol []int    // max width per column
	columns        int      // count
	headers        []string // [ "ID", "NAME", ...]
	entries        [][]string
}
