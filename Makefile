#
# no need to modify anything below
tool    = tablizer
version = $(shell egrep "^var version = " cmd/root.go | cut -d'=' -f2 | cut -d'"' -f 2)
archs   = android darwin freebsd linux netbsd openbsd windows

all:
	@echo "Type 'make install' to install $(tool)"

install:
	install -m 755 -d $(bindir)
	install -m 755 -d $(linkdir)
	install -m 755 $(tool) $(bindir)/$(tool)-$(version)
	ln -sf $(bindir)/$(tool)-$(version) $(linkdir)/$(tool)

release:
	mkdir -p releases
	$(foreach arch,$(archs), GOOS=$(arch) GOARCH=amd64 go build -x -o releases/$(tool)-$(arch)-amd64-$(version); sha256sum releases/$(tool)-$(arch)-amd64-$(version) | cut -d' ' -f1 > releases/$(tool)-$(arch)-amd64-$(version).sha256sum;)
