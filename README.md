[![Actions](https://github.com/tlinden/tablizer/actions/workflows/ci.yaml/badge.svg)](https://github.com/tlinden/tablizer/actions)
[![License](https://img.shields.io/badge/license-GPL-blue.svg)](https://github.com/tlinden/tablizer/blob/master/LICENSE)

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
% kubectl get pods | tablizer 
NAME(1)                      READY(2) STATUS(3)  RESTARTS(4)    AGE(5)
repldepl-7bcd8d5b64-7zq4l    1/1      Running    1 (69m ago)    5h26m
repldepl-7bcd8d5b64-m48n8    1/1      Running    1 (69m ago)    5h26m
repldepl-7bcd8d5b64-q2bf4    1/1      Running    1 (69m ago)    5h26m

% kubectl get pods | tablizer -c 1,3
NAME(1)                      STATUS(3)
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
NAME(1)                      READY(2) STATUS(3)  RESTARTS(4)    AGE(5)
repldepl-7bcd8d5b64-q2bf4    1/1      Running    1 (69m ago)    5h26m
```

There are more output modes like org-mode (orgtbl) and markdown.

## Installation

There are multiple ways to install **tablizer**:

- Go to the [latest release page](https://github.com/muesli/mango/releases/latest),
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
automatically installed if  you install from source.  However, you can
read the man-page online:

https://github.com/TLINDEN/tablizer/blob/main/tablizer.pod

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

## Copyright and license

This software is licensed under the GNU GENERAL PUBLIC LICENSE version 3.

## Authors

T.v.Dein <tom AT vondein DOT org>

## Project homepage

https://github.com/TLINDEN/tablizer
