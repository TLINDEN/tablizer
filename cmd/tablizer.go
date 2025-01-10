package cmd

var manpage = `
NAME
    tablizer - Manipulate tabular output of other programs

SYNOPSIS
        Usage:
          tablizer [regex] [file, ...] [flags]
    
        Operational Flags:
          -c, --columns string     Only show the speficied columns (separated by ,)
          -v, --invert-match       select non-matching rows
          -n, --no-numbering       Disable header numbering
          -N, --no-color           Disable pattern highlighting
          -H, --no-headers         Disable headers display
          -s, --separator string   Custom field separator
          -k, --sort-by int        Sort by column (default: 1)
          -z, --fuzzy              Use fuzzy search [experimental]
          -F, --filter field=reg   Filter given field with regex, can be used multiple times

        Output Flags (mutually exclusive):
          -X, --extended           Enable extended output
          -M, --markdown           Enable markdown table output
          -O, --orgtbl             Enable org-mode table output
          -S, --shell              Enable shell evaluable output
          -Y, --yaml               Enable yaml output
          -C, --csv                Enable CSV output
          -A, --ascii              Default output mode, ascii tabular
          -L, --hightlight-lines   Use alternating background colors for tables

        Sort Mode Flags (mutually exclusive):
          -a, --sort-age           sort according to age (duration) string
          -D, --sort-desc          Sort in descending order (default: ascending)
          -i, --sort-numeric       sort according to string numerical value
          -t, --sort-time          sort according to time string

        Other Flags:
              --completion <shell> Generate the autocompletion script for <shell>
          -f, --config <file>      Configuration file (default: ~/.config/tablizer/config)
          -l, --load-path <path>   Load path for lisp plugins (expects *.zy files)
          -d, --debug              Enable debugging
          -h, --help               help for tablizer
          -m, --man                Display manual page
          -V, --version            Print program version

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
    (as in GNU sort(1)). The default sort column is the first one. To
    disable sorting at all, supply 0 (Zero) to -k. The default sort order is
    ascending. You can change this to descending order using the option -D.
    The default sort order is by string, but there are other sort modes:

    -a --sort-age
        Sorts duration strings like "1d4h32m51s".

    -i --sort-numeric
        Sorts numeric fields.

    -t --sort-time
        Sorts timestamps.

    Finally the -d option enables debugging output which is mostly useful
    for the developer.

  PATTERNS AND FILTERING
    You can reduce the rows being displayed by using a regular expression
    pattern. The regexp is PCRE compatible, refer to the syntax cheat sheet
    here: <https://github.com/google/re2/wiki/Syntax>. If you want to read a
    more comprehensive documentation about the topic and have perl installed
    you can read it with:

        perldoc perlre

    Or read it online: <https://perldoc.perl.org/perlre>.

    A note on modifiers: the regexp engine used in tablizer uses another
    modifier syntax:

        (?MODIFIER)

    The most important modifiers are:

    "i" ignore case "m" multiline mode "s" single line mode

    Example for a case insensitive search:

        kubectl get pods -A | tablizer "(?i)account"

    You can use the experimental fuzzy search feature by providing the
    option -z, in which case the pattern is regarded as a fuzzy search term,
    not a regexp.

    Sometimes you want to filter by one or more columns. You can do that
    using the -F option. The option can be specified multiple times and has
    the following format:

        fieldname=regexp

    Fieldnames (== columns headers) are case insensitive.

    If you specify more than one filter, both filters have to match (AND
    operation).

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

  ENVIRONMENT VARIABLES
    tablizer supports certain environment variables which use can use to
    influence program behavior. Commandline flags have always precedence
    over environment variables.

    <T_NO_HEADER_NUMBERING> - disable numbering of header fields, like -n.
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

LISP PLUGINS [experimental]
    Tablizer supports plugins written in zygomys lisp. You can supply a
    directory to the "-l" parameter containing *.zy files or a single .zy
    file containing lisp code.

    You can put as much code as you want into the file, but you need to add
    one lips function to a hook at the end.

    The following hooks are available:

    filter
        The filter hook works one a whole line of the input. Your hook
        function is expected to return true or false. If you return true,
        the line will be included in the output, otherwise not.

        Multiple filter hook functions are supported.

        Example:

            /*
            Simple filter hook function. Splits the argument by whitespace,
            fetches the 2nd element, converts it to an int and returns true
            if it s larger than 5, false otherwise.
            */
            (defn uselarge [line]
              (cond (> (atoi (second (resplit line `\s+`))) 5) true false))
    
            /* Register the filter hook */
            (addhook %filter %uselarge)

    process
        The process hook function gets a table containing the parsed input
        data (see "lib/common.go:type Tabdata struct". It is expected to
        return a pair containing a bool to denote if the table has been
        modified, and the [modified] table. The resulting table may have
        less rows than the original and cells may have changed content but
        the number of columns must persist.

    transpose
        not yet implemented.

    append
        not yet implemented.

    Beside the existing language features, the following additional lisp
    functions are provided by tablizer:

        (resplit [string, regex]) => list
        (atoi    [string])        => int
        (matchre [string, regex]) => bool

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
  tablizer [regex] [file, ...] [flags]

Operational Flags:
  -c, --columns string     Only show the speficied columns (separated by ,)
  -v, --invert-match       select non-matching rows
  -n, --no-numbering       Disable header numbering
  -N, --no-color           Disable pattern highlighting
  -H, --no-headers         Disable headers display
  -s, --separator string   Custom field separator
  -k, --sort-by int        Sort by column (default: 1)
  -z, --fuzzy              Use fuzzy search [experimental]
  -F, --filter field=reg   Filter given field with regex, can be used multiple times

Output Flags (mutually exclusive):
  -X, --extended           Enable extended output
  -M, --markdown           Enable markdown table output
  -O, --orgtbl             Enable org-mode table output
  -S, --shell              Enable shell evaluable output
  -Y, --yaml               Enable yaml output
  -C, --csv                Enable CSV output
  -A, --ascii              Default output mode, ascii tabular
  -L, --hightlight-lines   Use alternating background colors for tables

Sort Mode Flags (mutually exclusive):
  -a, --sort-age           sort according to age (duration) string
  -D, --sort-desc          Sort in descending order (default: ascending)
  -i, --sort-numeric       sort according to string numerical value
  -t, --sort-time          sort according to time string

Other Flags:
      --completion <shell> Generate the autocompletion script for <shell>
  -f, --config <file>      Configuration file (default: ~/.config/tablizer/config)
  -l, --load-path <path>   Load path for lisp plugins (expects *.zy files)
  -d, --debug              Enable debugging
  -h, --help               help for tablizer
  -m, --man                Display manual page
  -V, --version            Print program version


`
