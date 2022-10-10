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
	"github.com/gookit/color"
	"github.com/olekukonko/tablewriter"
	"regexp"
	"strings"
)

func printData(data *Tabdata) {
	if OutputMode != "shell" {
		numberizeHeaders(data)
	}
	reduceColumns(data)

	switch OutputMode {
	case "extended":
		printExtendedData(data)
	case "ascii":
		printAsciiData(data)
	case "orgtbl":
		printOrgmodeData(data)
	case "markdown":
		printMarkdownData(data)
	case "shell":
		printShellData(data)
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

	color.Print(
		colorizeData(
			rightR.ReplaceAllString(
				leftR.ReplaceAllString(tableString.String(), "|"), "|")))
}

/*
   Markdown table
*/
func printMarkdownData(data *Tabdata) {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetHeader(data.headers)

	for _, row := range data.entries {
		table.Append(trimRow(row))
	}

	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.Render()
	color.Print(colorizeData(tableString.String()))
}

/*
   Simple ASCII table without any borders etc, just like the input we expect
*/
func printAsciiData(data *Tabdata) {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

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
	color.Print(colorizeData(tableString.String()))
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

				color.Printf(format, data.headers[idx], value)
				idx++
			}
			fmt.Println()
		}
	}
}

/*
   Shell output, ready to be eval'd. Just like FreeBSD stat(1)
*/
func printShellData(data *Tabdata) {
	if len(data.entries) > 0 {
		var idx int
		for _, entry := range data.entries {
			idx = 0
			shentries := []string{}
			for i, value := range entry {
				if len(Columns) > 0 {
					if !contains(UseColumns, i+1) {
						continue
					}
				}

				shentries = append(shentries, fmt.Sprintf("%s=\"%s\"", data.headers[idx], value))
				idx++
			}
			fmt.Println(strings.Join(shentries, " "))
		}
	}
}
