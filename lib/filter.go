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
	"io"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/tlinden/tablizer/cfg"
)

/*
* [!]Match a  line, use fuzzy  search for normal pattern  strings and
* regexp otherwise.

		'foo bar'  foo, /bar/!  => false => line contains foo and not (not bar)
	    'foo nix'  foo, /bar/!  => ture  => line contains foo and (not bar)
		'foo bar'  foo, /bar/   => true  => line contains both foo and bar
		'foo nix'  foo, /bar/   => false => line does not contain bar
		'foo bar'  foo, /nix/   => false => line does not contain nix
*/
func matchPattern(conf cfg.Config, line string) bool {
	if len(conf.Patterns) == 0 {
		// any line always matches ""
		return true
	}

	if conf.UseFuzzySearch {
		// fuzzy search only considers the 1st pattern
		return fuzzy.MatchFold(conf.Patterns[0].Pattern, line)
	}

	var match int

	//fmt.Printf("<%s>\n", line)
	for _, re := range conf.Patterns {
		patmatch := re.PatternRe.MatchString(line)
		if re.Negate {
			// toggle the meaning of match
			patmatch = !patmatch
		}

		if patmatch {
			match++
		}

		//fmt.Printf("patmatch: %t, match: %d, pattern: %s, negate: %t\n", patmatch, match, re.Pattern, re.Negate)
	}

	// fmt.Printf("result: %t\n", match == len(conf.Patterns))
	//fmt.Println()
	return match == len(conf.Patterns)
}

/*
 * Filter parsed data by fields. The  filter is positive, so if one or
 * more filters match on a row, it  will be kept, otherwise it will be
 * excluded.
 */
func FilterByFields(conf cfg.Config, data *Tabdata) (*Tabdata, bool, error) {
	if len(conf.Filters) == 0 {
		// no filters, no checking
		return nil, false, nil
	}

	newdata := data.CloneEmpty()

	for _, row := range data.entries {
		keep := true

		for idx, header := range data.headers {
			if !Exists(conf.Filters, strings.ToLower(header)) {
				// do not filter by unspecified field
				continue
			}

			if !conf.Filters[strings.ToLower(header)].MatchString(row[idx]) {
				// there IS a filter, but it doesn't match
				keep = false

				break
			}
		}

		if keep == !conf.InvertMatch {
			// also apply -v
			newdata.entries = append(newdata.entries, row)
		}
	}

	return &newdata, true, nil
}

/*
 * Transpose fields using search/replace regexp.
 */
func TransposeFields(conf cfg.Config, data *Tabdata) (*Tabdata, bool, error) {
	if len(conf.UseTransposers) == 0 {
		// nothing to be done
		return nil, false, nil
	}

	newdata := data.CloneEmpty()
	transposed := false

	for _, row := range data.entries {
		transposedrow := false

		for idx := range data.headers {
			transposeidx, hasone := findindex(conf.UseTransposeColumns, idx+1)
			if hasone {
				row[idx] =
					conf.UseTransposers[transposeidx].Search.ReplaceAllString(
						row[idx],
						conf.UseTransposers[transposeidx].Replace,
					)
				transposedrow = true
			}
		}

		if transposedrow {
			// also apply -v
			newdata.entries = append(newdata.entries, row)
			transposed = true
		}
	}

	return &newdata, transposed, nil
}

/* generic map.Exists(key) */
func Exists[K comparable, V any](m map[K]V, v K) bool {
	if _, ok := m[v]; ok {
		return true
	}

	return false
}

/*
 * Filters the whole input lines, returns filtered lines
 */
func FilterByPattern(conf cfg.Config, input io.Reader) (io.Reader, error) {
	if len(conf.Patterns) == 0 {
		return input, nil
	}

	scanner := bufio.NewScanner(input)
	lines := []string{}
	hadFirst := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if hadFirst {
			// don't match 1st line, it's the header
			if matchPattern(conf, line) == conf.InvertMatch {
				// by default  -v is false, so if a  line does NOT
				// match the pattern, we will ignore it. However,
				// if the user specified -v, the matching is inverted,
				// so we ignore all lines, which DO match.
				continue
			}
		}

		lines = append(lines, line)

		hadFirst = true
	}

	return strings.NewReader(strings.Join(lines, "\n")), nil
}
