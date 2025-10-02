package cmd

const shortusage = `tablizer [regex,...] [-r file] [flags]
-c col,...   show specified columns                -L  highlight matching lines
-k col,...   sort by specified columns             -j  read JSON input
-F col=reg   filter field with regexp              -v  invert match
-T col,...   transpose specified columns           -n  numberize columns
-R /from/to/ apply replacement to columns in -T    -N  do not use colors
-y col,...   yank columns to clipboard             -H  do not show headers
--ofs char   output field separator	               -s  specify field separator
-r file      read input from file                  -z  use fuzzy search
-f file      read config from file                 -I  interactive filter mode
                                                   -d  debug
-O org -C CSV -M md -X ext -S shell -Y yaml        -D  sort descending order
-m  show manual       --help  show detailed help   -v  show version
-a  sort by age       -i      sort numerically     -t  sort by time`
