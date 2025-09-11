#!/bin/sh

while read LINE; do
    eval "$LINE"; echo "$Customer ordered $Count ${Units}s"
done
