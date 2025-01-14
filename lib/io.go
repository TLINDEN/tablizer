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
	"fmt"
	"io"
	"os"

	"github.com/tlinden/tablizer/cfg"
)

const RWRR = 0755

func ProcessFiles(conf *cfg.Config, args []string) error {
	fd, pattern, err := determineIO(conf, args)

	if err != nil {
		return err
	}

	if err := conf.PreparePattern(pattern); err != nil {
		return err
	}

	data, err := Parse(*conf, fd)
	if err != nil {
		return err
	}

	if err = ValidateConsistency(&data); err != nil {
		return err
	}

	err = PrepareColumns(conf, &data)
	if err != nil {
		return err
	}

	printData(os.Stdout, *conf, &data)

	return nil
}

func determineIO(conf *cfg.Config, args []string) (io.Reader, string, error) {
	var filehandle io.Reader
	var pattern string
	var haveio bool

	switch {
	case conf.InputFile == "-":
		filehandle = os.Stdin
		haveio = true
	case conf.InputFile != "":
		fd, err := os.OpenFile(conf.InputFile, os.O_RDONLY, RWRR)

		if err != nil {
			return nil, "", fmt.Errorf("failed to read input file %s: %w", conf.InputFile, err)
		}

		filehandle = fd
		haveio = true
	}

	if !haveio {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// we're reading from STDIN, which takes precedence over file args
			filehandle = os.Stdin
			haveio = true
		}
	}

	if len(args) > 0 {
		pattern = args[0]
		conf.Pattern = args[0]
	}

	if !haveio {
		return nil, "", errors.New("no file specified and nothing to read on stdin")
	}

	return filehandle, pattern, nil
}
