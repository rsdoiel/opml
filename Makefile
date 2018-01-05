#
# Makefile to compile opmlsort for Mac OS X, Linux, Windows 7
# as well as Raspberry Pi Zero, 1, 2, and 3.
#

PROJECT = opml

MOTTO = "An OPML parser package plus opml cat and sort utilities"

VERSION = $(shell grep -m 1 'Version =' opml.go | cut -d\" -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

OS = $(shell uname)

EXT = 
ifeq ($(OS), Windows)
	EXT = .exe
endif


CLI_NAMES = opmlsort opmlcat opml2json

build: $(CLI_NAMES)

opmlsort: bin/opmlsort$(EXT)

opmlcat: bin/opmlcat$(EXT)

opml2json: bin/opml2json$(EXT)

bin/opmlsort$(EXT): opml.go cmd/opmlsort/opmlsort.go
	env CGO_ENABLED=0 go build -o bin/opmlsort$(EXT) cmd/opmlsort/opmlsort.go


bin/opmlcat$(EXT): opml.go cmd/opmlcat/opmlcat.go
	env CGO_ENABLED=0 go build -o bin/opmlcat$(EXT) cmd/opmlcat/opmlcat.go

bin/opml2json$(EXT): opml.go cmd/opml2json/opml2json.go
	env CGO_ENABLED=0 go build -o bin/opml2json$(EXT) cmd/opml2json/opml2json.go

test:
	go test

install:
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmd/opmlsort/opmlsort.go
	env CGO_ENABLES=0 GOBIN=$(HOME)/bin go install cmd/opmlcat/opmlcat.go

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROG)-$(VERSION)-release.zip ]; then /bin/rm $(PROG)-$(VERSION)-release.zip; fi

website:
	./mk-website.bash $(PROJECT) $(MOTTO) $(VERSION)

dist/linux-amd64:
	mkdir -p dist/bin
	env GOOS=linux GOARCH=amd64 go build -o dist/bin/omplsort cmd/opmlsort/opmlsort.go
	env GOOS=linux GOARCH=amd64 go build -o dist/bin/omplcat cmd/opmlcat/opmlcat.go
	env GOOS=linux GOARCH=amd64 go build -o dist/bin/ompl2json cmd/opml2json/opml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md docs
	rm -fR dist/bin

dist/windows-amd64:
	mkdir -p dist/bin
	env GOOS=windows GOARCH=amd64 go build -o dist/bin/opmlsort.exe cmd/opmlsort/opmlsort.go
	env GOOS=windows GOARCH=amd64 go build -o dist/bin/opmlcat.exe cmd/opmlcat/opmlcat.go
	env GOOS=windows GOARCH=amd64 go build -o dist/bin/opml2json.exe cmd/opml2json/opml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md docs
	rm -fR dist/bin

dist/macosx-amd64:
	mkdir -p dist/bin
	env GOOS=darwin GOARCH=amd64 go build -o dist/bin/opmlsort cmd/opmlsort/opmlsort.go
	env GOOS=darwin GOARCH=amd64 go build -o dist/bin/opmlcat cmd/opmlcat/opmlcat.go
	env GOOS=darwin GOARCH=amd64 go build -o dist/bin/opml2json cmd/opml2json/opml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macosx-amd64.zip README.md LICENSE INSTALL.md docs
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/opmlsort cmd/opmlsort/opmlsort.go
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/opmlcat cmd/opmlcat/opmlcat.go
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/opml2json cmd/opml2json/opml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm7.zip README.md LICENSE INSTALL.md docs
	rm -fR dist/bin

dist/raspbian-arm6:
	mkdir -p dist/bin
	env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/bin/opmlsort cmd/opmlsort/opmlsort.go
	env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/bin/opmlcat cmd/opmlcat/opmlcat.go
	env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/bin/opml2json cmd/opml2json/opml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm6.zip README.md LICENSE INSTALL.md docs
	rm -fR dist/bin

distribute_docs:
	mkdir -p dist
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -v docs/opmlsort.md dist/
	cp -v docs/opmlcat.md dist/
	cp -v docs/opml2json.md dist/

release: distribute_docs dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7 dist/linux-arm64

publish:
	./mk-website.bash
	./publish.bash

