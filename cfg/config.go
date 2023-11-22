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
	"log"
	"os"
	"regexp"

	"github.com/glycerine/zygomys/zygo"
	"github.com/gookit/color"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

const DefaultSeparator string = `(\s\s+|\t)`
const Version string = "v1.1.0"

var DefaultLoadPath string = os.Getenv("HOME") + "/.config/tablizer/lisp"
var DefaultConfigfile string = os.Getenv("HOME") + "/.config/tablizer/config"

var VERSION string // maintained by -x

// public config, set via config file or using defaults
type Configuration struct {
	FG             string `hcl:"FG"`
	BG             string `hcl:"BG"`
	HighlightFG    string `hcl:"HighlightFG"`
	HighlightBG    string `hcl:"HighlightBG"`
	NoHighlightFG  string `hcl:"NoHighlightFG"`
	NoHighlightBG  string `hcl:"NoHighlightBG"`
	HighlightHdrFG string `hcl:"HighlightHdrFG"`
	HighlightHdrBG string `hcl:"HighlightHdrBG"`
}

// internal config
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
	UseHighlight   bool

	SortMode       string
	SortDescending bool
	SortByColumn   int

	/*
	 FIXME: make configurable somehow, config file or ENV
	 see https://github.com/gookit/color.
	*/
	ColorStyle        color.Style
	HighlightStyle    color.Style
	NoHighlightStyle  color.Style
	HighlightHdrStyle color.Style

	NoColor bool

	// special  case: we use the  config struct to transport  the lisp
	// env trough the program
	Lisp *zygo.Zlisp

	// a path containing lisp scripts to be loaded on startup
	LispLoadPath string

	// config file, optional
	Configfile string

	Configuration Configuration
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
func (c *Config) Colors() map[color.Level]map[string]color.Color {
	colors := map[color.Level]map[string]color.Color{
		color.Level16: {
			"bg": color.BgGreen, "fg": color.FgWhite,
			"hlbg": color.BgGray, "hlfg": color.FgWhite,
		},
		color.Level256: {
			"bg": color.BgLightGreen, "fg": color.FgWhite,
			"hlbg": color.BgLightBlue, "hlfg": color.FgWhite,
		},
		color.LevelRgb: {
			"bg": color.BgLightGreen, "fg": color.FgWhite,
			"hlbg": color.BgHiGreen, "hlfg": color.FgWhite,
			"nohlbg": color.BgWhite, "nohlfg": color.FgLightGreen,
			"hdrbg": color.BgBlue, "hdrfg": color.FgWhite,
		},
	}

	if len(c.Configuration.BG) > 0 {
		colors[color.Level16]["bg"] = ColorStringToBGColor(c.Configuration.BG)
		colors[color.Level256]["bg"] = ColorStringToBGColor(c.Configuration.BG)
		colors[color.LevelRgb]["bg"] = ColorStringToBGColor(c.Configuration.BG)
	}

	if len(c.Configuration.FG) > 0 {
		colors[color.Level16]["fg"] = ColorStringToColor(c.Configuration.FG)
		colors[color.Level256]["fg"] = ColorStringToColor(c.Configuration.FG)
		colors[color.LevelRgb]["fg"] = ColorStringToColor(c.Configuration.FG)
	}

	if len(c.Configuration.HighlightBG) > 0 {
		colors[color.Level16]["hlbg"] = ColorStringToBGColor(c.Configuration.HighlightBG)
		colors[color.Level256]["hlbg"] = ColorStringToBGColor(c.Configuration.HighlightBG)
		colors[color.LevelRgb]["hlbg"] = ColorStringToBGColor(c.Configuration.HighlightBG)
	}

	if len(c.Configuration.HighlightFG) > 0 {
		colors[color.Level16]["hlfg"] = ColorStringToColor(c.Configuration.HighlightFG)
		colors[color.Level256]["hlfg"] = ColorStringToColor(c.Configuration.HighlightFG)
		colors[color.LevelRgb]["hlfg"] = ColorStringToColor(c.Configuration.HighlightFG)
	}

	if len(c.Configuration.NoHighlightBG) > 0 {
		colors[color.Level16]["nohlbg"] = ColorStringToBGColor(c.Configuration.NoHighlightBG)
		colors[color.Level256]["nohlbg"] = ColorStringToBGColor(c.Configuration.NoHighlightBG)
		colors[color.LevelRgb]["nohlbg"] = ColorStringToBGColor(c.Configuration.NoHighlightBG)
	}

	if len(c.Configuration.NoHighlightFG) > 0 {
		colors[color.Level16]["nohlfg"] = ColorStringToColor(c.Configuration.NoHighlightFG)
		colors[color.Level256]["nohlfg"] = ColorStringToColor(c.Configuration.NoHighlightFG)
		colors[color.LevelRgb]["nohlfg"] = ColorStringToColor(c.Configuration.NoHighlightFG)
	}

	if len(c.Configuration.HighlightHdrBG) > 0 {
		colors[color.Level16]["hdrbg"] = ColorStringToBGColor(c.Configuration.HighlightHdrBG)
		colors[color.Level256]["hdrbg"] = ColorStringToBGColor(c.Configuration.HighlightHdrBG)
		colors[color.LevelRgb]["hdrbg"] = ColorStringToBGColor(c.Configuration.HighlightHdrBG)
	}

	if len(c.Configuration.HighlightHdrFG) > 0 {
		colors[color.Level16]["hdrfg"] = ColorStringToColor(c.Configuration.HighlightHdrFG)
		colors[color.Level256]["hdrfg"] = ColorStringToColor(c.Configuration.HighlightHdrFG)
		colors[color.LevelRgb]["hdrfg"] = ColorStringToColor(c.Configuration.HighlightHdrFG)
	}

	return colors
}

// find supported color mode, modifies config based on constants
func (c *Config) DetermineColormode() {
	if !isTerminal(os.Stdout) {
		color.Disable()
	} else {
		level := color.TermColorLevel()
		colors := c.Colors()

		c.ColorStyle = color.New(colors[level]["bg"], colors[level]["fg"])
		c.HighlightStyle = color.New(colors[level]["hlbg"], colors[level]["hlfg"])
		c.NoHighlightStyle = color.New(colors[level]["nohlbg"], colors[level]["nohlfg"])
		c.HighlightHdrStyle = color.New(colors[level]["hdrbg"], colors[level]["hdrfg"])
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

func (c *Config) ParseConfigfile() error {
	if path, err := os.Stat(c.Configfile); !os.IsNotExist(err) {
		if !path.IsDir() {
			configstring, err := os.ReadFile(path.Name())
			if err != nil {
				return err
			}

			err = hclsimple.Decode(
				path.Name(), []byte(configstring),
				nil, &c.Configuration,
			)
			if err != nil {
				log.Fatalf("Failed to load configuration: %s", err)
			}
		}
	}

	return nil
}

// translate color string to internal color value
func ColorStringToColor(colorname string) color.Color {
	for name, color := range color.FgColors {
		if name == colorname {
			return color
		}
	}

	for name, color := range color.ExFgColors {
		if name == colorname {
			return color
		}
	}

	return color.Normal
}

// same, for background colors
func ColorStringToBGColor(colorname string) color.Color {
	for name, color := range color.BgColors {
		if name == colorname {
			return color
		}
	}

	for name, color := range color.ExBgColors {
		if name == colorname {
			return color
		}
	}

	return color.Normal
}
