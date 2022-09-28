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
package cmd

import (
	"fmt"
	"github.com/alecthomas/repr"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

var version = "v1.0.0"

var rootCmd = &cobra.Command{
	Use:   "tablizer [regex] [file, ...]",
	Short: "[Re-]tabularize tabular data",
	Long:  `Manipulate tabular output of other programs`,
	Run: func(cmd *cobra.Command, args []string) {
		if Version {
			fmt.Printf("This is tablizer version %s\n", version)
			return
		}

		var pattern string
		havefiles := false

		if len(Columns) > 0 {
			for _, use := range strings.Split(Columns, ",") {
				usenum, err := strconv.Atoi(use)
				if err != nil {
					die(err)
				}
				UseColumns = append(UseColumns, usenum)
			}
		}

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
					printTable(data)
				}
				havefiles = true
			}
		}

		if !havefiles {
			data := parseFile(os.Stdin, pattern)
			if Debug {
				repr.Print(data)
			}
			printTable(data)
		}
	},
}

var Debug bool
var XtendedOut bool
var NoNumbering bool
var Version bool
var Columns string
var UseColumns []int
var Separator string

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Enable debugging")
	rootCmd.PersistentFlags().BoolVarP(&XtendedOut, "extended", "x", false, "Enable extended output")
	rootCmd.PersistentFlags().BoolVarP(&NoNumbering, "no-numbering", "n", false, "Disable header numbering")
	rootCmd.PersistentFlags().BoolVarP(&Version, "version", "v", false, "Print program version")
	rootCmd.PersistentFlags().StringVarP(&Separator, "separator", "s", "", "Custom field separator")
	rootCmd.PersistentFlags().StringVarP(&Columns, "columns", "c", "", "Only show the speficied columns (separated by ,)")
}
