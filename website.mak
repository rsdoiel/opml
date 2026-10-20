
# generated with CMTools 0.0.9 facd557

#
# Makefile for running pandoc on all Markdown docs ending in .md
#
PROJECT = opml

PANDOC = $(shell which pandoc)

MD_PAGES = $(shell ls -1 *.md)

HTML_PAGES = $(shell ls -1 *.md | sed -E 's/.md/.html/g')

build: $(HTML_PAGES) $(MD_PAGES)

$(HTML_PAGES): $(MD_PAGES) .FORCE
	if [ -f $(PANDOC) ]; then $(PANDOC) --metadata title=$(basename $@) -s --to html5 $(basename $@).md -o $(basename $@).html \
		--lua-filter=links-to-html.lua \
	    --template=page.tmpl; fi
	@if [ $@ = "README.html" ]; then mv README.html index.html; fi

clean:
	@rm *.html

.FORCE:
