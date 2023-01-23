#!/bin/bash

# Copyright © 2022 Thomas von Dein

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.


# get list with: go tool dist list
DIST="darwin/amd64
freebsd/amd64
linux/amd64
netbsd/amd64
openbsd/amd64
windows/amd64"

tool="$1"
version="$2"

if test -z "$version"; then
  echo "Usage: $0 <tool name> <release version>"
  exit 1
fi

rm -rf releases
mkdir -p releases


for D in $DIST; do
    os=${D/\/*/}
    arch=${D/*\//}
    binfile="releases/${tool}-${os}-${arch}-${version}"
    tardir="${tool}-${os}-${arch}-${version}"
    tarfile="releases/${tool}-${os}-${arch}-${version}.tar.gz"
    set -x
    GOOS=${os} GOARCH=${arch} go build -o ${binfile} -ldflags "-X 'github.com/tlinden/tablizer/cfg.VERSION=${version}'"
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

