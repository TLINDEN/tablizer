/*
Copyright © 2022-2024 Thomas von Dein

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
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/alecthomas/repr"
	"github.com/tlinden/tablizer/cfg"
)

/*
Parser switch
*/
func Parse(conf cfg.Config, input io.Reader) (Tabdata, error) {
	if len(conf.Separator) == 1 {
		return parseCSV(conf, input)
	}

	return parseTabular(conf, input)
}

/*
Parse CSV input.
*/
func parseCSV(conf cfg.Config, input io.Reader) (Tabdata, error) {
	data := Tabdata{}

	// apply pattern, if any
	content, err := FilterByPattern(conf, input)
	if err != nil {
		return data, err
	}

	csvreader := csv.NewReader(content)
	csvreader.Comma = rune(conf.Separator[0])

	records, err := csvreader.ReadAll()
	if err != nil {
		return data, fmt.Errorf("could not parse CSV input: %w", err)
	}

	if len(records) >= 1 {
		data.headers = records[0]
		data.columns = len(records)

		for _, head := range data.headers {
			// register widest header field
			headerlen := len(head)
			if headerlen > data.maxwidthHeader {
				data.maxwidthHeader = headerlen
			}
		}

		if len(records) > 1 {
			data.entries = records[1:]
		}
	}

	// apply user defined lisp process hooks, if any
	userdata, changed, err := RunProcessHooks(conf, data)
	if err != nil {
		return data, fmt.Errorf("failed to apply filter hook: %w", err)
	}

	if changed {
		data = userdata
	}

	return data, nil
}

/*
Parse tabular input.
*/
func parseTabular(conf cfg.Config, input io.Reader) (Tabdata, error) {
	data := Tabdata{}

	var scanner *bufio.Scanner

	hadFirst := false
	separate := regexp.MustCompile(conf.Separator)

	scanner = bufio.NewScanner(input)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := separate.Split(line, -1)

		if !hadFirst {
			// header processing
			data.columns = len(parts)
			// if Debug {
			// 	fmt.Println(parts)
			// }

			// process all header fields
			for _, part := range parts {
				// register widest header field
				headerlen := len(part)
				if headerlen > data.maxwidthHeader {
					data.maxwidthHeader = headerlen
				}

				// register fields data
				data.headers = append(data.headers, strings.TrimSpace(part))

				// done
				hadFirst = true
			}
		} else {
			// data processing
			if conf.Pattern != "" && matchPattern(conf, line) == conf.InvertMatch {
				// by default  -v is false, so if a  line does NOT
				// match the pattern, we will ignore it. However,
				// if the user specified -v, the matching is inverted,
				// so we ignore all lines, which DO match.
				continue
			}

			// apply user defined lisp filters, if any
			accept, err := RunFilterHooks(conf, line)
			if err != nil {
				return data, fmt.Errorf("failed to apply filter hook: %w", err)
			}

			if !accept {
				//  IF there  are filter  hook[s] and  IF one  of them
				// returns false on the current line, reject it
				continue
			}

			idx := 0 // we cannot use the header index, because we could exclude columns
			values := []string{}
			for _, part := range parts {
				// if Debug {
				// 	fmt.Printf("<%s> ", value)
				// }
				values = append(values, strings.TrimSpace(part))

				idx++
			}

			// fill up missing fields, if any
			for i := len(values); i < len(data.headers); i++ {
				values = append(values, "")
			}

			data.entries = append(data.entries, values)
		}
	}

	if scanner.Err() != nil {
		return data, fmt.Errorf("failed to read from io.Reader: %w", scanner.Err())
	}

	// filter by field filters, if any
	filtereddata, changed, err := FilterByFields(conf, &data)
	if err != nil {
		return data, fmt.Errorf("failed to filter fields: %w", err)
	}

	if changed {
		data = *filtereddata
	}

	// transpose if demanded
	if err := PrepareTransposerColumns(&conf, &data); err != nil {
		return data, err
	}

	modifieddata, changed, err := TransposeFields(conf, &data)
	if err != nil {
		return data, fmt.Errorf("failed to transpose fields: %w", err)
	}

	if changed {
		data = *modifieddata
	}

	// apply user defined lisp process hooks, if any
	userdata, changed, err := RunProcessHooks(conf, data)
	if err != nil {
		return data, fmt.Errorf("failed to apply filter hook: %w", err)
	}

	if changed {
		data = userdata
	}

	if conf.Debug {
		repr.Print(data)
	}

	return data, nil
}
