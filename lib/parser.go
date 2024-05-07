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
	"errors"
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
func Parse(c cfg.Config, input io.Reader) (Tabdata, error) {
	if len(c.Separator) == 1 {
		return parseCSV(c, input)
	}

	return parseTabular(c, input)
}

/*
Parse CSV input.
*/
func parseCSV(c cfg.Config, input io.Reader) (Tabdata, error) {
	var content io.Reader = input
	data := Tabdata{}

	if len(c.Pattern) > 0 {
		scanner := bufio.NewScanner(input)
		lines := []string{}
		hadFirst := false
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if hadFirst {
				// don't match 1st line, it's the header
				if c.Pattern != "" && matchPattern(c, line) == c.InvertMatch {
					// by default  -v is false, so if a  line does NOT
					// match the pattern, we will ignore it. However,
					// if the user specified -v, the matching is inverted,
					// so we ignore all lines, which DO match.
					continue
				}

				// apply user defined lisp filters, if any
				accept, err := RunFilterHooks(c, line)
				if err != nil {
					return data, errors.Unwrap(fmt.Errorf("Failed to apply filter hook: %w", err))
				}

				if !accept {
					//  IF there  are filter  hook[s] and  IF one  of them
					// returns false on the current line, reject it
					continue
				}
			}
			lines = append(lines, line)
			hadFirst = true
		}
		content = strings.NewReader(strings.Join(lines, "\n"))
	}

	csvreader := csv.NewReader(content)
	csvreader.Comma = rune(c.Separator[0])

	records, err := csvreader.ReadAll()
	if err != nil {
		return data, errors.Unwrap(fmt.Errorf("Could not parse CSV input: %w", err))
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
	userdata, changed, err := RunProcessHooks(c, data)
	if err != nil {
		return data, errors.Unwrap(fmt.Errorf("Failed to apply filter hook: %w", err))
	}
	if changed {
		data = userdata
	}

	return data, nil
}

/*
Parse tabular input.
*/
func parseTabular(c cfg.Config, input io.Reader) (Tabdata, error) {
	data := Tabdata{}

	var scanner *bufio.Scanner

	hadFirst := false
	separate := regexp.MustCompile(c.Separator)

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
				// if Debug {
				// 	fmt.Printf("Part: <%s>\n", string(line[beg:part[0]]))
				//}

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
			if c.Pattern != "" && matchPattern(c, line) == c.InvertMatch {
				// by default  -v is false, so if a  line does NOT
				// match the pattern, we will ignore it. However,
				// if the user specified -v, the matching is inverted,
				// so we ignore all lines, which DO match.
				continue
			}

			// apply user defined lisp filters, if any
			accept, err := RunFilterHooks(c, line)
			if err != nil {
				return data, errors.Unwrap(fmt.Errorf("Failed to apply filter hook: %w", err))
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
		return data, errors.Unwrap(fmt.Errorf("Failed to read from io.Reader: %w", scanner.Err()))
	}

	// filter by field filters, if any
	filtereddata, changed, err := FilterByFields(c, data)
	if err != nil {
		return data, fmt.Errorf("Failed to filter fields: %w", err)
	}
	if changed {
		data = filtereddata
	}

	// apply user defined lisp process hooks, if any
	userdata, changed, err := RunProcessHooks(c, data)
	if err != nil {
		return data, errors.Unwrap(fmt.Errorf("Failed to apply filter hook: %w", err))
	}
	if changed {
		data = userdata
	}

	if c.Debug {
		repr.Print(data)
	}

	return data, nil
}
