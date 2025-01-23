/*
Copyright © 2022-2025 Thomas von Dein

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
	"log"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/tlinden/tablizer/cfg"
)

func yankColumns(conf cfg.Config, data *Tabdata) {
	var yank []string

	if len(data.entries) == 0 || len(conf.UseYankColumns) == 0 {
		return
	}

	for _, row := range data.entries {
		for i, field := range row {
			for _, idx := range conf.UseYankColumns {
				if i == idx-1 {
					yank = append(yank, field)
				}
			}
		}
	}

	if len(yank) > 0 {
		setprimary()
		if err := clipboard.WriteAll(strings.Join(yank, " ")); err != nil {
			log.Fatalln("error writing string to clipboard:", err)
		}
	}
}
