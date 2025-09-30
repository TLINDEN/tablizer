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
	"encoding/json"
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
	var data Tabdata
	var err error

	// first step, parse the data
	if len(conf.Separator) == 1 {
		data, err = parseCSV(conf, input)
	} else if conf.InputJSON {
		data, err = parseJSON(conf, input)
	} else {
		data, err = parseTabular(conf, input)
	}

	if err != nil {
		return data, err
	}

	// 2nd step, apply filters, code or transposers, if any
	postdata, changed, err := PostProcess(conf, &data)
	if err != nil {
		return data, err
	}

	if changed {
		return *postdata, nil
	}

	return data, err
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
			if matchPattern(conf, line) == conf.InvertMatch {
				// by default  -v is false, so if a  line does NOT
				// match the pattern, we will ignore it. However,
				// if the user specified -v, the matching is inverted,
				// so we ignore all lines, which DO match.
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

	return data, nil
}

/*
Parse JSON input.  We only support an array of  maps.

FIXME:  does not  preserve order,  so,  columns are  added in  some
random order as JSON maps are unordered
*/
func parseJSON(conf cfg.Config, input io.Reader) (Tabdata, error) {
	var data Tabdata
	var rawdata []map[string]string

	scanner := bufio.NewScanner(input)
	var rawjson string

	for scanner.Scan() {
		rawjson += scanner.Text()
	}

	if scanner.Err() != nil {
		return data, fmt.Errorf("failed to read from io.Reader: %w", scanner.Err())
	}

	if err := json.Unmarshal([]byte(rawjson), &rawdata); err != nil {
		return data, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	if len(rawdata) == 0 {
		return data, nil
	}

	// setup header fields
	headers := make(map[string]int, len(rawdata))
	var idx int

	for key := range rawdata[0] {
		data.headers = append(data.headers, key)
		headers[key] = idx

		idx++
	}

	// setup data entries
	for _, entry := range rawdata {
		row := make([]string, len(data.headers))
		var line string

		for idx, field := range data.headers {
			if Exists(entry, field) {
				row[idx] = entry[field]
			} else {
				row[idx] = ""
			}
			line += " " + row[idx]
		}

		// apply line filter
		if matchPattern(conf, line) == conf.InvertMatch {
			continue
		}

		data.entries = append(data.entries, row)
	}

	return data, nil
}

func PostProcess(conf cfg.Config, data *Tabdata) (*Tabdata, bool, error) {
	var modified bool

	// filter by field filters, if any
	filtereddata, changed, err := FilterByFields(conf, data)
	if err != nil {
		return data, false, fmt.Errorf("failed to filter fields: %w", err)
	}

	if changed {
		data = filtereddata
		modified = true
	}

	// check if transposers are valid and turn into Transposer structs
	if err := PrepareTransposerColumns(&conf, data); err != nil {
		return data, false, err
	}

	// transpose if demanded
	modifieddata, changed, err := TransposeFields(conf, data)
	if err != nil {
		return data, false, fmt.Errorf("failed to transpose fields: %w", err)
	}

	if changed {
		data = modifieddata
		modified = true
	}

	if conf.Debug {
		repr.Print(data)
	}

	return data, modified, nil
}
