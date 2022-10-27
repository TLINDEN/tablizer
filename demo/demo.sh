#!/bin/bash

prompt() {
    if test -n "$1"; then
        echo
        echo -n "% $*"
        sleep 1
        echo
        $*
        echo
        echo -n "% "
    else
        echo -n "% "
    fi
}

PATH=..:$PATH
clear
while IFS=$'\t' read -r flags table msg source _; do
    echo "#"
    echo "#   source tabular data:"
    cat $table
    echo
    echo "#"
    echo "#   $msg:"
    prompt "tablizer $flags $table"
  
    sleep 4
    clear
done < <(yq -r tables.yaml \
         | yq -r '.tables[] | [.flags, .table, .msg, .source] | @tsv')
