## tablizer - Manipulate tabular output of other programs

Tablizer  can   be  used   to  re-format   tabular  output   of  other
programs. While you  could do this using standard unix  tools, in some
cases it's a hard job.

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
% kubectl get pods | ./tablizer 
NAME(1)                      READY(2) STATUS(3)  RESTARTS(4)    AGE(5)
repldepl-7bcd8d5b64-7zq4l    1/1      Running    1 (69m ago)    5h26m
repldepl-7bcd8d5b64-m48n8    1/1      Running    1 (69m ago)    5h26m
repldepl-7bcd8d5b64-q2bf4    1/1      Running    1 (69m ago)    5h26m

% kubectl get pods | ./tablizer -c 1,3
NAME(1)                      STATUS(3)
repldepl-7bcd8d5b64-7zq4l    Running
repldepl-7bcd8d5b64-m48n8    Running
repldepl-7bcd8d5b64-q2bf4    Running 
```

Another use case is when the tabular  output is so wide that lines are
being broken and  the whole output is completely distorted.  In such a
case you can use the `-x` flag to get an output similar to `\x` in `psql`:

```
% kubectl get pods | ./tablizer -x    
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
% kubectl get pods | ./tablizer q2bf4
NAME(1)                      READY(2) STATUS(3)  RESTARTS(4)    AGE(5)
repldepl-7bcd8d5b64-q2bf4    1/1      Running    1 (69m ago)    5h26m
```


## Installation

Download the latest release file for your architecture and put it into
a directory within your `$PATH`.

## Getting help

Although I'm happy to hear from tablizer users in private email,
that's the best way for me to forget to do something.

In order to report a bug, unexpected behavior, feature requests
or to submit a patch, please open an issue on github:
https://github.com/TLINDEN/tablizer/issues.

## Copyright and license

This software is licensed under the GNU GENERAL PUBLIC LICENSE version 3.

## Authors

T.v.Dein <tom AT vondein DOT org>

## Project homepage

https://github.com/TLINDEN/tablizer