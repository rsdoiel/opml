
# USAGE

	opmlcat [OPTIONS] OPML_FILE [OPML_FILE ...]

## SYNOPSIS


opmlcat concatenates one or more opml files as siblings to standard out.


## OPTIONS

```
    -examples                 display examples
    -generate-markdown-docs   generate Markdown documentation
    -h, -help                 display help
    -i, -input                set input filename
    -l, -license              display license
    -nl, -newline             add a trailing newline
    -o, -output               set output filename
    -p, -pretty               pretty print XML output
    -quiet                    suppress error messages
    -v, -version              display version
```


## EXAMPLES


This is an example of using opmlcat and opmlsort together to 
create a combined sorted opml file.

    opmlcat file1.opml file1.opml | opmlsort -o combined-sorted.opml


opmlcat v0.0.5
