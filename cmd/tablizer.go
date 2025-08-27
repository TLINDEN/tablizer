package cmd

var manpage = `
NAME
    tablizer - Manipulate tabular output of other programs

SYNOPSIS
        Usage:
          tablizer [regex,...] [file, ...] [flags]
    
        Operational Flags:
          -c, --columns string              Only show the speficied columns (separated by ,)
          -v, --invert-match                select non-matching rows
          -n, --numbering                   Enable header numbering
          -N, --no-color                    Disable pattern highlighting
          -H, --no-headers                  Disable headers display
          -s, --separator string            Custom field separator
          -k, --sort-by int|name            Sort by column (default: 1)
          -z, --fuzzy                       Use fuzzy search [experimental]
          -F, --filter field[!]=reg         Filter given field with regex, can be used multiple times
          -T, --transpose-columns string    Transpose the speficied columns (separated by ,)
          -R, --regex-transposer /from/to/  Apply /search/replace/ regexp to fields given in -T
          -I, --interactive                 Interactively filter and select rows

        Output Flags (mutually exclusive):
          -X, --extended                    Enable extended output
          -M, --markdown                    Enable markdown table output
          -O, --orgtbl                      Enable org-mode table output
          -S, --shell                       Enable shell evaluable output
          -Y, --yaml                        Enable yaml output
          -C, --csv                         Enable CSV output
          -A, --ascii                       Default output mode, ascii tabular
          -L, --hightlight-lines            Use alternating background colors for tables
          -y, --yank-columns                Yank specified columns (separated by ,) to clipboard,
                                            space separated

        Sort Mode Flags (mutually exclusive):
          -a, --sort-age                    sort according to age (duration) string
          -D, --sort-desc                   Sort in descending order (default: ascending)
          -i, --sort-numeric                sort according to string numerical value
          -t, --sort-time                   sort according to time string

        Other Flags:
              --completion <shell>         Generate the autocompletion script for <shell>
          -f, --config <file>              Configuration file (default: ~/.config/tablizer/config)
          -d, --debug                      Enable debugging
          -h, --help                       help for tablizer
          -m, --man                        Display manual page
          -V, --version                    Print program version

DESCRIPTION
    Many programs generate tabular output. But sometimes you need to
    post-process these tables, you may need to remove one or more columns or
    you may want to filter for some pattern (See PATTERNS) or you may need
    the output in another program and need to parse it somehow. Standard
    unix tools such as awk(1), grep(1) or column(1) may help, but sometimes
    it's a tedious business.

    Let's take the output of the tool kubectl. It contains cells with
    withespace and they do not separate columns by TAB characters. This is
    not easy to process.

    You can use tablizer to do these and more things.

    tablizer analyses the header fields of a table, registers the column
    positions of each header field and separates columns by those positions.

    Without any options it reads its input from "STDIN", but you can also
    specify a file as a parameter. If you want to reduce the output by some
    regular expression, just specify it as its first parameter. You may also
    use the -v option to exclude all rows which match the pattern. Hence:

       # read from STDIN
       kubectl get pods | tablizer

       # read a file
       tablizer filename

       # search for pattern in a file (works like grep)
       tablizer regex filename

       # search for pattern in STDIN
       kubectl get pods | tablizer regex

    The output looks like the original one but every header field will have
    a numer associated with it, e.g.:

       NAME(1) READY(2) STATUS(3) RESTARTS(4) AGE(5)

    These numbers denote the column and you can use them to specify which
    columns you want to have in your output (see COLUMNS:

       kubectl get pods | tablizer -c1,3

    You can specify the numbers in any order but output will always follow
    the original order.

    The numbering can be suppressed by using the -n option.

    By default tablizer shows a header containing the names of each column.
    This can be disabled using the -H option. Be aware that this only
    affects tabular output modes. Shell, Extended, Yaml and CSV output modes
    always use the column names.

    By default, if a pattern has been speficied, matches will be
    highlighted. You can disable this behavior with the -N option.

    Use the -k option to specify by which column to sort the tabular data
    (as in GNU sort(1)). The default sort column is the first one. You can
    specify column numbers or names. Column numbers start with 1, names are
    case insensitive. You can specify multiple columns separated by comma to
    sort, but the type must be the same. For example if you want to sort
    numerically, all columns must be numbers. If you use column numbers,
    then be aware, that these are the numbers before column extraction. For
    example if you have a table with 4 columns and specify "-c4", then only
    1 column (the fourth) will be printed, however if you want to sort by
    this column, you'll have to specify "-k4".

    The default sort order is ascending. You can change this to descending
    order using the option -D. The default sort order is by alphanumeric
    string, but there are other sort modes:

    -a --sort-age
        Sorts duration strings like "1d4h32m51s".

    -i --sort-numeric
        Sorts numeric fields.

    -t --sort-time
        Sorts timestamps.

    Finally the -d option enables debugging output which is mostly useful
    for the developer.

  PATTERNS AND FILTERING
    You can reduce the rows being displayed by using one or more regular
    expression patterns. The regexp language being used is the one of
    GOLANG, refer to the syntax cheat sheet here:
    <https://pkg.go.dev/regexp/syntax>.

    If you want to read a more comprehensive documentation about the topic
    and have perl installed you can read it with:

        perldoc perlre

    Or read it online: <https://perldoc.perl.org/perlre>. But please note
    that the GO regexp engine does NOT support all perl regex terms,
    especially look-ahead and look-behind.

    If you want to supply flags to a regex, then surround it with slashes
    and append the flag. The following flags are supported:

        i => case insensitive
        ! => negative match

    Example for a case insensitive search:

        kubectl get pods -A | tablizer "/account/i"

    If you use the "!" flag, then the regex match will be negated, that is,
    if a line in the input matches the given regex, but "!" is supplied,
    tablizer will NOT include it in the output.

    For example, here we want to get all lines matching "foo" but not "bar":

        cat table | tablizer foo '/bar/!'

    This would match a line "foo zorro" but not "foo bar".

    The flags can also be combined.

    You can also use the experimental fuzzy search feature by providing the
    option -z, in which case the pattern is regarded as a fuzzy search term,
    not a regexp.

    Sometimes you want to filter by one or more columns. You can do that
    using the -F option. The option can be specified multiple times and has
    the following format:

        fieldname=regexp

    Fieldnames (== columns headers) are case insensitive.

    If you specify more than one filter, both filters have to match (AND
    operation).

    These field filters can also be negated:

        fieldname!=regexp

    If the option -v is specified, the filtering is inverted.

  COLUMNS
    The parameter -c can be used to specify, which columns to display. By
    default tablizer numerizes the header names and these numbers can be
    used to specify which header to display, see example above.

    However, beside numbers, you can also use regular expressions with -c,
    also separated by comma. And you can mix column numbers with regexps.

    Lets take this table:

            PID TTY          TIME CMD
          14001 pts/0    00:00:00 bash
          42871 pts/0    00:00:00 ps
          42872 pts/0    00:00:00 sed

    We want to see only the CMD column and use a regex for this:

        ps | tablizer -s '\s+' -c C
        CMD(4)
        bash
        ps
        tablizer
        sed

    where "C" is our regexp which matches CMD.

    If a column specifier doesn't look like a regular expression, matching
    against header fields will be case insensitive. So, if you have a field
    with the name "ID" then these will all match: "-c id", "-c Id". The same
    rule applies to the options "-T" and "-F".

  TRANSPOSE FIELDS USING REGEXPS
    You can manipulate field contents using regular expressions. You have to
    tell tablizer which field[s] to operate on using the option "-T" and the
    search/replace pattern using "-R". The number of columns and patterns
    must match.

    A search/replace pattern consists of the following elements:

        /search-regexp/replace-string/

    The separator can be any valid character. Especially if you want to use
    a regexp containing the "/" character, eg:

        |search-regexp|replace-string|

    Example:

        cat t/testtable2
        NAME  DURATION
        x     10
        a     100
        z     0
        u     4
        k     6
    
        cat t/testtable2 | tablizer -T2 -R '/^\d/4/' -n
        NAME    DURATION 
        x       40      
        a       400     
        z       4       
        u       4       
        k       4

  OUTPUT MODES
    There might be cases when the tabular output of a program is way too
    large for your current terminal but you still need to see every column.
    In such cases the -o extended or -X option can be useful which enables
    *extended mode*. In this mode, each row will be printed vertically,
    header left, value right, aligned by the field widths. Here's an
    example:

        kubectl get pods | ./tablizer -o extended
            NAME: repldepl-7bcd8d5b64-7zq4l  
           READY: 1/1    
          STATUS: Running  
        RESTARTS: 1 (71m ago)  
             AGE: 5h28m

    You can of course still use a regex to reduce the number of rows
    displayed.

    The option -o shell can be used if the output has to be processed by the
    shell, it prints variable assignments for each cell, one line per row:

        kubectl get pods | ./tablizer -o extended ./tablizer -o shell
        NAME="repldepl-7bcd8d5b64-7zq4l" READY="1/1" STATUS="Running" RESTARTS="9 (47m ago)" AGE="4d23h" 
        NAME="repldepl-7bcd8d5b64-m48n8" READY="1/1" STATUS="Running" RESTARTS="9 (47m ago)" AGE="4d23h" 
        NAME="repldepl-7bcd8d5b64-q2bf4" READY="1/1" STATUS="Running" RESTARTS="9 (47m ago)" AGE="4d23h"

    You can use this in an eval loop.

    Beside normal ascii mode (the default) and extended mode there are more
    output modes available: orgtbl which prints an Emacs org-mode table and
    markdown which prints a Markdown table, yaml, which prints yaml encoding
    and CSV mode, which prints a comma separated value file.

  PUT FIELDS TO CLIPBOARD
    You can let tablizer put fields to the clipboard using the option "-y".
    This best fits the use-case when the result of your filtering yields
    just one row. For example:

        cloudctl cluster ls | tablizer -yid matchbox

    If "matchbox" matches one cluster, you can immediately use the id of
    that cluster somewhere else and paste it. Of course, if there are
    multiple matches, then all id's will be put into the clipboard separated
    by one space.

  ENVIRONMENT VARIABLES
    tablizer supports certain environment variables which use can use to
    influence program behavior. Commandline flags have always precedence
    over environment variables.

    <T_HEADER_NUMBERING> - enable numbering of header fields, like -n.
    <T_COLUMNS> - comma separated list of columns to output, like -c
    <NO_COLORS> - disable colorization of matches, like -N

  COMPLETION
    Shell completion for command line options can be enabled by using the
    --completion flag. The required parameter is the name of your shell.
    Currently supported are: bash, zsh, fish and powershell.

    Detailed instructions:

    Bash:
           source <(tablizer --completion bash)

        To load completions for each session, execute once:

          # Linux:
          $ tablizer --completion bash > /etc/bash_completion.d/tablizer

          # macOS:
          $ tablizer --completion bash > $(brew --prefix)/etc/bash_completion.d/tablizer

    Zsh:
        If shell completion is not already enabled in your environment, you
        will need to enable it. You can execute the following once:

          echo "autoload -U compinit; compinit" >> ~/.zshrc

        To load completions for each session, execute once:

          $ tablizer --completion zsh > "${fpath[1]}/_tablizer"

        You will need to start a new shell for this setup to take effect.

    fish:
           tablizer --completion fish | source

        To load completions for each session, execute once:

           tablizer --completion fish > ~/.config/fish/completions/tablizer.fish

    PowerShell:
           tablizer --completion powershell | Out-String | Invoke-Expression

        To load completions for every new session, run:

           tablizer --completion powershell > tablizer.ps1

        and source this file from your PowerShell profile.

CONFIGURATION AND COLORS
    YOu can put certain configuration values into a configuration file in
    HCL format. By default tablizer looks for
    "$HOME/.config/tablizer/config", but you can provide one using the
    parameter "-f".

    In the configuration the following variables can be defined:

        BG             = "lightGreen"
        FG             = "white"
        HighlightBG    = "lightGreen"
        HighlightFG    = "white"
        NoHighlightBG  = "white"
        NoHighlightFG  = "lightGreen"
        HighlightHdrBG = "red"
        HighlightHdrFG = "white"

    The following color definitions are available:

    black, blue, cyan, darkGray, default, green, lightBlue, lightCyan,
    lightGreen, lightMagenta, lightRed, lightWhite, lightYellow, magenta,
    red, white, yellow

    The Variables FG and BG are being used to highlight matches. The other
    *FG and *BG variables are for colored table output (enabled with the
    "-L" parameter).

    Colorization can be turned off completely either by setting the
    parameter "-N" or the environment variable NO_COLOR to a true value.

BUGS
    In order to report a bug, unexpected behavior, feature requests or to
    submit a patch, please open an issue on github:
    <https://github.com/TLINDEN/tablizer/issues>.

LICENSE
    This software is licensed under the GNU GENERAL PUBLIC LICENSE version
    3.

    Copyright (c) 2022-2024 by Thomas von Dein

    This software uses the following GO modules:

    repr (https://github.com/alecthomas/repr)
        Released under the MIT License, Copyright (c) 2016 Alec Thomas

    cobra (https://github.com/spf13/cobra)
        Released under the Apache 2.0 license, Copyright 2013-2022 The Cobra
        Authors

    dateparse (github.com/araddon/dateparse)
        Released under the MIT License, Copyright (c) 2015-2017 Aaron Raddon

    color (github.com/gookit/color)
        Released under the MIT License, Copyright (c) 2016 inhere

    tablewriter (github.com/olekukonko/tablewriter)
        Released under the MIT License, Copyright (c) 201 by Oleku Konko

    yaml (gopkg.in/yaml.v3)
        Released under the MIT License, Copyright (c) 2006-2011 Kirill
        Simonov

AUTHORS
    Thomas von Dein tom AT vondein DOT org

`
var usage = `

Usage:
  tablizer [regex,...] [file, ...] [flags]

Operational Flags:
  -c, --columns string              Only show the speficied columns (separated by ,)
  -v, --invert-match                select non-matching rows
  -n, --numbering                   Enable header numbering
  -N, --no-color                    Disable pattern highlighting
  -H, --no-headers                  Disable headers display
  -s, --separator string            Custom field separator
  -k, --sort-by int|name            Sort by column (default: 1)
  -z, --fuzzy                       Use fuzzy search [experimental]
  -F, --filter field[!]=reg         Filter given field with regex, can be used multiple times
  -T, --transpose-columns string    Transpose the speficied columns (separated by ,)
  -R, --regex-transposer /from/to/  Apply /search/replace/ regexp to fields given in -T
  -I, --interactive                 Interactively filter and select rows

Output Flags (mutually exclusive):
  -X, --extended                    Enable extended output
  -M, --markdown                    Enable markdown table output
  -O, --orgtbl                      Enable org-mode table output
  -S, --shell                       Enable shell evaluable output
  -Y, --yaml                        Enable yaml output
  -C, --csv                         Enable CSV output
  -A, --ascii                       Default output mode, ascii tabular
  -L, --hightlight-lines            Use alternating background colors for tables
  -y, --yank-columns                Yank specified columns (separated by ,) to clipboard,
                                    space separated

Sort Mode Flags (mutually exclusive):
  -a, --sort-age                    sort according to age (duration) string
  -D, --sort-desc                   Sort in descending order (default: ascending)
  -i, --sort-numeric                sort according to string numerical value
  -t, --sort-time                   sort according to time string

Other Flags:
      --completion <shell>         Generate the autocompletion script for <shell>
  -f, --config <file>              Configuration file (default: ~/.config/tablizer/config)
  -d, --debug                      Enable debugging
  -h, --help                       help for tablizer
  -m, --man                        Display manual page
  -V, --version                    Print program version


`
