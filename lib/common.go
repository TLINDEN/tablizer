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

// contains a whole parsed table
type Tabdata struct {
	maxwidthHeader int      // longest header
	columns        int      // count
	headers        []string // [ "ID", "NAME", ...]
	entries        [][]string
}

func (data *Tabdata) CloneEmpty() Tabdata {
	newdata := Tabdata{
		maxwidthHeader: data.maxwidthHeader,
		columns:        data.columns,
		headers:        data.headers,
	}

	return newdata
}

// add a TAB (\t) in front of every cell, but not the first
func (data *Tabdata) TabEntries() [][]string {
	newentries := make([][]string, len(data.entries))

	for rowidx, row := range data.entries {
		newentries[rowidx] = make([]string, len(row))

		for colidx, cell := range row {
			switch colidx {
			case 0:
				newentries[rowidx][colidx] = cell
			default:
				newentries[rowidx][colidx] = "\t" + cell
			}
		}
	}

	return newentries
}

// add a TAB (\t) in front of every header, but not the first
func (data *Tabdata) TabHeaders() []string {
	newheaders := make([]string, len(data.headers))

	for colidx, cell := range data.headers {
		switch colidx {
		case 0:
			newheaders[colidx] = cell
		default:
			newheaders[colidx] = "\t" + cell
		}
	}

	return newheaders
}
