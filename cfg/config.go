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
	"errors"
	"fmt"
	"github.com/gookit/color"
	"regexp"
)

const DefaultSeparator string = `(\s\s+|\t)`
const ValidOutputModes string = "(orgtbl|markdown|extended|ascii|yaml)"
const Version string = "v1.0.11"

var VERSION string // maintained by -x

type Config struct {
	Debug       bool
	NoNumbering bool
	Columns     string
	UseColumns  []int
	Separator   string
	OutputMode  string
	InvertMatch bool
	Pattern     string

	SortMode       string
	SortDescending bool
	SortByColumn   int

	/*
	 FIXME: make configurable somehow, config file or ENV
	 see https://github.com/gookit/color will be set by
	 io.ProcessFiles() according to currently supported
	 color mode.
	*/
	MatchFG string
	MatchBG string
	NoColor bool
}

// maps outputmode short flags to output mode, ie. -O => -o orgtbl
type Modeflag struct {
	X bool
	O bool
	M bool
	S bool
	Y bool
	A bool
}

// various sort types
type Sortmode struct {
	Numeric bool
	Time    bool
	Age     bool
}

func Colors() map[color.Level]map[string]string {
	// default color schemes
	return map[color.Level]map[string]string{
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
}

func Getversion() string {
	// main program version

	// generated  version string, used  by -v contains  lib.Version on
	//  main  branch,   and  lib.Version-$branch-$lastcommit-$date  on
	// development branch

	return fmt.Sprintf("This is tablizer version %s", VERSION)
}

func (conf *Config) PrepareModeFlags(flag Modeflag, mode string) error {
	if len(mode) == 0 {
		// associate short flags like -X with mode selector
		switch {
		case flag.X:
			conf.OutputMode = "extended"
		case flag.M:
			conf.OutputMode = "markdown"
		case flag.O:
			conf.OutputMode = "orgtbl"
		case flag.S:
			conf.OutputMode = "shell"
			conf.NoNumbering = true
		case flag.Y:
			conf.OutputMode = "yaml"
			conf.NoNumbering = true
		default:
			conf.OutputMode = "ascii"
		}
	} else {
		r, err := regexp.Compile(ValidOutputModes)

		if err != nil {
			return errors.New("Failed to validate output mode spec!")
		}

		match := r.MatchString(mode)

		if !match {
			return errors.New("Invalid output mode!")
		}
	}

	return nil
}

func (conf *Config) PrepareSortFlags(flag Sortmode) {
	switch {
	case flag.Numeric:
		conf.SortMode = "numeric"
	case flag.Age:
		conf.SortMode = "duration"
	case flag.Time:
		conf.SortMode = "time"
	default:
		conf.SortMode = "string"
	}
}
