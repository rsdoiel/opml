
# Ideas about an OPML

## background

OPML is a spec that comes out of the tradition of blogging, sharing links and relates to lists of RSS feeds.
It is also the document format of a very nice outline editor called [Fargo](http://fargo.io) written by
[Dave Winer](http://scripting.com). I've used other outliners in the past but I think Fargo is the simplest.
This has left me two personal use cases - OPML for feed reading and outlines in Fargo. 

## Why a Go package for OPML

[The feed read I use](http://goread.io) is driven from an OPML file. Over time I tend to collect allot of feeds in my list.
The feed reader supports manual sorting of the feeds but doesn't support automatic alphabetical sorting. The feed
reader does support both import and export of the feed list in OPML format. 

I few weeks ago I got tired of sorting the feed list by hand.  I exported my feed list in OPML format. Wrote
a small Go language package (I like write code in Go) that does three things

+ Reads in an OPML file to a in memory data structure
+ Writes out that data structure in OPML format
+ Sorts the feed list (and sub lists) alphabetically

With that I created a command line tool called _opmlsort_.  Since then I've been think about what else I
should do with this little library.

## So what's next

I enjoy writing.  I enjoy outlining my projects, presentations and even structure for my personal creative
writing projects.  When I was in grad school I used to use a wonderful program called Scrivener. The trouble
is since grad school I stopped using a Mac and Scrivener is no longer available to me. Scrivener has many
clever features but the three I really miss having are

+ A Notebook (aka outline) structure for organizing material
+ The abililty to render the "Notebook" to a single document or set of web pages
    + It support ePub as well as specialize formats like Final Draft fdx files
+ Each sub list in the structure could stand as either a part of a document itself

Dave Winer wrote a wonderful outliner called Fargo. When I am in a web browser I can use it. It's nice. The
file format for Fargo is OPML. This is not suprising if you follow Dave. 

Most of the computers I use at home are [tiny Raspberry Pi computers](http://raspberrypi.org) and running 
a web browser isn't always an option or convient. I like using them for writing. They are powerful enough 
to run your usual Unix text processing tools and simple enough not to be distracting (e.g. no web 
distractions).

When writing prose in the Unix shell I tend to format by text in Markdown. Often my Markdown starts out as
an outline. This got me to thinking about OPML and how that might help me create a simple scrivener like
environment for writing that worked in the Unix shell.

A step in this direction would be to write a simple conversion tool for opml2markdown and markdown2opml.
This would give me the flexibility of using Fargo when avialable or Markdown when not. But how should that
mapping work?





Basic text processing via the veneral Unix shell and vi is very reasonable to do on a Raspberry Pi. I can
also write programs in Go and compile them for the Raspberry Pi.
For writing since my return to the Unix shell I spend allot of time in vi. I can easily git nicely formatted
text for the web with Markdown. Markdown also supports creating simple list. This got me thinking. Could
I use an OPML file to create 





OIt is open, I can mine it for interesting information,
and it is both human and but the feed reader I current prefer [goread.io](http://goread.io) has some limitations.
One of the limitations is there is no automatic sorting of the feed list. You can manually sort the list but when 
the list is long it's a challenge to sort alphabetically.  I like the command line, I currently write a fair
amount of code in Go, so I wrote a small Go language package to read and write a OPML file as well as sort it
alphabetically by the text attribute in the outline element.  This made it trivial to implement a command line
tool called opmlsort that does what it says.

## So what's next

Now that I have this simple little package I've been itching to do something else with it. I spend most
of my computing time at the Unix shell prompt. I am very comfortable there.  I have also been thinking about
how I turned away from wordpressors and GUI tools over the years for writing purposes.  I enjoy the simplicity
of the what I can do in text processing in the shell.  I have also noticed how writing prose in Markdown format
has creaped in even to my creative writing time. This got me thinking. What's missing from my command line
world besides the richness of graphics?

My one of my favorite creative writing tools is Scrivener.  I used it write my thesis, to write a
presentations, to edit my unpublished screenplays and novellas. The trouble is I don't get to use it
anymore because I don't use a Mac at home.  The computers I personally own these days are tiny and inexpensive. 
They all run Linux for the most part and Scrivener isn't available. What I liked about Scrivener is three featuers

+ You can organize material in an outline structure called a Notebook
+ Each node can be thought of as a single file or as a part of a larger document
+ You can render that outline structure to a single document (including Final Drafts' fdx format)







at the console.
I like to keep my
feeds sorted alphabetically.  
I wrote a small OPML
package in Go

I like outlines but these days I tend to write allot in Markdown format.
