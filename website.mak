
#
# Makefile for building the website pages with pandoc.
#
# The pages are listed by name, not globbed, so the site does not change with
# whatever .md files happen to sit in the directory (TODO.md and other internal
# notes stay unpublished). opml_test.go checks every listed page has a source.
# README.md becomes index.html.
#
PROJECT = opml

PANDOC = $(shell which pandoc)

PROGRAMS = opml2cmark opml2json opml2urls opmlcat opmlsort urls2opml

SITE_PAGES = README about INSTALL INSTALL_NOTES_macOS INSTALL_NOTES_Windows NOTES user-manual \
	$(addsuffix .1,$(PROGRAMS))

HTML_PAGES = $(patsubst README.html,index.html,$(addsuffix .html,$(SITE_PAGES)))

build: $(HTML_PAGES)

# Pandoc is optional; without it the existing pages are left as they are.
$(HTML_PAGES): .FORCE
	@if [ -z "$(PANDOC)" ]; then echo "website: pandoc not found, skipping $@" >&2; else \
		SRC="$(patsubst index.html,README.html,$@)"; SRC="$${SRC%.html}"; \
		$(PANDOC) --metadata title="$${SRC}" -s --to html5 "$${SRC}.md" -o "$@" \
			--lua-filter=links-to-html.lua \
			--lua-filter=add-col-scope.lua \
			--template=page.tmpl; fi

clean:
	@rm -f $(HTML_PAGES)

.FORCE:
