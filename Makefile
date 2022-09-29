#
# no need to modify anything below
tool    = tablizer
version = $(shell egrep "^var version = " cmd/root.go | cut -d'=' -f2 | cut -d'"' -f 2)
archs   = android darwin freebsd linux netbsd openbsd windows
PREFIX = /usr/local
UID    = root
GID    = 0

all: buildlocal man

man:
	pod2man -c "User Commands" -r 1 -s 1 $(tool).pod > $(tool).1

buildlocal:
	go build

release:
	mkdir -p releases
	$(foreach arch,$(archs), GOOS=$(arch) GOARCH=amd64 go build -x -o releases/$(tool)-$(arch)-amd64-$(version); sha256sum releases/$(tool)-$(arch)-amd64-$(version) | cut -d' ' -f1 > releases/$(tool)-$(arch)-amd64-$(version).sha256sum;)

install: buildlocal
	install -d -o $(UID) -g $(GID) $(PREFIX)/bin
	install -d -o $(UID) -g $(GID) $(PREFIX)/man/man1
	install -o $(UID) -g $(GID) -m 555 $(tool)  $(PREFIX)/sbin/
	install -o $(UID) -g $(GID) -m 444 $(tool).1 $(PREFIX)/man/man1/

clean:
	rm -f $(tool) $(tool).1
