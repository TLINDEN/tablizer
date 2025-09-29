/*
Copyright © 2022-2025 Thomas von Dein

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
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/tlinden/tablizer/cfg"
	"gopkg.in/yaml.v3"
)

func printData(writer io.Writer, conf cfg.Config, data *Tabdata) {
	//  Sort   the  data  first,  before   headers+entries  are  being
	// reduced. That way the user can specify any valid column to sort
	// by, independently if it's being used for display or not.
	sortTable(conf, data)

	// put one or more columns into clipboard
	yankColumns(conf, data)

	// add numbers to headers and remove those we're not interested in
	numberizeAndReduceHeaders(conf, data)

	// remove unwanted columns, if any
	reduceColumns(conf, data)

	switch conf.OutputMode {
	case cfg.Extended:
		printExtendedData(writer, conf, data)
	case cfg.ASCII:
		printASCIIData(writer, conf, data)
	case cfg.Orgtbl:
		printOrgmodeData(writer, conf, data)
	case cfg.Markdown:
		printMarkdownData(writer, conf, data)
	case cfg.Shell:
		printShellData(writer, data)
	case cfg.Yaml:
		printYamlData(writer, data)
	case cfg.CSV:
		printCSVData(writer, data)
	default:
		printASCIIData(writer, conf, data)
	}
}

func output(writer io.Writer, str string) {
	fmt.Fprint(writer, str)
}

/*
Emacs org-mode compatible table (also orgtbl-mode)
*/
func printOrgmodeData(writer io.Writer, conf cfg.Config, data *Tabdata) {
	tableString := &strings.Builder{}

	table := tablewriter.NewTable(tableString,
		tablewriter.WithRenderer(
			renderer.NewBlueprint(
				tw.Rendition{
					Borders: tw.Border{
						Left:   tw.On,
						Right:  tw.On,
						Top:    tw.On,
						Bottom: tw.On,
					},
					Settings: tw.Settings{
						Separators: tw.Separators{
							ShowHeader:     tw.On,
							ShowFooter:     tw.Off,
							BetweenRows:    tw.Off,
							BetweenColumns: 0,
						},
					},
					Symbols: tw.NewSymbols(tw.StyleASCII),
				})),

		tablewriter.WithConfig(
			tablewriter.Config{
				Header: tw.CellConfig{
					Formatting: tw.CellFormatting{
						Alignment:  tw.AlignLeft,
						AutoFormat: tw.Off,
					},
				},
				Row: tw.CellConfig{
					Formatting: tw.CellFormatting{
						Alignment: tw.AlignLeft,
					},
				},
			},
		),
	)

	if !conf.NoHeaders {
		table.Header(data.headers)
	}

	if err := table.Bulk(data.entries); err != nil {
		log.Fatalf("Failed to add data to table renderer: %s", err)
	}

	if err := table.Render(); err != nil {
		log.Fatalf("Failed to render table: %s", err)
	}

	output(writer, color.Sprint(colorizeData(conf, tableString.String())))
}

/*
Markdown table
*/
func printMarkdownData(writer io.Writer, conf cfg.Config, data *Tabdata) {
	tableString := &strings.Builder{}

	table := tablewriter.NewTable(tableString,
		tablewriter.WithRenderer(
			renderer.NewBlueprint(
				tw.Rendition{
					Borders: tw.Border{
						Left:   tw.On,
						Right:  tw.On,
						Top:    tw.Off,
						Bottom: tw.Off,
					},
					Settings: tw.Settings{
						Separators: tw.Separators{
							ShowHeader:     tw.On,
							ShowFooter:     tw.Off,
							BetweenRows:    tw.Off,
							BetweenColumns: 0,
						},
					},
					Symbols: tw.NewSymbols(tw.StyleMarkdown),
				})),

		tablewriter.WithConfig(
			tablewriter.Config{
				Header: tw.CellConfig{
					Formatting: tw.CellFormatting{
						Alignment:  tw.AlignLeft,
						AutoFormat: tw.Off,
					},
				},
				Row: tw.CellConfig{
					Formatting: tw.CellFormatting{
						Alignment: tw.AlignLeft,
					},
				},
			},
		),
	)

	if !conf.NoHeaders {
		table.Header(data.headers)
	}

	if err := table.Bulk(data.entries); err != nil {
		log.Fatalf("Failed to add data to table renderer: %s", err)
	}

	if err := table.Render(); err != nil {
		log.Fatalf("Failed to render table: %s", err)
	}

	output(writer, color.Sprint(colorizeData(conf, tableString.String())))
}

/*
Simple ASCII table without any borders etc, just like the input we expect
*/
func printASCIIData(writer io.Writer, conf cfg.Config, data *Tabdata) {
	tableString := &strings.Builder{}

	styleTSV := tw.NewSymbolCustom("space").WithColumn("\t")

	table := tablewriter.NewTable(tableString,
		tablewriter.WithRenderer(
			renderer.NewBlueprint(tw.Rendition{
				Borders: tw.BorderNone,
				Symbols: styleTSV,
				Settings: tw.Settings{
					Separators: tw.SeparatorsNone,
					Lines:      tw.LinesNone,
				},
			})),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{
					AutoFormat: tw.Off,
				},
				Padding: tw.CellPadding{Global: tw.Padding{Left: "", Right: " "}},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{
					AutoWrap:  tw.WrapNone,
					Alignment: tw.AlignLeft,
				},
				Padding: tw.CellPadding{Global: tw.Padding{Right: " "}},
			},

			Debug: true,
		}),
	)

	if !conf.NoHeaders {
		table.Header(data.headers)
	}

	if err := table.Bulk(data.entries); err != nil {
		log.Fatalf("Failed to add data to table renderer: %s", err)
	}

	if err := table.Render(); err != nil {
		log.Fatalf("Failed to render table: %s", err)
	}

	output(writer, color.Sprint(colorizeData(conf, tableString.String())))
}

/*
We simulate the \x command of psql (the PostgreSQL client)
*/
func printExtendedData(writer io.Writer, conf cfg.Config, data *Tabdata) {
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

	output(writer, colorizeData(conf, out))
}

/*
Shell output, ready to be eval'd. Just like FreeBSD stat(1)
*/
func printShellData(writer io.Writer, data *Tabdata) {
	out := ""

	if len(data.entries) > 0 {
		for _, entry := range data.entries {
			shentries := []string{}

			for idx, value := range entry {
				shentries = append(shentries, fmt.Sprintf("%s=\"%s\"",
					data.headers[idx], value))
			}

			out += strings.Join(shentries, " ") + "\n"
		}
	}

	// no colorization here
	output(writer, out)
}

func printYamlData(writer io.Writer, data *Tabdata) {
	type Data struct {
		Entries []map[string]interface{} `yaml:"entries"`
	}

	yamlout := Data{}

	for _, entry := range data.entries {
		yamldata := map[string]interface{}{}

		for idx, entry := range entry {
			style := yaml.TaggedStyle

			_, err := strconv.Atoi(entry)
			if err != nil {
				style = yaml.DoubleQuotedStyle
			}

			yamldata[strings.ToLower(data.headers[idx])] =
				&yaml.Node{
					Kind:  yaml.ScalarNode,
					Style: style,
					Value: entry}
		}

		yamlout.Entries = append(yamlout.Entries, yamldata)
	}

	yamlstr, err := yaml.Marshal(&yamlout)

	if err != nil {
		log.Fatal(err)
	}

	output(writer, string(yamlstr))
}

func printCSVData(writer io.Writer, data *Tabdata) {
	csvout := csv.NewWriter(writer)

	if err := csvout.Write(data.headers); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for _, entry := range data.entries {
		if err := csvout.Write(entry); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	csvout.Flush()

	if err := csvout.Error(); err != nil {
		log.Fatal(err)
	}
}
