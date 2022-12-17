#
# Makefile to compile opmlsort for Mac OS X, Linux, Windows 7
# as well as Raspberry Pi Zero, 1, 2, and 3.
#

PROJECT = opml

VERSION = $(shell grep '"version":' codemeta.json | cut -d\"  -f 4)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

#PREFIX = /usr/local/bin
PREFIX = $(HOME)

PKGASSETS = $(shell which pkgassets)

PROGRAMS = $(shell ls -1 cmd)

MANPAGES = opmlsort.1 opml2json.1 opmlcat.1 opml2urls.1 urls2opml.1

PANDOC = $(shell which pandoc)

OS = $(shell uname)

EXT =
ifeq ($(OS), Windows)
	EXT = .exe
endif

build: version.go $(PROGRAMS) CITATION.cff

version.go: .FORCE
	@echo "package $(PROJECT)" >version.go
	@echo '' >>version.go
	@echo '// Version of package' >>version.go
	@echo 'const (' >>version.go
	@echo '    Version = `$(VERSION)`' >>version.go
	@echo '' >>version.go
	@echo '    LicenseText = `' >>version.go
	@cat LICENSE >>version.go
	@echo '`' >>version.go
	@echo ')' >>version.go
	@echo '' >>version.go
	@git add version.go

about.md: codemeta.json .FORCE
	cat codemeta.json | sed -E   's/"@context"/"at__context"/g;s/"@type"/"at__type"/g;s/"@id"/"at__id"/g' >_codemeta.json
	if [ -f $(PANDOC) ]; then echo "" | $(PANDOC) --metadata title="About $(PROJECT)" --metadata-file=_codemeta.json --template=codemeta-md.tmpl >about.md; fi

about.html: about.md

CITATION.cff: codemeta.json .FORCE
	cat codemeta.json | sed -E   's/"@context"/"at__context"/g;s/"@type"/"at__type"/g;s/"@id"/"at__id"/g' >_codemeta.json
	if [ -f $(PANDOC) ]; then echo "" | $(PANDOC) --metadata title="Cite $(PROJECT)" --metadata-file=_codemeta.json --template=codemeta-cff.tmpl >CITATION.cff; fi

$(PROGRAMS): cmd/*/*.go $(PACKAGE)
	@mkdir -p bin
	go build -o bin/$@$(EXT) cmd/$@/*.go

test:
	go test

man: .FORCE 
	@mkdir -p man/man1
	@for FNAME in $(MANPAGES); do pandoc $$FNAME.md -s --from markdown --to man >./man/man1/$$FNAME; done


install: build man
	@for FNAME in $(PROGRAMS); do mv bin/$$FNAME $(PREFIX)/bin/; done
	@for FNAME in $(MANPAGES); do cp ./man/man1/$$FNAME $(PREFIX)/man/man1/; done

uninstall: .FORCE
	@for FNAME in $(CLI_NAMES); do if [ -f $(PREFIX)/bin/$$FNAME$(EXT) ]; then rm $(PREFIX)/bin/$$FNAME$(EXT); fi; done
	@for FNAME in $(CLI_NAMES); do if [ -f $(PREFIX)/man/man1/$$FNAME.1 ]; then rm $(PREFIX)/man/man1/$$FNAME.1; fi; done

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

clean:
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi
	if [ -d man ]; then rm -fR man; fi

website: about.html
	make -f website.mak

dist/linux-amd64:
	mkdir -p dist/bin
	env GOOS=linux GOARCH=amd64 go build -o dist/bin/opmlsort cmd/opmlsort/opmlsort.go
	env GOOS=linux GOARCH=amd64 go build -o dist/bin/opmlcat cmd/opmlcat/opmlcat.go
	env GOOS=linux GOARCH=amd64 go build -o dist/bin/opml2json cmd/opml2json/opml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md docs/* bin/*
	rm -fR dist/bin

dist/windows-amd64:
	mkdir -p dist/bin
	env GOOS=windows GOARCH=amd64 go build -o dist/bin/opmlsort.exe cmd/opmlsort/opmlsort.go
	env GOOS=windows GOARCH=amd64 go build -o dist/bin/opmlcat.exe cmd/opmlcat/opmlcat.go
	env GOOS=windows GOARCH=amd64 go build -o dist/bin/opml2json.exe cmd/opml2json/opml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md docs/* bin/*
	rm -fR dist/bin

dist/windows-arm64:
	mkdir -p dist/bin
	env GOOS=windows GOARCH=arm64 go build -o dist/bin/opmlsort.exe cmd/opmlsort/opmlsort.go
	env GOOS=windows GOARCH=arm64 go build -o dist/bin/opmlcat.exe cmd/opmlcat/opmlcat.go
	env GOOS=windows GOARCH=arm64 go build -o dist/bin/opml2json.exe cmd/opml2json/opml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-arm64.zip README.md LICENSE INSTALL.md docs/* bin/*
	rm -fR dist/bin

dist/macos-amd64:
	mkdir -p dist/bin
	env GOOS=darwin GOARCH=amd64 go build -o dist/bin/opmlsort cmd/opmlsort/opmlsort.go
	env GOOS=darwin GOARCH=amd64 go build -o dist/bin/opmlcat cmd/opmlcat/opmlcat.go
	env GOOS=darwin GOARCH=amd64 go build -o dist/bin/opml2json cmd/opml2json/opml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macos-amd64.zip README.md LICENSE INSTALL.md docs/* bin/*
	rm -fR dist/bin

dist/macos-arm64:
	mkdir -p dist/bin
	env GOOS=darwin GOARCH=arm64 go build -o dist/bin/opmlsort cmd/opmlsort/opmlsort.go
	env GOOS=darwin GOARCH=arm64 go build -o dist/bin/opmlcat cmd/opmlcat/opmlcat.go
	env GOOS=darwin GOARCH=arm64 go build -o dist/bin/opml2json cmd/opml2json/opml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macos-arm64.zip README.md LICENSE INSTALL.md docs/* bin/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/opmlsort cmd/opmlsort/opmlsort.go
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/opmlcat cmd/opmlcat/opmlcat.go
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/opml2json cmd/opml2json/opml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm7.zip README.md LICENSE INSTALL.md docs/* bin/*
	rm -fR dist/bin

dist/raspbian-arm6:
	mkdir -p dist/bin
	env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/bin/opmlsort cmd/opmlsort/opmlsort.go
	env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/bin/opmlcat cmd/opmlcat/opmlcat.go
	env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/bin/opml2json cmd/opml2json/opml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm6.zip README.md LICENSE INSTALL.md docs/* bin/*
	rm -fR dist/bin

dist/linux-arm64:
	mkdir -p dist/bin
	env GOOS=linux GOARCH=arm64 GOARM=6 go build -o dist/bin/opmlsort cmd/opmlsort/opmlsort.go
	env GOOS=linux GOARCH=arm64 GOARM=6 go build -o dist/bin/opmlcat cmd/opmlcat/opmlcat.go
	env GOOS=linux GOARCH=arm64 GOARM=6 go build -o dist/bin/opml2json cmd/opml2json/opml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-arm64.zip README.md LICENSE INSTSALL.md docs/* bin/*
	rm -fR dist/bin

generate_usage_pages: opmlsort opmlcat opml2json
	bash gen-usage-pages.bash

distribute_docs:
	mkdir -p dist/docs
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	bash gen-usage-pages.bash
	cp -v docs/opmlsort.md dist/docs/
	cp -v docs/opmlcat.md dist/docs/
	cp -v docs/opml2json.md dist/docs/

release: generate_usage_pages distribute_docs dist/linux-amd64 dist/windows-amd64 dist/macos-amd64 dist/macos-arm64 dist/raspbian-arm7 dist/raspbian-arm6 dist/linux-arm64

publish: website
	./publish.bash


.FORCE:
