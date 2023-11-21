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
func PrepareColumns(c *cfg.Config, data *Tabdata) error {
	if len(c.Columns) > 0 {
		for _, use := range strings.Split(c.Columns, ",") {
			if len(use) == 0 {
				msg := fmt.Sprintf("Could not parse columns list %s: empty column", c.Columns)
				return errors.New(msg)
			}

			usenum, err := strconv.Atoi(use)
			if err != nil {
				// might be a regexp
				colPattern, err := regexp.Compile(use)
				if err != nil {
					msg := fmt.Sprintf("Could not parse columns list %s: %v", c.Columns, err)
					return errors.New(msg)
				}

				// find matching header fields
				for i, head := range data.headers {
					if colPattern.MatchString(head) {
						c.UseColumns = append(c.UseColumns, i+1)
					}

				}
			} else {
				// we digress from go  best practises here, because if
				// a colum spec is not a number, we process them above
				// inside the err handler  for atoi(). so only add the
				// number, if it's really just a number.
				c.UseColumns = append(c.UseColumns, usenum)
			}
		}

		// deduplicate: put all values into a map (value gets map key)
		// thereby  removing duplicates,  extract keys into  new slice
		// and sort it
		imap := make(map[int]int, len(c.UseColumns))
		for _, i := range c.UseColumns {
			imap[i] = 0
		}
		c.UseColumns = nil
		for k := range imap {
			c.UseColumns = append(c.UseColumns, k)
		}
		sort.Ints(c.UseColumns)
	}
	return nil
}

// prepare headers: add numbers to headers
func numberizeAndReduceHeaders(c cfg.Config, data *Tabdata) {
	numberedHeaders := []string{}
	maxwidth := 0 // start from scratch, so we only look at displayed column widths

	for i, head := range data.headers {
		headlen := 0
		if len(c.Columns) > 0 {
			// -c specified
			if !contains(c.UseColumns, i+1) {
				// ignore this one
				continue
			}
		}
		if c.NoNumbering {
			numberedHeaders = append(numberedHeaders, head)
			headlen = len(head)
		} else {
			numhead := fmt.Sprintf("%s(%d)", head, i+1)
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
func reduceColumns(c cfg.Config, data *Tabdata) {
	if len(c.Columns) > 0 {
		reducedEntries := [][]string{}
		var reducedEntry []string
		for _, entry := range data.entries {
			reducedEntry = nil
			for i, value := range entry {
				if !contains(c.UseColumns, i+1) {
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

func colorizeData(c cfg.Config, output string) string {
	if len(c.Pattern) > 0 && !c.NoColor && color.IsConsole(os.Stdout) {
		r := regexp.MustCompile("(" + c.Pattern + ")")
		highlight := true
		colorized := ""

		for _, line := range strings.Split(output, "\n") {
			if c.UseHighlight {
				if highlight {
					line = c.HighlightStyle.Sprint(line)
				}
				highlight = !highlight
			} else {
				line = r.ReplaceAllStringFunc(line, func(in string) string {
					return c.ColorStyle.Sprint(in)
				})
			}

			colorized += line + "\n"
		}

		return colorized
	} else {
		return output
	}
}
