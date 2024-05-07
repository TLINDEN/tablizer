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
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/tlinden/tablizer/cfg"
)

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// parse columns list given  with -c, modifies config.UseColumns based
// on eventually given regex
func PrepareColumns(conf *cfg.Config, data *Tabdata) error {
	if len(conf.Columns) > 0 {
		for _, use := range strings.Split(conf.Columns, ",") {
			if len(use) == 0 {
				msg := fmt.Sprintf("Could not parse columns list %s: empty column", conf.Columns)
				return errors.New(msg)
			}

			usenum, err := strconv.Atoi(use)
			if err != nil {
				// might be a regexp
				colPattern, err := regexp.Compile(use)
				if err != nil {
					msg := fmt.Sprintf("Could not parse columns list %s: %v", conf.Columns, err)
					return errors.New(msg)
				}

				// find matching header fields
				for i, head := range data.headers {
					if colPattern.MatchString(head) {
						conf.UseColumns = append(conf.UseColumns, i+1)
					}

				}
			} else {
				// we digress from go  best practises here, because if
				// a colum spec is not a number, we process them above
				// inside the err handler  for atoi(). so only add the
				// number, if it's really just a number.
				conf.UseColumns = append(conf.UseColumns, usenum)
			}
		}

		// deduplicate: put all values into a map (value gets map key)
		// thereby  removing duplicates,  extract keys into  new slice
		// and sort it
		imap := make(map[int]int, len(conf.UseColumns))
		for _, i := range conf.UseColumns {
			imap[i] = 0
		}
		conf.UseColumns = nil
		for k := range imap {
			conf.UseColumns = append(conf.UseColumns, k)
		}
		sort.Ints(conf.UseColumns)
	}
	return nil
}

// prepare headers: add numbers to headers
func numberizeAndReduceHeaders(conf cfg.Config, data *Tabdata) {
	numberedHeaders := []string{}
	maxwidth := 0 // start from scratch, so we only look at displayed column widths

	for idx, head := range data.headers {
		headlen := 0
		if len(conf.Columns) > 0 {
			// -c specified
			if !contains(conf.UseColumns, idx+1) {
				// ignore this one
				continue
			}
		}
		if conf.NoNumbering {
			numberedHeaders = append(numberedHeaders, head)
			headlen = len(head)
		} else {
			numhead := fmt.Sprintf("%s(%d)", head, idx+1)
			headlen = len(numhead)
			numberedHeaders = append(numberedHeaders, numhead)
		}

		if headlen > maxwidth {
			maxwidth = headlen
		}
	}
	data.headers = numberedHeaders
	if data.maxwidthHeader != maxwidth && maxwidth > 0 {
		data.maxwidthHeader = maxwidth
	}
}

// exclude columns, if any
func reduceColumns(conf cfg.Config, data *Tabdata) {
	if len(conf.Columns) > 0 {
		reducedEntries := [][]string{}
		var reducedEntry []string
		for _, entry := range data.entries {
			reducedEntry = nil
			for i, value := range entry {
				if !contains(conf.UseColumns, i+1) {
					continue
				}

				reducedEntry = append(reducedEntry, value)
			}
			reducedEntries = append(reducedEntries, reducedEntry)
		}
		data.entries = reducedEntries
	}
}

func trimRow(row []string) []string {
	// FIXME: remove this when we only use Tablewriter and strip in ParseFile()!
	var fixedrow []string
	for _, cell := range row {
		fixedrow = append(fixedrow, strings.TrimSpace(cell))
	}

	return fixedrow
}

func colorizeData(conf cfg.Config, output string) string {
	switch {
	case conf.UseHighlight && color.IsConsole(os.Stdout):
		highlight := true
		colorized := ""
		first := true

		for _, line := range strings.Split(output, "\n") {
			if highlight {
				if first {
					// we  need to add  two spaces to the  header line
					//  because tablewriter omits them for some reason
					//  in pprint mode. This doesn't matter as long as
					//  we don't use colorization. But with colors the
					// missing spaces can be seen.
					if conf.OutputMode == cfg.ASCII {
						line += "  "
					}

					line = conf.HighlightHdrStyle.Sprint(line)
					first = false
				} else {
					line = conf.HighlightStyle.Sprint(line)
				}
			} else {
				line = conf.NoHighlightStyle.Sprint(line)
			}
			highlight = !highlight

			colorized += line + "\n"
		}

		return colorized
	case len(conf.Pattern) > 0 && !conf.NoColor && color.IsConsole(os.Stdout):
		r := regexp.MustCompile("(" + conf.Pattern + ")")
		return r.ReplaceAllStringFunc(output, func(in string) string {
			return conf.ColorStyle.Sprint(in)
		})
	default:
		return output
	}
}
