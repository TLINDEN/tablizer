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
	"errors"
	"github.com/alecthomas/repr"
	"os"
)

func ProcessFiles(args []string) error {
	var pattern string
	havefiles := false

	prepareColumns()

	if len(args) > 0 {
		if _, err := os.Stat(args[0]); err != nil {
			pattern = args[0]
			args = args[1:]
		}

		if len(args) > 0 {
			for _, file := range args {
				fd, err := os.OpenFile(file, os.O_RDONLY, 0755)
				if err != nil {
					die(err)
				}

				data := parseFile(fd, pattern)
				if Debug {
					repr.Print(data)
				}
				printData(data)
			}
			havefiles = true
		}
	}

	if !havefiles {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			data := parseFile(os.Stdin, pattern)
			if Debug {
				repr.Print(data)
			}
			printData(data)
		} else {
			return errors.New("No file specified and nothing to read on stdin!")
		}
	}

	return nil
}
