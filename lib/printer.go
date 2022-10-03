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
	"github.com/olekukonko/tablewriter"
	"os"
	"regexp"
	"strings"
)

func printData(data *Tabdata) {
	// prepare headers: add numbers to headers
	if !NoNumbering {
		numberedHeaders := []string{}
		for i, head := range data.headers {
			if len(Columns) > 0 {
				// -c specified
				if !contains(UseColumns, i+1) {
					// ignore this one
					continue
				}
			}
			numberedHeaders = append(numberedHeaders, fmt.Sprintf("%s(%d)", head, i+1))
		}
		data.headers = numberedHeaders
	}

	// prepare data
	if len(Columns) > 0 {
		reducedEntries := [][]string{}
		reducedEntry := []string{}
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

	switch OutputMode {
	case "extended":
		printExtendedData(data)
	case "ascii":
		printAsciiData(data)
	case "orgtbl":
		printOrgmodeData(data)
	case "markdown":
		printMarkdownData(data)
	default:
		printAsciiData(data)
	}
}

/*
   Emacs org-mode compatible table (also orgtbl-mode)
*/
func printOrgmodeData(data *Tabdata) {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetHeader(data.headers)

	for _, row := range data.entries {
		table.Append(trimRow(row))
	}

	table.Render()

	/* fix output for org-mode (orgtbl)
	   tableWriter output:
	   +------+------+
	   | cell | cell |
	   +------+------+

	   Needed for org-mode compatibility:
	   |------+------|
	   | cell | cell |
	   |------+------|
	*/
	leftR := regexp.MustCompile("(?m)^\\+")
	rightR := regexp.MustCompile("\\+(?m)$")

	fmt.Print(rightR.ReplaceAllString(leftR.ReplaceAllString(tableString.String(), "|"), "|"))
}

/*
   Markdown table
*/
func printMarkdownData(data *Tabdata) {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(data.headers)

	for _, row := range data.entries {
		table.Append(trimRow(row))
	}

	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.Render()
}

/*
   Simple ASCII table without any borders etc, just like the input we expect
*/
func printAsciiData(data *Tabdata) {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(data.headers)
	table.AppendBulk(data.entries)

	// for _, row := range data.entries {
	// 	table.Append(trimRow(row))
	// }

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	table.Render()
}

/*
   We simulate the \x command of psql (the PostgreSQL client)
*/
func printExtendedData(data *Tabdata) {
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
