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
	"os"
	"regexp"

	"github.com/glycerine/zygomys/zygo"
	"github.com/gookit/color"
)

const DefaultSeparator string = `(\s\s+|\t)`
const Version string = "v1.0.17"

var DefaultLoadPath string = os.Getenv("HOME") + "/.config/tablizer/lisp"

var VERSION string // maintained by -x

type Config struct {
	Debug          bool
	NoNumbering    bool
	NoHeaders      bool
	Columns        string
	UseColumns     []int
	Separator      string
	OutputMode     int
	InvertMatch    bool
	Pattern        string
	PatternR       *regexp.Regexp
	UseFuzzySearch bool

	SortMode       string
	SortDescending bool
	SortByColumn   int

	/*
	 FIXME: make configurable somehow, config file or ENV
	 see https://github.com/gookit/color.
	*/
	ColorStyle color.Style

	NoColor bool

	// special  case: we use the  config struct to transport  the lisp
	// env trough the program
	Lisp *zygo.Zlisp

	// a path containing lisp scripts to be loaded on startup
	LispLoadPath string
}

// maps outputmode short flags to output mode, ie. -O => -o orgtbl
type Modeflag struct {
	X bool
	O bool
	M bool
	S bool
	Y bool
	A bool
	C bool
}

// used for switching printers
const (
	Extended = iota + 1
	Orgtbl
	Markdown
	Shell
	Yaml
	CSV
	Ascii
)

// various sort types
type Sortmode struct {
	Numeric bool
	Time    bool
	Age     bool
}

// valid lisp hooks
var ValidHooks []string

// default color schemes
func Colors() map[color.Level]map[string]color.Color {
	return map[color.Level]map[string]color.Color{
		color.Level16: {
			"bg": color.BgGreen, "fg": color.FgBlack,
		},
		color.Level256: {
			"bg": color.BgLightGreen, "fg": color.FgBlack,
		},
		color.LevelRgb: {
			// FIXME: maybe use something nicer
			"bg": color.BgLightGreen, "fg": color.FgBlack,
		},
	}
}

// find supported color mode, modifies config based on constants
func (c *Config) DetermineColormode() {
	if !isTerminal(os.Stdout) {
		color.Disable()
	} else {
		level := color.TermColorLevel()
		colors := Colors()
		c.ColorStyle = color.New(colors[level]["bg"], colors[level]["fg"])
	}
}

// Return true if current terminal is interactive
func isTerminal(f *os.File) bool {
	o, _ := f.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		return true
	} else {
		return false
	}
}

func Getversion() string {
	// main program version

	// generated  version string, used  by -v contains  lib.Version on
	//  main  branch,   and  lib.Version-$branch-$lastcommit-$date  on
	// development branch

	return fmt.Sprintf("This is tablizer version %s", VERSION)
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

func (conf *Config) PrepareModeFlags(flag Modeflag) {
	switch {
	case flag.X:
		conf.OutputMode = Extended
	case flag.O:
		conf.OutputMode = Orgtbl
	case flag.M:
		conf.OutputMode = Markdown
	case flag.S:
		conf.OutputMode = Shell
	case flag.Y:
		conf.OutputMode = Yaml
	case flag.C:
		conf.OutputMode = CSV
	default:
		conf.OutputMode = Ascii
	}
}

func (c *Config) CheckEnv() {
	// check for environment vars, command line flags have precedence,
	// NO_COLOR is being checked by the color module itself.
	if !c.NoNumbering {
		_, set := os.LookupEnv("T_NO_HEADER_NUMBERING")
		if set {
			c.NoNumbering = true
		}
	}

	if len(c.Columns) == 0 {
		cols := os.Getenv("T_COLUMNS")
		if len(cols) > 1 {
			c.Columns = cols
		}
	}
}

func (c *Config) ApplyDefaults() {
	// mode specific defaults
	if c.OutputMode == Yaml || c.OutputMode == CSV {
		c.NoNumbering = true
	}

	ValidHooks = []string{"filter", "process", "transpose", "append"}
}

func (c *Config) PreparePattern(pattern string) error {
	PatternR, err := regexp.Compile(pattern)

	if err != nil {
		return errors.Unwrap(fmt.Errorf("Regexp pattern %s is invalid: %w", c.Pattern, err))
	}

	c.PatternR = PatternR
	c.Pattern = pattern

	return nil
}
