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
	"fmt"
	"strings"
)

func printData(data Tabdata) {
	if XtendedOut {
		printExtendedData(data)
	} else {
		printTabularData(data)
	}
}

func printTabularData(data Tabdata) {
	// needed for data output
	var formats []string

	if len(data.entries) > 0 {
		// headers
		for i, head := range data.headers {
			if len(Columns) > 0 {
				if !contains(UseColumns, i+1) {
					continue
				}
			}

			// calculate column width
			var width int
			var iwidth int
			var format string

			// generate format string
			if len(head) > data.maxwidthPerCol[i] {
				width = len(head)
			} else {
				width = data.maxwidthPerCol[i]
			}

			if NoNumbering {
				iwidth = 0
			} else {
				iwidth = len(fmt.Sprintf("%d", i)) // in case i > 9
			}

			format = fmt.Sprintf("%%-%ds", 3+iwidth+width)

			if NoNumbering {
				fmt.Printf(format, fmt.Sprintf("%s ", head))
			} else {
				fmt.Printf(format, fmt.Sprintf("%s(%d) ", head, i+1))
			}

			// register
			formats = append(formats, format)
		}
		fmt.Println()

		// entries
		var idx int
		for _, entry := range data.entries {
			idx = 0
			//fmt.Println(entry)
			for i, value := range entry {
				if len(Columns) > 0 {
					if !contains(UseColumns, i+1) {
						continue
					}
				}
				fmt.Printf(formats[idx], strings.TrimSpace(value))
				idx++
			}
			fmt.Println()
		}
	}
}

/*
   We simulate the \x command of psql (the PostgreSQL client)
*/
func printExtendedData(data Tabdata) {
	// needed for data output
	format := fmt.Sprintf("%%%ds: %%s\n", data.maxwidthHeader) // FIXME: re-calculate if -c has been set

	if len(data.entries) > 0 {
		var idx int
		for _, entry := range data.entries {
			idx = 0
			for i, value := range entry {
				if len(Columns) > 0 {
					if !contains(UseColumns, i+1) {
						continue
					}
				}

				fmt.Printf(format, data.headers[idx], value)
				idx++
			}
			fmt.Println()
		}
	}
}
