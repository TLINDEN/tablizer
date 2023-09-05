## Fixes to be implemented

## Features to be implemented

- add comment support (csf.NewReader().Comment = '#')

- add --no-headers option

### Lisp Plugin Infrastructure using zygo

Hooks:

| Filter    | Purpose                                                     | Args                | Return |
|-----------|-------------------------------------------------------------|---------------------|--------|
| filter    | include or exclude lines                                    | row as hash         | bool   |
| process   | do calculations with data, store results in global lisp env | whole dataset       | nil    |
| transpose | modify a cell                                               | headername and cell | cell   |
| append    | add one or more rows to the dataset (use this to add stats) | nil                 | rows   |
