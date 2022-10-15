package cmd

var manpage = `
NAME
    tablizer - Manipulate tabular output of other programs

SYNOPSIS
        Usage:
          tablizer [regex] [file, ...] [flags]
    
        Flags:
          -c, --columns string     Only show the speficied columns (separated by ,)
          -d, --debug              Enable debugging
          -h, --help               help for tablizer
          -v, --invert-match       select non-matching rows
          -m, --man                Display manual page
          -n, --no-numbering       Disable header numbering
          -N, --no-color           Disable pattern highlighting
          -o, --output string      Output mode - one of: orgtbl, markdown, extended, ascii(default)
          -X, --extended           Enable extended output
          -M, --markdown           Enable markdown table output
          -O, --orgtbl             Enable org-mode table output
          -s, --separator string   Custom field separator
          -a, --sort-age           sort according to age (duration) string
          -k, --sort-by int        Sort by column (default: 1)
          -D, --sort-desc          Sort in descending order (default: ascending)
          -i, --sort-numeric       sort according to string numerical value
          -t, --sort-time          sort according to time string
          -v, --version            Print program version

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

  PATTERNS
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
    markdown which prints a Markdown table.

BUGS
    In order to report a bug, unexpected behavior, feature requests or to
    submit a patch, please open an issue on github:
    <https://github.com/TLINDEN/tablizer/issues>.

LICENSE
    This software is licensed under the GNU GENERAL PUBLIC LICENSE version
    3.

    Copyright (c) 2022 by Thomas von Dein

    This software uses the following GO libraries:

    repr (https://github.com/alecthomas/repr)
        Released under the MIT License, Copyright (c) 2016 Alec Thomas

    cobra (https://github.com/spf13/cobra)
        Released under the Apache 2.0 license, Copyright 2013-2022 The Cobra
        Authors

AUTHORS
    Thomas von Dein tom AT vondein DOT org

`
