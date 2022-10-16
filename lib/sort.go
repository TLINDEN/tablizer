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
	"github.com/araddon/dateparse"
	"regexp"
	"sort"
	"strconv"
)

func sortTable(data *Tabdata, col int) {
	if col <= 0 {
		// no sorting wanted
		return
	}

	col-- // ui starts counting by 1, but use 0 internally

	// sanity checks
	if len(data.entries) == 0 {
		return
	}

	if col >= len(data.headers) {
		// fall back to default column
		col = 0
	}

	// actual sorting
	sort.SliceStable(data.entries, func(i, j int) bool {
		return compare(data.entries[i][col], data.entries[j][col])
	})
}

func compare(a string, b string) bool {
	var comp bool

	switch SortMode {
	case "numeric":
		left, err := strconv.Atoi(a)
		if err != nil {
			left = 0
		}
		right, err := strconv.Atoi(b)
		if err != nil {
			right = 0
		}
		comp = left < right
	case "duration":
		left := duration2int(a)
		right := duration2int(b)
		comp = left < right
	case "time":
		left, _ := dateparse.ParseAny(a)
		right, _ := dateparse.ParseAny(b)
		comp = left.Unix() < right.Unix()
	default:
		comp = a < b
	}

	if SortDescending {
		comp = !comp
	}

	return comp
}

/*
   We could use time.ParseDuration(), but this doesn't support days.

   We  could also  use github.com/xhit/go-str2duration/v2,  which does
   the job,  but it's  just another dependency,  just for  this little
   gem. And  we don't need a  time.Time value. And int  is good enough
   for duration comparision.

   Convert a  durartion into  an integer.  Valid  time units  are "s",
   "m", "h" and "d".
*/
func duration2int(duration string) int {
	re := regexp.MustCompile(`(\d+)([dhms])`)
	seconds := 0

	for _, match := range re.FindAllStringSubmatch(duration, -1) {
		if len(match) == 3 {
			v, _ := strconv.Atoi(match[1])
			switch match[2][0] {
			case 'd':
				seconds += v * 86400
			case 'h':
				seconds += v * 3600
			case 'm':
				seconds += v * 60
			case 's':
				seconds += v
			}
		}
	}

	return seconds
}
