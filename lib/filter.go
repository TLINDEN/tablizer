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
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/tlinden/tablizer/cfg"
)

/*
 * [!]Match a  line, use fuzzy  search for normal pattern  strings and
 * regexp otherwise.
 */
func matchPattern(conf cfg.Config, line string) bool {
	if conf.UseFuzzySearch {
		return fuzzy.MatchFold(conf.Pattern, line)
	}

	return conf.PatternR.MatchString(line)
}

/*
 * Filter parsed data by fields. The  filter is positive, so if one or
 * more filters match on a row, it  will be kept, otherwise it will be
 * excluded.
 */
func FilterByFields(conf cfg.Config, data Tabdata) (Tabdata, bool, error) {
	if len(conf.Filters) == 0 {
		// no filters, no checking
		return Tabdata{}, false, nil
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

	return newdata, true, nil
}

func Exists[K comparable, V any](m map[K]V, v K) bool {
	if _, ok := m[v]; ok {
		return true
	}
	return false
}
