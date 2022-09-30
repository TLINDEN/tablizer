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
	"github.com/spf13/cobra"
	"os"
)

var version = "v1.0.1"

var rootCmd = &cobra.Command{
	Use:   "tablizer [regex] [file, ...]",
	Short: "[Re-]tabularize tabular data",
	Long:  `Manipulate tabular output of other programs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if Version {
			fmt.Printf("This is tablizer version %s\n", version)
			return nil
		}

		return process(args)
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
