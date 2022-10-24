## Fixes to be implemented

- highlighting does not repeat, only 1 will be highlighted, see #3

- rm printYamlData() log.Fatal(), maybe return error on all printers?

- refactor parser, there's some duplicate code

## Features to be implemented

- add comment support (csf.NewReader().Comment = '#')

- add output mode csv

- add --no-headers option

-  add input  parsing support  for CSV  including unquoting  of stuff
  like: `"xxx","1919 b"` etc, maybe an extra option for unquoting

