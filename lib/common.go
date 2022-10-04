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

var (
	// command line flags
	Debug           bool
	XtendedOut      bool
	NoNumbering     bool
	ShowVersion     bool
	Columns         string
	UseColumns      []int
	Separator       string
	OutflagExtended bool
	OutflagMarkdown bool
	OutflagOrgtable bool
	OutflagShell    bool
	OutputMode      string

	// used for validation
	validOutputmodes = "(orgtbl|markdown|extended|ascii)"

	// main program version
	Version = "v1.0.4"

	// generated  version string, used  by -v contains  lib.Version on
	//  main  branch,   and  lib.Version-$branch-$lastcommit-$date  on
	// development branch
	VERSION string
)
