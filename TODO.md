## Fixes to be implemented

- catch  rows with less  fields as headers and  fill them up  to avoid
  a panic

## Features to be implemented

- ability to change sort order (ascending vs descending)

- sorting by: numerical, time, duration, string(default)

- allow  regexp in -c like:  `-c N.*,ST,8` which means,  match "NAME",
  "NAMESPACE", "STATUS", 8th Header
  
- add output modes yaml and csv

- add --no-headers option

-  add input  parsing support  for CSV  including unquoting  of stuff
  like: `"xxx","1919 b"` etc, maybe an extra option for unquoting

