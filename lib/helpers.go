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
	"github.com/gookit/color"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func PrepareColumns() error {
	if len(Columns) > 0 {
		for _, use := range strings.Split(Columns, ",") {
			usenum, err := strconv.Atoi(use)
			if err != nil {
				msg := fmt.Sprintf("Could not parse columns list %s: %v", Columns, err)
				return errors.New(msg)
			}
			UseColumns = append(UseColumns, usenum)
		}
	}
	return nil
}

func numberizeHeaders(data *Tabdata) {
	// prepare headers: add numbers to headers
	numberedHeaders := []string{}
	maxwidth := 0 // start from scratch, so we only look at displayed column widths

	for i, head := range data.headers {
		headlen := 0
		if len(Columns) > 0 {
			// -c specified
			if !contains(UseColumns, i+1) {
				// ignore this one
				continue
			}
		}
		if NoNumbering {
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

func reduceColumns(data *Tabdata) {
	// exclude columns, if any
	if len(Columns) > 0 {
		reducedEntries := [][]string{}
		var reducedEntry []string
		for _, entry := range data.entries {
			reducedEntry = nil
			for i, value := range entry {
				if !contains(UseColumns, i+1) {
					continue
				}

				reducedEntry = append(reducedEntry, value)
			}
			reducedEntries = append(reducedEntries, reducedEntry)
		}
		data.entries = reducedEntries
	}
}

func PrepareModeFlags() error {
	if len(OutputMode) == 0 {
		// associate short flags like -X with mode selector
		switch {
		case OutflagExtended:
			OutputMode = "extended"
		case OutflagMarkdown:
			OutputMode = "markdown"
		case OutflagOrgtable:
			OutputMode = "orgtbl"
		case OutflagShell:
			OutputMode = "shell"
			NoNumbering = true
		default:
			OutputMode = "ascii"
		}
	} else {
		r, err := regexp.Compile(validOutputmodes)

		if err != nil {
			return errors.New("Failed to validate output mode spec!")
		}

		match := r.MatchString(OutputMode)

		if !match {
			return errors.New("Invalid output mode!")
		}
	}

	return nil
}

func trimRow(row []string) []string {
	// FIXME: remove this when we only use Tablewriter and strip in ParseFile()!
	var fixedrow []string
	for _, cell := range row {
		fixedrow = append(fixedrow, strings.TrimSpace(cell))
	}

	return fixedrow
}

func colorizeData(output string) string {
	if len(Pattern) > 0 && !NoColor && color.IsConsole(os.Stdout) {
		r := regexp.MustCompile("(" + Pattern + ")")
		return r.ReplaceAllString(output, "<bg="+MatchBG+";fg="+MatchFG+">$1</>")
	} else {
		return output
	}
}

func isTerminal(f *os.File) bool {
	o, _ := f.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		return true
	} else {
		return false
	}
}
