/*
Copyright © 2022-2024 Thomas von Dein

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
	"strings"

	"github.com/glycerine/zygomys/zygo"
	"github.com/gookit/color"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

const DefaultSeparator string = `(\s\s+|\t)`
const Version string = "v1.2.3"
const MAXPARTS = 2

var DefaultLoadPath = os.Getenv("HOME") + "/.config/tablizer/lisp"
var DefaultConfigfile = os.Getenv("HOME") + "/.config/tablizer/config"

var VERSION string // maintained by -x

// public config, set via config file or using defaults
type Settings struct {
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

	Settings Settings

	// used for field filtering
	Rawfilters []string
	Filters    map[string]*regexp.Regexp
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
	ASCII
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
func (conf *Config) Colors() map[color.Level]map[string]color.Color {
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

	if len(conf.Settings.BG) > 0 {
		colors[color.Level16]["bg"] = ColorStringToBGColor(conf.Settings.BG)
		colors[color.Level256]["bg"] = ColorStringToBGColor(conf.Settings.BG)
		colors[color.LevelRgb]["bg"] = ColorStringToBGColor(conf.Settings.BG)
	}

	if len(conf.Settings.FG) > 0 {
		colors[color.Level16]["fg"] = ColorStringToColor(conf.Settings.FG)
		colors[color.Level256]["fg"] = ColorStringToColor(conf.Settings.FG)
		colors[color.LevelRgb]["fg"] = ColorStringToColor(conf.Settings.FG)
	}

	if len(conf.Settings.HighlightBG) > 0 {
		colors[color.Level16]["hlbg"] = ColorStringToBGColor(conf.Settings.HighlightBG)
		colors[color.Level256]["hlbg"] = ColorStringToBGColor(conf.Settings.HighlightBG)
		colors[color.LevelRgb]["hlbg"] = ColorStringToBGColor(conf.Settings.HighlightBG)
	}

	if len(conf.Settings.HighlightFG) > 0 {
		colors[color.Level16]["hlfg"] = ColorStringToColor(conf.Settings.HighlightFG)
		colors[color.Level256]["hlfg"] = ColorStringToColor(conf.Settings.HighlightFG)
		colors[color.LevelRgb]["hlfg"] = ColorStringToColor(conf.Settings.HighlightFG)
	}

	if len(conf.Settings.NoHighlightBG) > 0 {
		colors[color.Level16]["nohlbg"] = ColorStringToBGColor(conf.Settings.NoHighlightBG)
		colors[color.Level256]["nohlbg"] = ColorStringToBGColor(conf.Settings.NoHighlightBG)
		colors[color.LevelRgb]["nohlbg"] = ColorStringToBGColor(conf.Settings.NoHighlightBG)
	}

	if len(conf.Settings.NoHighlightFG) > 0 {
		colors[color.Level16]["nohlfg"] = ColorStringToColor(conf.Settings.NoHighlightFG)
		colors[color.Level256]["nohlfg"] = ColorStringToColor(conf.Settings.NoHighlightFG)
		colors[color.LevelRgb]["nohlfg"] = ColorStringToColor(conf.Settings.NoHighlightFG)
	}

	if len(conf.Settings.HighlightHdrBG) > 0 {
		colors[color.Level16]["hdrbg"] = ColorStringToBGColor(conf.Settings.HighlightHdrBG)
		colors[color.Level256]["hdrbg"] = ColorStringToBGColor(conf.Settings.HighlightHdrBG)
		colors[color.LevelRgb]["hdrbg"] = ColorStringToBGColor(conf.Settings.HighlightHdrBG)
	}

	if len(conf.Settings.HighlightHdrFG) > 0 {
		colors[color.Level16]["hdrfg"] = ColorStringToColor(conf.Settings.HighlightHdrFG)
		colors[color.Level256]["hdrfg"] = ColorStringToColor(conf.Settings.HighlightHdrFG)
		colors[color.LevelRgb]["hdrfg"] = ColorStringToColor(conf.Settings.HighlightHdrFG)
	}

	return colors
}

// find supported color mode, modifies config based on constants
func (conf *Config) DetermineColormode() {
	if !isTerminal(os.Stdout) {
		color.Disable()
	} else {
		level := color.TermColorLevel()
		colors := conf.Colors()

		conf.ColorStyle = color.New(colors[level]["bg"], colors[level]["fg"])
		conf.HighlightStyle = color.New(colors[level]["hlbg"], colors[level]["hlfg"])
		conf.NoHighlightStyle = color.New(colors[level]["nohlbg"], colors[level]["nohlfg"])
		conf.HighlightHdrStyle = color.New(colors[level]["hdrbg"], colors[level]["hdrfg"])
	}
}

// Return true if current terminal is interactive
func isTerminal(f *os.File) bool {
	o, _ := f.Stat()

	return (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice
}

// main program version
// generated  version string, used  by -v contains  lib.Version on
//
//	main branch, and lib.Version-$branch-$lastcommit-$date on
//
// development branch
func Getversion() string {
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
		conf.OutputMode = ASCII
	}
}

func (conf *Config) PrepareFilters() error {
	conf.Filters = make(map[string]*regexp.Regexp, len(conf.Rawfilters))

	for _, filter := range conf.Rawfilters {
		parts := strings.Split(filter, "=")
		if len(parts) != MAXPARTS {
			return errors.New("filter field and value must be separated by =")
		}

		reg, err := regexp.Compile(parts[1])
		if err != nil {
			return fmt.Errorf("failed to compile filter regex for field %s: %w",
				parts[0], err)
		}

		conf.Filters[strings.ToLower(parts[0])] = reg
	}

	return nil
}

func (conf *Config) CheckEnv() {
	// check for environment vars, command line flags have precedence,
	// NO_COLOR is being checked by the color module itself.
	if !conf.NoNumbering {
		_, set := os.LookupEnv("T_NO_HEADER_NUMBERING")
		if set {
			conf.NoNumbering = true
		}
	}

	if len(conf.Columns) == 0 {
		cols := os.Getenv("T_COLUMNS")
		if len(cols) > 1 {
			conf.Columns = cols
		}
	}
}

func (conf *Config) ApplyDefaults() {
	// mode specific defaults
	if conf.OutputMode == Yaml || conf.OutputMode == CSV {
		conf.NoNumbering = true
	}

	ValidHooks = []string{"filter", "process", "transpose", "append"}
}

func (conf *Config) PreparePattern(pattern string) error {
	PatternR, err := regexp.Compile(pattern)

	if err != nil {
		return fmt.Errorf("regexp pattern %s is invalid: %w", conf.Pattern, err)
	}

	conf.PatternR = PatternR
	conf.Pattern = pattern

	return nil
}

// Parse config file.  Ignore if the file doesn't exist  but return an
// error if it exists but fails to read or parse
func (conf *Config) ParseConfigfile() error {
	path, err := os.Stat(conf.Configfile)

	if err != nil {
		if os.IsNotExist(err) {
			// ignore non-existent files
			return nil
		}

		return fmt.Errorf("failed to stat config file: %w", err)
	}

	if path.IsDir() {
		// ignore non-existent or dirs
		return nil
	}

	configstring, err := os.ReadFile(path.Name())
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", path.Name(), err)
	}

	err = hclsimple.Decode(
		path.Name(),
		configstring,
		nil,
		&conf.Settings)
	if err != nil {
		return fmt.Errorf("failed to load configuration file %s: %w",
			path.Name(), err)
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
