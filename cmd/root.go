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
	"daemon.de/tablizer/lib"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "tablizer [regex] [file, ...]",
	Short: "[Re-]tabularize tabular data",
	Long:  `Manipulate tabular output of other programs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if lib.ShowVersion {
			fmt.Printf("This is tablizer version %s\n", lib.Version)
			return nil
		}

		lib.PrepareColumns()

		err := lib.PrepareModeFlags()
		if err != nil {
			return err
		}

		return lib.ProcessFiles(args)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&lib.Debug, "debug", "d", false, "Enable debugging")
	rootCmd.PersistentFlags().BoolVarP(&lib.NoNumbering, "no-numbering", "n", false, "Disable header numbering")
	rootCmd.PersistentFlags().BoolVarP(&lib.ShowVersion, "version", "v", false, "Print program version")
	rootCmd.PersistentFlags().StringVarP(&lib.Separator, "separator", "s", "", "Custom field separator")
	rootCmd.PersistentFlags().StringVarP(&lib.Columns, "columns", "c", "", "Only show the speficied columns (separated by ,)")

	// output flags, only 1 allowed
	rootCmd.PersistentFlags().BoolVarP(&lib.OutflagExtended, "extended", "X", false, "Enable extended output")
	rootCmd.PersistentFlags().BoolVarP(&lib.OutflagMarkdown, "markdown", "M", false, "Enable markdown table output")
	rootCmd.PersistentFlags().BoolVarP(&lib.OutflagOrgtable, "orgtbl", "O", false, "Enable org-mode table output")
	rootCmd.MarkFlagsMutuallyExclusive("extended", "markdown", "orgtbl")

	// same thing but more common, takes precedence over above group
	rootCmd.PersistentFlags().StringVarP(&lib.OutputMode, "output", "o", "", "Output mode - one of: orgtbl, markdown, extended, ascii(default)")
}
