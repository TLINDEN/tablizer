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
	"bufio"
	"errors"
	"fmt"
	"github.com/alecthomas/repr"
	"io"
	"regexp"
	"strings"
)

// contains a whole parsed table
type Tabdata struct {
	maxwidthHeader int      // longest header
	maxwidthPerCol []int    // max width per column
	columns        int      // count
	headers        []string // [ "ID", "NAME", ...]
	entries        [][]string
}

/*
   Parse tabular input.
*/
func parseFile(input io.Reader, pattern string) (Tabdata, error) {
	data := Tabdata{}

	var scanner *bufio.Scanner

	hadFirst := false
	separate := regexp.MustCompile(Separator)
	patternR, err := regexp.Compile(pattern)
	if err != nil {
		return data, errors.Unwrap(fmt.Errorf("Regexp pattern %s is invalid: %w", pattern, err))
	}

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
			if len(pattern) > 0 {
				if patternR.MatchString(line) == InvertMatch {
					// by default  -v is false, so if a  line does NOT
					// match the pattern, we will ignore it. However,
					// if the user specified -v, the matching is inverted,
					// so we ignore all lines, which DO match.
					continue
				}
			}

			idx := 0 // we cannot use the header index, because we could exclude columns
			values := []string{}
			for _, part := range parts {
				width := len(strings.TrimSpace(part))

				if len(data.maxwidthPerCol)-1 < idx {
					data.maxwidthPerCol = append(data.maxwidthPerCol, width)
				} else {
					if width > data.maxwidthPerCol[idx] {
						data.maxwidthPerCol[idx] = width
					}
				}

				// if Debug {
				// 	fmt.Printf("<%s> ", value)
				// }
				values = append(values, strings.TrimSpace(part))

				idx++
			}
			data.entries = append(data.entries, values)
		}
	}

	if scanner.Err() != nil {
		return data, errors.Unwrap(fmt.Errorf("Failed to read from io.Reader: %w", scanner.Err()))
	}

	if Debug {
		repr.Print(data)
	}

	return data, nil
}
