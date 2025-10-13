[![Actions](https://github.com/tlinden/tablizer/actions/workflows/ci.yaml/badge.svg)](https://github.com/tlinden/tablizer/actions)
[![License](https://img.shields.io/badge/license-GPL-blue.svg)](https://github.com/tlinden/tablizer/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tlinden/tablizer)](https://goreportcard.com/report/github.com/tlinden/tablizer)

## tablizer - Manipulate tabular output of other programs

Tablizer  can   be  used   to  re-format   tabular  output   of  other
programs. While you  could do this using standard unix  tools, in some
cases it's a hard job. With tablizer you can filter by column[s],
ignore certain column[s] by regex, name or number. It can output the
tabular data in a range of formats (see below). There's even an
interactive filter/selection tool available.

## FEATURES

- supports csv, json or ascii format input from files or stdin
- split any tabular input data by character or regular expression into columns
- add headers if input data doesn't contain them (automatically or manually)
- print tabular data as ascii table, org-mode, markdown, csv, shell-evaluable or yaml format
- filter rows by regular expression (saves a call to `| grep ...`)
- filter rows by column filter
- filters may also be negations eg `-Fname!=cow.*` or `-v`
- modify cells wih regular expressions
- reduce columns by specifying which columns to show, with regex support
- color support
- sort by any field[s], multiple sort modes are supported
- shell completion for options
- regular used options can be put into a config file
- filter TUI where where you can interactively sort and filter rows

## Demo

![demo cast](vhsdemo/demo.gif)


## Usage

```default
Usage:
  tablizer [regex,...] [file, ...] [flags]

Operational Flags:
  -c, --columns string               Only show the speficied columns (separated by ,)
  -v, --invert-match                 select non-matching rows
  -n, --numbering                    Enable header numbering
  -N, --no-color                     Disable pattern highlighting
  -H, --no-headers                   Disable headers display
  -s, --separator <string>           Custom field separator
  -k, --sort-by <int|name>           Sort by column (default: 1)
  -z, --fuzzy                        Use fuzzy search [experimental]
  -F, --filter <field[!]=reg>        Filter given field with regex, can be used multiple times
  -T, --transpose-columns string     Transpose the speficied columns (separated by ,)
  -R, --regex-transposer </from/to/> Apply /search/replace/ regexp to fields given in -T
  -j, --json                         Read JSON input (must be array of hashes)
  -I, --interactive                  Interactively filter and select rows
      --auto-headers                 Generate headers if there are none present in input
      --custom-headers a,b,...       Use custom headers, separated by comma


Output Flags (mutually exclusive):
  -X, --extended                     Enable extended output
  -M, --markdown                     Enable markdown table output
  -O, --orgtbl                       Enable org-mode table output
  -S, --shell                        Enable shell evaluable output
  -Y, --yaml                         Enable yaml output
  -C, --csv                          Enable CSV output
  -A, --ascii                        Default output mode, ascii tabular
  -L, --hightlight-lines             Use alternating background colors for tables
  -y, --yank-columns                 Yank specified columns (separated by ,) to clipboard,
                                     space separated
      --ofs <char>                   Output field separator, used by -A and -C. 

Sort Mode Flags (mutually exclusive):
  -a, --sort-age                     sort according to age (duration) string
  -D, --sort-desc                    Sort in descending order (default: ascending)
  -i, --sort-numeric                 sort according to string numerical value
  -t, --sort-time                    sort according to time string

Other Flags:
  -r  --read-file <file>             Use <file> as input instead of STDIN
      --completion <shell>           Generate the autocompletion script for <shell>
  -f, --config <file>                Configuration file (default: ~/.config/tablizer/config)
  -d, --debug                        Enable debugging
  -h, --help                         help for tablizer
  -m, --man                          Display manual page
  -V, --version                      Print program version
```

Let's take this output:
```
% kubectl get pods -o wide
NAME                        READY   STATUS    RESTARTS      AGE
repldepl-7bcd8d5b64-7zq4l   1/1     Running   1 (69m ago)   5h26m
repldepl-7bcd8d5b64-m48n8   1/1     Running   1 (69m ago)   5h26m
repldepl-7bcd8d5b64-q2bf4   1/1     Running   1 (69m ago)   5h26m
```

But you're only interested in the  NAME and STATUS columns. Here's how
to do this with tablizer:

```
% kubectl get pods | tablizer 
NAME                         READY    STATUS     RESTARTS       AGE
repldepl-7bcd8d5b64-7zq4l    1/1      Running    1 (69m ago)    5h26m
repldepl-7bcd8d5b64-m48n8    1/1      Running    1 (69m ago)    5h26m
repldepl-7bcd8d5b64-q2bf4    1/1      Running    1 (69m ago)    5h26m

% kubectl get pods | tablizer -c 1,3
NAME                         STATUS
repldepl-7bcd8d5b64-7zq4l    Running
repldepl-7bcd8d5b64-m48n8    Running
repldepl-7bcd8d5b64-q2bf4    Running 
```

Another use case is when the tabular  output is so wide that lines are
being broken and the whole output  is completely distorted.  In such a
case you can use the `-o extended  | -X` flag to get an output similar
to `\x` in `psql`:

```
% kubectl get pods | tablizer -X
    NAME: repldepl-7bcd8d5b64-7zq4l
   READY: 1/1
  STATUS: Running
RESTARTS: 1 (71m ago)
     AGE: 5h28m

    NAME: repldepl-7bcd8d5b64-m48n8
   READY: 1/1
  STATUS: Running
RESTARTS: 1 (71m ago)
     AGE: 5h28m

    NAME: repldepl-7bcd8d5b64-q2bf4
   READY: 1/1
  STATUS: Running
RESTARTS: 1 (71m ago)
     AGE: 5h28m
```

Tablize can read one or more files or - if none specified - from STDIN.

You can also specify a regex pattern to reduce the output:

```
% kubectl get pods | tablizer q2bf4
NAME                         READY    STATUS     RESTARTS       AGE
repldepl-7bcd8d5b64-q2bf4    1/1      Running    1 (69m ago)    5h26m
```

Sometimes a filter regex is to broad  and you wish to filter only on a
particular column. This is possible using `-F`:
```
% kubectl get pods | tablizer -Fname=2
NAME                            READY   STATUS  RESTARTS        AGE
repldepl-7bcd8d5b64-q2bf4       1/1     Running 1 (69m ago)     5h26m
```

Here we filtered  the `NAME` column for `2`, which  would have matched
otherwise on all rows.

There are more output modes like org-mode (orgtbl) and markdown.

You can also use it to modify certain cells using regular expression
matching. For example:

```shell
kubectl get pods | tablizer -T4 -R '/ /-/'
NAME                            READY   STATUS  RESTARTS        AGE
repldepl-7bcd8d5b64-7zq4l       1/1     Running 1-(69m-ago)     5h26m
repldepl-7bcd8d5b64-m48n8       1/1     Running 1-(69m-ago)     5h26m
repldepl-7bcd8d5b64-q2bf4       1/1     Running 1-(69m-ago)     5h26m
```

Here, we modified the 4th column (`-T4`) by replacing every space with
a dash. If you need to work with `/` characters, you can also use any
other separator, for instance: `-R '| |-|'`.

There's also an interactive mode, invoked with the option B<-I>, where
you can interactively filter and select rows:

<img width="937" height="293" alt="interactive" src="https://github.com/user-attachments/assets/0d4d65e2-d156-43ed-8021-39047c7939ed" />



## Installation

There are multiple ways to install **tablizer**:

- Go to the [latest release page](https://github.com/tlinden/tablizer/releases/latest),
  locate the binary for your operating system and platform.
  
  Download it and put it into some directory within your `$PATH` variable.
  
- The release page also contains a tarball for every supported platform. Unpack it
  to some temporary directory, extract it and execute the following command inside:
  ```
  sudo make install
  ```
  
- You can also install from source. Issue the following commands in your shell:
  ```
  git clone https://github.com/TLINDEN/tablizer.git
  cd tablizer
  make
  sudo make install
  ```

If you  do not find a  binary release for your  platform, please don't
hesitate to ask me about it, I'll add it.

## Documentation

The  documentation  is  provided  as  a unix  man-page.   It  will  be
automatically installed if  you install from source.

[However, you can read the man-page online](https://github.com/TLINDEN/tablizer/blob/main/tablizer.pod).

Or if you cloned  the repository you can read it  this way (perl needs
to be installed though): `perldoc tablizer.pod`.

If you have the binary installed, you  can also read the man page with
this command:

    tablizer --man

## Getting help

Although I'm happy to hear from tablizer users in private email,
that's the best way for me to forget to do something.

In order to report a bug, unexpected behavior, feature requests
or to submit a patch, please open an issue on github:
https://github.com/TLINDEN/tablizer/issues.

## Prior Art

When I started with tablizer I was not aware that other tools
exist. Here is a non-exhausive list of the ones I find especially
awesome:

### [miller](https://github.com/johnkerl/miller)

This is a really powerful tool to work with tabular data and it also
allows other inputs as json, csv etc. You can filter, manipulate,
create pipelines, there's even a programming language builtin to do
even more amazing things.

### [csvq](https://github.com/mithrandie/csvq)

Csvq allows you to query CSV and TSV data using SQL queries. How nice
is that? Highly recommended if you have to work with a large (and
wide) dataset and need to apply a complicated set of rules.

### [goawk](https://github.com/benhoyt/goawk)

Goawk is a 100% POSIX compliant AWK implementation in GO, which also
supports CSV and TSV data as input (using `-i csv` for example). You
can apply any kind of awk code to your tabular data, there are no
limit to your creativity!

### [teip](https://github.com/greymd/teip)

I particularly like teip, it's a real gem. You can use it to drill
"holes" into your tabular data and modify these "holes" using small
external unix commands such as grep or sed. The possibilities are
endless, you can even use teip to modify data inside a hole created by
teip. Highly recommended.


## Copyright and license

This software is licensed under the GNU GENERAL PUBLIC LICENSE version 3.

## Authors

T.v.Dein <tom AT vondein DOT org>

## Project homepage

https://github.com/TLINDEN/tablizer
