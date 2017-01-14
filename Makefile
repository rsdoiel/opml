#
# Makefile to compile opmlsort for Mac OS X, Linux, Windows 7
# as well as Raspberry Pi Zero, 1, 2, and 3.
#

PROJECT = opml

VERSION = $(shell grep -m 1 'Version =' opml.go | cut -d\" -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

CLI_NAMES = opmlsort opmlcat

build: $(CLI_NAMES)

opmlsort: bin/opmlsort

opmlcat: bin/opmlcat

bin/opmlsort: opml.go cmds/opmlsort/opmlsort.go
	env CGO_ENABLED=0 go build -o bin/opmlsort cmds/opmlsort/opmlsort.go


bin/opmlcat: opml.go cmds/opmlcat/opmlcat.go
	env CGO_ENABLED=0 go build -o bin/opmlcat cmds/opmlcat/opmlcat.go

test:
	go test

install:
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/opmlsort/opmlsort.go
	env CGO_ENABLES=0 GOBIN=$(HOME)/bin go install cmds/opmlcat/opmlcat.go

status:
	git status

save:
	git commit -am "Quick Save"
	git push origin $(BRANCH)

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROG)-$(VERSION)-release.zip ]; then /bin/rm $(PROG)-$(VERSION)-release.zip; fi

website:
	./mk-website.bash

publish:
	./publish.bash

dist/linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/omplsort cmds/opmlsort/opmlsort.go
	env GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/omplcat cmds/opmlcat/opmlcat.go

dist/windows-amd64:
	env GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/opmlsort.exe cmds/opmlsort/opmlsort.go
	env GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/opmlcat.exe cmds/opmlcat/opmlcat.go

dist/macosx-amd64:
	env GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/opmlsort cmds/opmlsort/opmlsort.go
	env GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/opmlcat cmds/opmlcat/opmlcat.go

dist/raspbian-arm7:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/opmlsort cmds/opmlsort/opmlsort.go
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/opmlcat cmds/opmlcat/opmlcat.go

dist/raspbian-arm6:
	env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspberrypi-arm6/opmlsort cmds/opmlsort/opmlsort.go
	env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspberrypi-arm6/opmlcat cmds/opmlcat/opmlcat.go

release: dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7 dist/raspbian-arm6
	cp -v README.md dist/
	cp -v INSTALL.md dist/
	cp -v LICENSE dist/
	zip -r $(PROJECT)-$(VERSION)-release.zip dist/*
	
