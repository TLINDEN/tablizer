package cmd

import (
	"fmt"
	"strings"
)

func printTable(data Tabdata) {
	if XtendedOut {
		printExtended(data)
		return
	}

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
func printExtended(data Tabdata) {
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

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
