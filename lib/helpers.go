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

func PrepareModeFlags() error {
	if len(OutputMode) == 0 {
		switch {
		case OutflagExtended:
			OutputMode = "extended"
		case OutflagMarkdown:
			OutputMode = "markdown"
		case OutflagOrgtable:
			OutputMode = "orgtbl"
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
