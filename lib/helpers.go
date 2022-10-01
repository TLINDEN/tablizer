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
	"strconv"
	"strings"
)

func die(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func PrepareColumns() {
	if len(Columns) > 0 {
		for _, use := range strings.Split(Columns, ",") {
			usenum, err := strconv.Atoi(use)
			if err != nil {
				die(err)
			}
			UseColumns = append(UseColumns, usenum)
		}
	}
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
