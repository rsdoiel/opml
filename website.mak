
#
# Makefile for the opml website: pandoc builds an HTML page beside each Markdown
# page, the Markdown pages are published as they are, for BonBon browsers and
# cmarkbrowser, and toc.bash writes /.well-known/toc.opml.
#
# The pages are listed by name, not globbed, so the site does not change with
# whatever .md files happen to sit in the directory. site_test.go checks that
# every published page is listed and that the built files are current.
# README.md becomes index.html and index.md.
#
# Written by hand from cmarkbrowser's website.mak; do not regenerate it with cmt.
#
PROJECT = opml

# Where the site is published (GitHub Pages, gh-pages branch).
BASE_URL = https://rsdoiel.github.io/opml/

PANDOC = $(shell which pandoc)

PROGRAMS = cmark2opml opml2cmark opml2json opml2urls opmlcat opmlsort urls2opml

# Every published page, in the order of the table of contents, README.md first:
# it is the home page. TODO.md and index.md (a copy of README.md) are not pages.
PAGE_ORDER = README.md user-manual.md cmark2opml.1.md opml2cmark.1.md opml2json.1.md opml2urls.1.md opmlcat.1.md opmlsort.1.md urls2opml.1.md INSTALL.md INSTALL_NOTES_macOS.md INSTALL_NOTES_Windows.md NOTES.md CHANGES.md about.md

NO_PUBLISH = TODO.md index.md

MD_PAGES = $(filter-out $(NO_PUBLISH),$(shell ls -1 *.md))

MISSING = $(filter-out $(PAGE_ORDER),$(MD_PAGES))

HTML_PAGES = $(patsubst README.html,index.html,$(patsubst %.md,%.html,$(PAGE_ORDER)))

TOC = .well-known/toc.opml

build: check-order $(HTML_PAGES) index.md $(TOC) .nojekyll

check-order:
	@if [ -n "$(MISSING)" ]; then echo "website.mak: add to PAGE_ORDER: $(MISSING)" >&2; exit 1; fi

# Pandoc is optional; without it the existing pages are left as they are.
# Each HTML page names its Markdown page, for a reader that finds the HTML first.
$(HTML_PAGES): .FORCE
	@if [ -z "$(PANDOC)" ]; then echo "website: pandoc not found, skipping $@" >&2; else \
		SRC="$(patsubst index.html,README.html,$@)"; SRC="$${SRC%.html}"; \
		PUBLISHED="$${SRC}.md"; if [ "$${SRC}" = "README" ]; then PUBLISHED="index.md"; fi; \
		$(PANDOC) --metadata title="$${SRC}" -V markdown-source="$${PUBLISHED}" \
			-s --to html5 "$${SRC}.md" -o "$@" \
			--lua-filter=links-to-html.lua \
			--lua-filter=add-col-scope.lua \
			--template=page.tmpl; fi

# A BonBon browser asks for the site's address and gets index.md.
index.md: README.md
	cp README.md index.md

# GitHub Pages hides folders that start with a dot (and so /.well-known/)
# unless the site says it is not a Jekyll site.
.nojekyll:
	touch .nojekyll

$(TOC): $(MD_PAGES) toc.bash
	mkdir -p $(dir $(TOC))
	./toc.bash $(BASE_URL) $(PROJECT) $(PAGE_ORDER) > $(TOC).tmp
	mv $(TOC).tmp $(TOC)

clean:
	@rm -f $(HTML_PAGES)

.FORCE:
