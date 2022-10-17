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
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func printData(w io.Writer, data *Tabdata) {
	// some output preparations:

	// add numbers to headers and remove this we're not interested in
	numberizeAndReduceHeaders(data)

	// remove unwanted columns, if any
	reduceColumns(data)

	// sort the data
	sortTable(data, SortByColumn)

	switch OutputMode {
	case "extended":
		printExtendedData(w, data)
	case "ascii":
		printAsciiData(w, data)
	case "orgtbl":
		printOrgmodeData(w, data)
	case "markdown":
		printMarkdownData(w, data)
	case "shell":
		printShellData(w, data)
	case "yaml":
		printYamlData(w, data)
	default:
		printAsciiData(w, data)
	}
}

func output(w io.Writer, str string) {
	fmt.Fprint(w, str)
}

/*
   Emacs org-mode compatible table (also orgtbl-mode)
*/
func printOrgmodeData(w io.Writer, data *Tabdata) {
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

	output(w, color.Sprint(
		colorizeData(
			rightR.ReplaceAllString(
				leftR.ReplaceAllString(tableString.String(), "|"), "|"))))
}

/*
   Markdown table
*/
func printMarkdownData(w io.Writer, data *Tabdata) {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetHeader(data.headers)

	for _, row := range data.entries {
		table.Append(trimRow(row))
	}

	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.Render()
	output(w, color.Sprint(colorizeData(tableString.String())))
}

/*
   Simple ASCII table without any borders etc, just like the input we expect
*/
func printAsciiData(w io.Writer, data *Tabdata) {
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
	output(w, color.Sprint(colorizeData(tableString.String())))
}

/*
   We simulate the \x command of psql (the PostgreSQL client)
*/
func printExtendedData(w io.Writer, data *Tabdata) {
	// needed for data output
	format := fmt.Sprintf("%%%ds: %%s\n", data.maxwidthHeader)
	out := ""
	if len(data.entries) > 0 {
		for _, entry := range data.entries {
			for i, value := range entry {
				out += color.Sprintf(format, data.headers[i], value)
			}

			out += "\n"
		}
	}

	output(w, out)
}

/*
   Shell output, ready to be eval'd. Just like FreeBSD stat(1)
*/
func printShellData(w io.Writer, data *Tabdata) {
	out := ""
	if len(data.entries) > 0 {
		for _, entry := range data.entries {
			shentries := []string{}
			for i, value := range entry {
				shentries = append(shentries, fmt.Sprintf("%s=\"%s\"",
					data.headers[i], value))
			}
			out += fmt.Sprint(strings.Join(shentries, " ")) + "\n"
		}
	}
	output(w, out)
}

func printYamlData(w io.Writer, data *Tabdata) {
	type D struct {
		Entries []map[string]interface{} `yaml:"entries"`
	}

	d := D{}

	for _, entry := range data.entries {
		ml := map[string]interface{}{}

		for i, entry := range entry {
			style := yaml.TaggedStyle
			_, err := strconv.Atoi(entry)
			if err != nil {
				style = yaml.DoubleQuotedStyle
			}

			ml[strings.ToLower(data.headers[i])] =
				&yaml.Node{
					Kind:  yaml.ScalarNode,
					Style: style,
					Value: entry}
		}

		d.Entries = append(d.Entries, ml)
	}

	yamlstr, err := yaml.Marshal(&d)

	if err != nil {
		log.Fatal(err)
	}

	output(w, string(yamlstr))
}
