/*
Copyright © 2025 Thomas von Dein

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
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tiagomelo/go-clipboard/clipboard"
	"github.com/tlinden/tablizer/cfg"
)

var yanktests = []struct {
	name   string
	yank   []int // -y$colum,$column... after processing
	filter string
	expect string
}{
	{
		name:   "one",
		yank:   []int{1},
		filter: "beta",
	},
}

func DISABLED_TestYankColumns(t *testing.T) {
	cb := clipboard.New()

	for _, testdata := range yanktests {
		testname := fmt.Sprintf("yank-%s-filter-%s",
			testdata.name, testdata.filter)
		t.Run(testname, func(t *testing.T) {
			conf := cfg.Config{
				OutputMode:     cfg.ASCII,
				UseYankColumns: testdata.yank,
				NoColor:        true,
			}

			conf.ApplyDefaults()
			data := newData() // defined in printer_test.go, reused here

			var writer bytes.Buffer
			printData(&writer, conf, &data)

			got, err := cb.PasteText()

			assert.NoError(t, err)
			assert.EqualValues(t, testdata.expect, got)
		})
	}
}
