/*
Copyright © 2022-2024 Thomas von Dein

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
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tlinden/tablizer/cfg"
	"github.com/tlinden/tablizer/lib"
)

func man() {
	man := exec.Command("less", "-")

	var buffer bytes.Buffer

	buffer.Write([]byte(manpage))

	man.Stdout = os.Stdout
	man.Stdin = &buffer
	man.Stderr = os.Stderr

	err := man.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func completion(cmd *cobra.Command, mode string) error {
	switch mode {
	case "bash":
		return cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		return cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		return cmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	default:
		return errors.New("invalid shell parameter! Valid ones: bash|zsh|fish|powershell")
	}
}

// we die with exit 1 if there's an error
func wrapE(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func Execute() {
	var (
		conf           cfg.Config
		ShowManual     bool
		ShowVersion    bool
		ShowCompletion string
		modeflag       cfg.Modeflag
		sortmode       cfg.Sortmode
	)

	var rootCmd = &cobra.Command{
		Use:   "tablizer [regex] [file, ...]",
		Short: "[Re-]tabularize tabular data",
		Long:  `Manipulate tabular output of other programs`,

		Run: func(cmd *cobra.Command, args []string) {
			if ShowVersion {
				fmt.Println(cfg.Getversion())

				return
			}

			if ShowManual {
				man()

				return
			}

			if len(ShowCompletion) > 0 {
				wrapE(completion(cmd, ShowCompletion))

				return
			}

			// Setup
			wrapE(conf.ParseConfigfile())

			conf.CheckEnv()
			conf.PrepareModeFlags(modeflag)
			conf.PrepareSortFlags(sortmode)

			wrapE(conf.PrepareFilters())

			conf.DetermineColormode()
			conf.ApplyDefaults()

			// setup lisp env, load plugins etc
			wrapE(lib.SetupLisp(&conf))

			// actual execution starts here
			wrapE(lib.ProcessFiles(&conf, args))
		},
	}

	// options
	rootCmd.PersistentFlags().BoolVarP(&conf.Debug, "debug", "d", false,
		"Enable debugging")
	rootCmd.PersistentFlags().BoolVarP(&conf.NoNumbering, "no-numbering", "n", false,
		"Disable header numbering")
	rootCmd.PersistentFlags().BoolVarP(&conf.NoHeaders, "no-headers", "H", false,
		"Disable header display")
	rootCmd.PersistentFlags().BoolVarP(&conf.NoColor, "no-color", "N", false,
		"Disable pattern highlighting")
	rootCmd.PersistentFlags().BoolVarP(&ShowVersion, "version", "V", false,
		"Print program version")
	rootCmd.PersistentFlags().BoolVarP(&conf.InvertMatch, "invert-match", "v", false,
		"select non-matching rows")
	rootCmd.PersistentFlags().BoolVarP(&ShowManual, "man", "m", false,
		"Display manual page")
	rootCmd.PersistentFlags().BoolVarP(&conf.UseFuzzySearch, "fuzzy", "z", false,
		"Use fuzzy searching")
	rootCmd.PersistentFlags().BoolVarP(&conf.UseHighlight, "highlight-lines", "L", false,
		"Use alternating background colors")
	rootCmd.PersistentFlags().StringVarP(&ShowCompletion, "completion", "", "",
		"Display completion code")
	rootCmd.PersistentFlags().StringVarP(&conf.Separator, "separator", "s", cfg.DefaultSeparator,
		"Custom field separator")
	rootCmd.PersistentFlags().StringVarP(&conf.Columns, "columns", "c", "",
		"Only show the speficied columns (separated by ,)")
	rootCmd.PersistentFlags().StringVarP(&conf.TransposeColumns, "transpose-columns", "T", "",
		"Transpose the speficied columns (separated by ,)")

	// sort options
	rootCmd.PersistentFlags().IntVarP(&conf.SortByColumn, "sort-by", "k", 0,
		"Sort by column (default: 1)")

	// sort mode, only 1 allowed
	rootCmd.PersistentFlags().BoolVarP(&conf.SortDescending, "sort-desc", "D", false,
		"Sort in descending order (default: ascending)")
	rootCmd.PersistentFlags().BoolVarP(&sortmode.Numeric, "sort-numeric", "i", false,
		"sort according to string numerical value")
	rootCmd.PersistentFlags().BoolVarP(&sortmode.Time, "sort-time", "t", false,
		"sort according to time string")
	rootCmd.PersistentFlags().BoolVarP(&sortmode.Age, "sort-age", "a", false,
		"sort according to age (duration) string")
	rootCmd.MarkFlagsMutuallyExclusive("sort-numeric", "sort-time",
		"sort-age")

	// output flags, only 1 allowed
	rootCmd.PersistentFlags().BoolVarP(&modeflag.X, "extended", "X", false,
		"Enable extended output")
	rootCmd.PersistentFlags().BoolVarP(&modeflag.M, "markdown", "M", false,
		"Enable markdown table output")
	rootCmd.PersistentFlags().BoolVarP(&modeflag.O, "orgtbl", "O", false,
		"Enable org-mode table output")
	rootCmd.PersistentFlags().BoolVarP(&modeflag.S, "shell", "S", false,
		"Enable shell mode output")
	rootCmd.PersistentFlags().BoolVarP(&modeflag.Y, "yaml", "Y", false,
		"Enable yaml output")
	rootCmd.PersistentFlags().BoolVarP(&modeflag.C, "csv", "C", false,
		"Enable CSV output")
	rootCmd.PersistentFlags().BoolVarP(&modeflag.A, "ascii", "A", false,
		"Enable ASCII output (default)")
	rootCmd.MarkFlagsMutuallyExclusive("extended", "markdown", "orgtbl",
		"shell", "yaml", "csv")

	// lisp options
	rootCmd.PersistentFlags().StringVarP(&conf.LispLoadPath, "load-path", "l", cfg.DefaultLoadPath,
		"Load path for lisp plugins (expects *.zy files)")

	// config file
	rootCmd.PersistentFlags().StringVarP(&conf.Configfile, "config", "f", cfg.DefaultConfigfile,
		"config file (default: ~/.config/tablizer/config)")

	// filters
	rootCmd.PersistentFlags().StringArrayVarP(&conf.Rawfilters,
		"filter", "F", nil, "Filter by field (field=regexp)")
	rootCmd.PersistentFlags().StringArrayVarP(&conf.Transposers,
		"regex-transposer", "R", nil, "apply /search/replace/ regexp to fields given in -T")

	// input
	rootCmd.PersistentFlags().StringVarP(&conf.InputFile, "read-file", "r", "",
		"Read input data from file")

	rootCmd.SetUsageTemplate(strings.TrimSpace(usage) + "\n")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
