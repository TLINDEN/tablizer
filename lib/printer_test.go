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
	"os"
	"strings"
	"testing"
)

func TestPrinter(t *testing.T) {
	table := `ONE    TWO    THREE  
asd    igig   cxxxncnc  
19191  EDD 1  X`

	expects := map[string]string{
		"ascii": `ONE(1)	TWO(2)	THREE(3) 
asd   	igig  	cxxxncnc	
19191 	EDD 1 	X`,
		"orgtbl": `|--------+--------+----------|
| ONE(1) | TWO(2) | THREE(3) |
|--------+--------+----------|
| asd    | igig   | cxxxncnc |
|  19191 | EDD 1  | X        |
|--------+--------+----------|`,
		"markdown": `| ONE(1) | TWO(2) | THREE(3) |
|--------|--------|----------|
| asd    | igig   | cxxxncnc |
|  19191 | EDD 1  | X        |`,
	}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	origStdout := os.Stdout
	os.Stdout = w

	for mode, expect := range expects {
		OutputMode = mode
		fd := strings.NewReader(table)
		data, err := parseFile(fd, "")

		if err != nil {
			t.Errorf("Parser returned error: %s\nData processed so far: %+v", err, data)
		}

		printData(data)

		buf := make([]byte, 1024)
		n, err := r.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		buf = buf[:n]
		output := strings.TrimSpace(string(buf))

		if output != expect {
			t.Errorf("output mode: %s, got:\n%s\nwant:\n%s\n (%d <=> %d)", mode, output, expect, len(output), len(expect))
		}
	}

	// Restore
	os.Stdout = origStdout

}
