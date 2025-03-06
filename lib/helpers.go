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

func findindex(s []int, e int) (int, bool) {
	for i, a := range s {
		if a == e {
			return i, true
		}
	}

	return 0, false
}

// validate the consitency of parsed data
func ValidateConsistency(data *Tabdata) error {
	expectedfields := len(data.headers)

	for idx, row := range data.entries {
		if len(row) != expectedfields {
			return fmt.Errorf("row %d does not contain expected %d elements, but %d",
				idx, expectedfields, len(row))
		}
	}

	return nil
}

// parse columns list given  with -c, modifies config.UseColumns based
// on eventually given regex.
// This is an output filter, because -cN,N,... is being applied AFTER
// processing of the input data.
func PrepareColumns(conf *cfg.Config, data *Tabdata) error {
	// -c columns
	usecolumns, err := PrepareColumnVars(conf.Columns, data)
	if err != nil {
		return err
	}

	conf.UseColumns = usecolumns

	// -y columns
	useyankcolumns, err := PrepareColumnVars(conf.YankColumns, data)
	if err != nil {
		return err
	}

	conf.UseYankColumns = useyankcolumns

	return nil
}

// Same thing as above but for -T option, which is an input option,
// because transposers are being applied before output.
func PrepareTransposerColumns(conf *cfg.Config, data *Tabdata) error {
	// -T columns
	usetransposecolumns, err := PrepareColumnVars(conf.TransposeColumns, data)
	if err != nil {
		return err
	}

	conf.UseTransposeColumns = usetransposecolumns

	// verify that columns and transposers match and prepare transposer structs
	if err := conf.PrepareTransposers(); err != nil {
		return err
	}

	return nil
}

// output option, prepare -k1,2 sort fields
func PrepareSortColumns(conf *cfg.Config, data *Tabdata) error {
	// -c columns
	usecolumns, err := PrepareColumnVars(conf.SortByColumn, data)
	if err != nil {
		return err
	}

	conf.UseSortByColumn = usecolumns

	return nil
}

func PrepareColumnVars(columns string, data *Tabdata) ([]int, error) {
	if columns == "" {
		return nil, nil
	}

	usecolumns := []int{}

	isregex := regexp.MustCompile(`\W`)

	for _, columnpattern := range strings.Split(columns, ",") {
		if len(columnpattern) == 0 {
			return nil, fmt.Errorf("could not parse columns list %s: empty column", columns)
		}

		usenum, err := strconv.Atoi(columnpattern)
		if err != nil {
			// not a number

			if !isregex.MatchString(columnpattern) {
				// is not a regexp (contains no non-word chars)
				// lc() it so that word searches are case insensitive
				columnpattern = strings.ToLower(columnpattern)

				for i, head := range data.headers {
					if columnpattern == strings.ToLower(head) {
						usecolumns = append(usecolumns, i+1)
					}
				}
			} else {
				colPattern, err := regexp.Compile("(?i)" + columnpattern)
				if err != nil {
					msg := fmt.Sprintf("Could not parse columns list %s: %v", columns, err)

					return nil, errors.New(msg)
				}

				// find matching header fields, ignoring case
				for i, head := range data.headers {
					if colPattern.MatchString(strings.ToLower(head)) {
						usecolumns = append(usecolumns, i+1)
					}
				}
			}
		} else {
			// we digress from go  best practises here, because if
			// a colum spec is not a number, we process them above
			// inside the err handler  for atoi(). so only add the
			// number, if it's really just a number.
			usecolumns = append(usecolumns, usenum)
		}
	}

	// deduplicate: put all values into a map (value gets map key)
	// thereby  removing duplicates,  extract keys into  new slice
	// and sort it
	imap := make(map[int]int, len(usecolumns))
	for _, i := range usecolumns {
		imap[i] = 0
	}

	// fill with deduplicated columns
	usecolumns = nil

	for k := range imap {
		usecolumns = append(usecolumns, k)
	}

	sort.Ints(usecolumns)

	return usecolumns, nil
}

// prepare headers: add numbers to headers
func numberizeAndReduceHeaders(conf cfg.Config, data *Tabdata) {
	numberedHeaders := []string{}
	maxwidth := 0 // start from scratch, so we only look at displayed column widths

	for idx, head := range data.headers {
		var headlen int

		if len(conf.Columns) > 0 {
			// -c specified
			if !contains(conf.UseColumns, idx+1) {
				// ignore this one
				continue
			}
		}

		if conf.Numbering {
			numhead := fmt.Sprintf("%s(%d)", head, idx+1)
			headlen = len(numhead)
			numberedHeaders = append(numberedHeaders, numhead)
		} else {
			numberedHeaders = append(numberedHeaders, head)
			headlen = len(head)
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

// FIXME: remove this when we only use Tablewriter and strip in ParseFile()!
func trimRow(row []string) []string {
	var fixedrow = make([]string, len(row))

	for idx, cell := range row {
		fixedrow[idx] = strings.TrimSpace(cell)
	}

	return fixedrow
}

// FIXME: refactor this beast!
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

	case len(conf.Patterns) > 0 && !conf.NoColor && color.IsConsole(os.Stdout):
		out := output

		for _, re := range conf.Patterns {
			if !re.Negate {
				r := regexp.MustCompile("(" + re.Pattern + ")")

				out = r.ReplaceAllStringFunc(out, func(in string) string {
					return conf.ColorStyle.Sprint(in)
				})
			}
		}

		return out

	default:
		return output
	}
}
