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
	new := Tabdata{
		maxwidthHeader: data.maxwidthHeader,
		columns:        data.columns,
		headers:        data.headers,
	}

	return new
}
