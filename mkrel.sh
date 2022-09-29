#!/bin/bash

# get list with: go tool dist list
DIST="darwin/amd64
freebsd/amd64
linux/amd64
netbsd/amd64
openbsd/amd64
windows/amd64"

tool="$1"
version="$2"

rm -rf releases
mkdir -p releases


for D in $DIST; do
    os=${D/\/*/}
    arch=${D/*\//}
    binfile="releases/${tool}-${os}-${arch}-${version}"
    tardir="${tool}-${os}-${arch}-${version}"
    tarfile="releases/${tool}-${os}-${arch}-${version}.tar.gz"
    set -x
    GOOS=${os} GOARCH=${arch} go build -o ${binfile}
    mkdir -p ${tardir}
    cp ${binfile} README.md LICENSE ${tardir}/
    echo 'tool = tablizer
PREFIX = /usr/local
UID    = root
GID    = 0

install:
	install -d -o $(UID) -g $(GID) $(PREFIX)/bin
	install -d -o $(UID) -g $(GID) $(PREFIX)/man/man1
	install -o $(UID) -g $(GID) -m 555 $(tool)  $(PREFIX)/sbin/
	install -o $(UID) -g $(GID) -m 444 $(tool).1 $(PREFIX)/man/man1/' > ${tardir}/Makefile
    tar cpzf ${tarfile} ${tardir}
    sha256sum ${binfile} | cut -d' ' -f1 > ${binfile}.sha256
    sha256sum ${tarfile} | cut -d' ' -f1 > ${tarfile}.sha256
    rm -rf ${tardir}
    set +x
done

