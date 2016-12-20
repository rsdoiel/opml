#
# Makefile to compile opmlsort for Mac OS X, Linux, Windows 7
# as well as Raspberry Pi Zero, 1, 2, and 3.
#

PROJECT = opml

VERSION = $(shell grep 'Version =' opml.go | cut -d\"  -f 2)

BRANCH = $(shell git branch | grep '*' | cut -d \  -f 2)

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
	git commit -am "quick save"
	git push origin $(BRANCH)

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROG)-$(VERSION)-release.zip ]; then /bin/rm $(PROG)-$(VERSION)-release.zip; fi

release:
	./mk-release.bash	
	
publish:
	./publish.bash
