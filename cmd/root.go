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
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tlinden/tablizer/lib"
	"log"
	"os"
	"os/exec"
)

var ShowManual = false

func man() {
	man := exec.Command("less", "-")

	var b bytes.Buffer
	b.Write([]byte(manpage))

	man.Stdout = os.Stdout
	man.Stdin = &b
	man.Stderr = os.Stderr

	err := man.Run()

	if err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "tablizer [regex] [file, ...]",
	Short: "[Re-]tabularize tabular data",
	Long:  `Manipulate tabular output of other programs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if lib.ShowVersion {
			fmt.Printf("This is tablizer version %s\n", lib.VERSION)
			return nil
		}

		if ShowManual {
			man()
			return nil
		}

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
	rootCmd.PersistentFlags().BoolVarP(&lib.NoColor, "no-color", "N", false, "Disable pattern highlighting")
	rootCmd.PersistentFlags().BoolVarP(&lib.ShowVersion, "version", "V", false, "Print program version")
	rootCmd.PersistentFlags().BoolVarP(&lib.InvertMatch, "invert-match", "v", false, "select non-matching rows")
	rootCmd.PersistentFlags().BoolVarP(&ShowManual, "man", "m", false, "Display manual page")
	rootCmd.PersistentFlags().StringVarP(&lib.Separator, "separator", "s", lib.DefaultSeparator, "Custom field separator")
	rootCmd.PersistentFlags().StringVarP(&lib.Columns, "columns", "c", "", "Only show the speficied columns (separated by ,)")

	// sort options
	rootCmd.PersistentFlags().IntVarP(&lib.SortByColumn, "sort-by", "k", 0, "Sort by column (default: 1)")
	rootCmd.PersistentFlags().BoolVarP(&lib.SortDescending, "sort-desc", "D", false, "Sort in descending order (default: ascending)")
	rootCmd.PersistentFlags().BoolVarP(&lib.SortNumeric, "sort-numeric", "i", false, "sort according to string numerical value")
	rootCmd.PersistentFlags().BoolVarP(&lib.SortTime, "sort-time", "t", false, "sort according to time string")
	rootCmd.PersistentFlags().BoolVarP(&lib.SortAge, "sort-age", "a", false, "sort according to age (duration) string")

	// output flags, only 1 allowed, hidden, since just short cuts
	rootCmd.PersistentFlags().BoolVarP(&lib.OutflagExtended, "extended", "X", false, "Enable extended output")
	rootCmd.PersistentFlags().BoolVarP(&lib.OutflagMarkdown, "markdown", "M", false, "Enable markdown table output")
	rootCmd.PersistentFlags().BoolVarP(&lib.OutflagOrgtable, "orgtbl", "O", false, "Enable org-mode table output")
	rootCmd.PersistentFlags().BoolVarP(&lib.OutflagShell, "shell", "S", false, "Enable shell mode output")
	rootCmd.MarkFlagsMutuallyExclusive("extended", "markdown", "orgtbl", "shell")
	rootCmd.Flags().MarkHidden("extended")
	rootCmd.Flags().MarkHidden("orgtbl")
	rootCmd.Flags().MarkHidden("markdown")
	rootCmd.Flags().MarkHidden("shell")

	// same thing but more common, takes precedence over above group
	rootCmd.PersistentFlags().StringVarP(&lib.OutputMode, "output", "o", "", "Output mode - one of: orgtbl, markdown, extended, shell, ascii(default)")
}
