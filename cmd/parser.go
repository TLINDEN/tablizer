package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// contains a whole parsed table
type Tabdata struct {
	maxwidthHeader int   // longest header
	maxwidthPerCol []int // max width per column
	columns        int
	headerIndices  []map[string]int // [ {beg=>0, end=>17}, ... ]
	headers        []string         // [ "ID", "NAME", ...]
	entries        [][]string
}

func die(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

/*
   Parse tabular input. We split the  header (first line) by 2 or more
   spaces, remember the positions of  the header fields. We then split
   the data (everything after the first line) by those positions. That
   way we can turn "tabular data" (with fields containing whitespaces)
   into real tabular data. We re-tabulate our input if you will.
*/
func parseFile(input io.Reader, pattern string) Tabdata {
	data := Tabdata{}

	var scanner *bufio.Scanner
	var spaces = `\s\s+|$`

	if len(Separator) > 0 {
		spaces = Separator
	}

	hadFirst := false
	spacefinder := regexp.MustCompile(spaces)
	beg := 0

	scanner = bufio.NewScanner(input)

	for scanner.Scan() {
		line := scanner.Text()
		values := []string{}

		patternR, err := regexp.Compile(pattern)
		if err != nil {
			die(err)
		}

		if !hadFirst {
			// header processing
			parts := spacefinder.FindAllStringIndex(line, -1)
			data.columns = len(parts)
			// if Debug {
			// 	fmt.Println(parts)
			// }

			// process all header fields
			for _, part := range parts {
				// if Debug {
				// 	fmt.Printf("Part: <%s>\n", string(line[beg:part[0]]))
				//}

				// current field
				head := string(line[beg:part[0]])

				// register begin and end of field within line
				indices := make(map[string]int)
				indices["beg"] = beg
				if part[0] == part[1] {
					indices["end"] = 0
				} else {
					indices["end"] = part[1] - 1
				}

				// register widest header field
				headerlen := len(head)
				if headerlen > data.maxwidthHeader {
					data.maxwidthHeader = headerlen
				}

				// register fields data
				data.headerIndices = append(data.headerIndices, indices)
				data.headers = append(data.headers, head)

				// end of current field == begin of next one
				beg = part[1]

				// done
				hadFirst = true
			}
			// if Debug {
			// 	fmt.Println(data.headerIndices)
			// }
		} else {
			// data processing
			if len(pattern) > 0 {
				//fmt.Println(patternR.MatchString(line))
				if !patternR.MatchString(line) {
					continue
				}
			}

			idx := 0 // we cannot use the header index, because we could exclude columns

			for _, index := range data.headerIndices {
				value := ""
				if index["end"] == 0 {
					value = string(line[index["beg"]:])
				} else {
					value = string(line[index["beg"]:index["end"]])
				}

				width := len(strings.TrimSpace(value))

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
				values = append(values, value)

				idx++
			}
			if Debug {
				fmt.Println()
			}
			data.entries = append(data.entries, values)
		}
	}

	if scanner.Err() != nil {
		die(scanner.Err())
	}

	return data
}
