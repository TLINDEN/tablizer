
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


#
# no need to modify anything below
tool    = tablizer
version = $(shell egrep "= .v" lib/common.go | cut -d'=' -f2 | cut -d'"' -f 2)
archs   = android darwin freebsd linux netbsd openbsd windows
PREFIX  = /usr/local
UID     = root
GID     = 0
BRANCH  = $(shell git describe --all | cut -d/ -f2)
COMMIT  = $(shell git rev-parse --short=8 HEAD)
BUILD   = $(shell date +%Y.%m.%d.%H%M%S) 
VERSION:= $(if $(filter $(BRANCH), development),$(version)-$(BRANCH)-$(COMMIT)-$(BUILD),$(version))


all: $(tool).1 cmd/$(tool).go buildlocal

%.1: %.pod
	pod2man -c "User Commands" -r 1 -s 1 $*.pod > $*.1

cmd/%.go: %.pod
	echo "package cmd" > cmd/$*.go
	echo >> cmd/$*.go
	echo "var manpage = \`" >> cmd/$*.go
	pod2text $*.pod >> cmd/$*.go
	echo "\`" >> cmd/$*.go

buildlocal:
	go build -ldflags "-X 'github.com/tlinden/tablizer/lib.VERSION=$(VERSION)'"

release:
	./mkrel.sh $(tool) $(version)
	gh release create $(version) --generate-notes releases/*

install: buildlocal
	install -d -o $(UID) -g $(GID) $(PREFIX)/bin
	install -d -o $(UID) -g $(GID) $(PREFIX)/man/man1
	install -o $(UID) -g $(GID) -m 555 $(tool)  $(PREFIX)/sbin/
	install -o $(UID) -g $(GID) -m 444 $(tool).1 $(PREFIX)/man/man1/

clean:
	rm -rf $(tool) releases

test:
	go test -v ./...
