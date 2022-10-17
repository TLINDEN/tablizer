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
	"github.com/gookit/color"
	"io"
	"os"
)

func ProcessFiles(args []string) error {
	fds, pattern, err := determineIO(args)

	if !isTerminal(os.Stdout) {
		color.Disable()
	} else {
		level := color.TermColorLevel()
		MatchFG = Colors[level]["fg"]
		MatchBG = Colors[level]["bg"]
	}

	if err != nil {
		return err
	}

	for _, fd := range fds {
		data, err := parseFile(fd, pattern)
		if err != nil {
			return err
		}

		err = PrepareColumns(&data)
		if err != nil {
			return err
		}

		printData(os.Stdout, &data)
	}

	return nil
}

func determineIO(args []string) ([]io.Reader, string, error) {
	var pattern string
	var fds []io.Reader
	var havefiles bool

	if len(args) > 0 {
		// threre were args left, take a look
		if _, err := os.Stat(args[0]); err != nil {
			// first  one is  not a  file, consider  it as  regexp and
			// shift arg list
			pattern = args[0]
			Pattern = args[0] // FIXME
			args = args[1:]
		}

		if len(args) > 0 {
			// only files
			for _, file := range args {
				fd, err := os.OpenFile(file, os.O_RDONLY, 0755)

				if err != nil {
					return nil, "", err
				}

				fds = append(fds, fd)
			}
			havefiles = true
		}
	}

	if !havefiles {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			fds = append(fds, os.Stdin)
		} else {
			return nil, "", errors.New("No file specified and nothing to read on stdin!")
		}
	}

	return fds, pattern, nil
}
