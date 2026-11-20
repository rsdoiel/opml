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
version: 0.0.12
license_url: https://www.gnu.org/licenses/agpl-3.0.txt

programming_language:
  - Go >= 1.27.3

keywords:
  - xml
  - opml
  - text markup
  - list markup
  - markdown
  - commonmark

date_released: 2026-11-20
---

About this software
===================

## opml 0.0.12

- new opml2cmark and cmark2opml write an OPML outline as a CommonMark list and read it back, to edit in a text editor; new cmark package (opml DR-0001)
- Parse, ReadFile and Unmarshal read ISO-8859-1, US-ASCII and windows-1252 documents; Dave Winer's toc.opml could not be read before
- new Flatten() and Includes() for reading a table of contents; new MarshalIndent and WriteFileIndent; new Head.Generator
- opml2urls lists the pages of a table of contents (-url, on by itself for a document with no xmlUrl or htmlUrl; -includes) and exits 0, not 70, when there is nothing to list
- WriteFile returns its errors and no longer panics on a nil Body
- CHANGED: Head.OwnerID is read and written as ownerId (a spec-correct `<ownerId>` was dropped before), and opml2json writes ownerId and dateModified (were OwnerId and the typo dataModified)
- removed the TypeScript port (opml.ts)
- now requires github.com/yuin/goldmark and gopkg.in/yaml.v3, imported by the cmark package and the two new programs only

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
- github.com/yuin/goldmark >= 1.8.7
- gopkg.in/yaml.v3 >= 3.0.1


## Software Suggestions

- CMTools >= 0.0.45b
- Pandoc >= 3.9
- GNU Make >= 3


