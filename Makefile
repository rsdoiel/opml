#
# Makefile to compile opmlsort for Mac OS X, Linux, Windows 7
# as well as Raspberry Pi Zero, 1, 2, and 3.
#

build:
	go build -o bin/opmlsort cmds/opmlsort/opmlsort.go

test:
	go test

install:
	env GOBIN=$(HOME)/bin go install cmds/opmlsort/opmlsort.go

status:
	git status

save:
	./mk-website.bash
	git commit -am "quick save"
	git push origin master

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f opml-binary-release.zip ]; then /bin/rm opml-binary-release.zip; fi

release:
	./mk-release.bash	
	
publish:
	./publish.bash
