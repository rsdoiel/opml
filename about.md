---
title: opml
abstract: |-
  A library for working with OPML XML files and example command line programs.
authors:
  - family_name: Doiel
    given_name: R. S.
    id: https://orcid.org/0000-0003-0900-6903


maintainer:
  - family_name: Doiel
    given_name: R. S.
    id: https://orcid.org/0000-0003-0900-6903

repository_code: https://github.com/rsdoiel/opml
version: 0.0.11
license_url: https://www.gnu.org/licenses/agpl-3.0.txt

programming_language:
  - Go >= 1.27.3

keywords:
  - xml
  - opml
  - text markup
  - list markup

date_released: 2026-11-18
---

About this software
===================

## opml 0.0.11

- fixed Walk(), which visited only the first child at each level and did not stop when the callback returned false at depth
- fixed the Outline.Language tag typo so the language attribute is read and written
- the programs moved into internal/cli with documented exit codes (2, 65, 66, 73, 74, 77) and an EXIT STATUS section in each help page
- added -h, -l and -v; new -nl, -input and -output options
- opmlcat and opml2urls no longer empty an input file named by mistake; opmlsort exits non-zero on a missing input
- fixed -help, -license and -version output in opml2json and opmlsort
- dropped the unused caltechlibrary/cli requirement

## Authors

- [R. S. Doiel](https://orcid.org/0000-0003-0900-6903)




## Maintainers

- [R. S. Doiel](https://orcid.org/0000-0003-0900-6903)


A library for working with OPML XML files and example command line programs.

- [License](https://www.gnu.org/licenses/agpl-3.0.txt)
- [Code Repository](https://github.com/rsdoiel/opml)
  - [Issue Tracker](https://github.com/rsdoiel/opml/issues)

## Programming languages

- Go >= 1.27.3




## Software Requirements

- Go >= 1.27.3


## Software Suggestions

- CMTools >= 0.0.45b
- Pandoc >= 3.9
- GNU Make >= 3


