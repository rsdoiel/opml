
# OPML to Markdown round trip

## Overview

I wrote a Go language package to sort [OPML](http://dev.opml.org/spec2.html) outlines. 
I wrote this because my preferred [feed reader ](http://goread.io) supports manual 
sorting but not automatic alpha sorting by the _outline_ element's _text_ attribute. 

## Observations

Out of the box the OPML 2 Spec provides attributes indicating inclusion of other OPML files,
scripts, basic metadata (create, modified, authorship), and even directory structures.

[Fargo](http://fargo.io) allows userdefined attributes to be applied to the _outline_ 
element in OPML. This could be used in support some of the 
[Scrivener](https://www.literatureandlatte.com/scrivener.php)
features I miss such as describing how to render a project to various formats such as
rtf, pdf, ePub, web pages or even [Final Draft fdx](https://www.finaldraft.com/) files.

I write allot of Markdown formatted text.  Markdown is simple to index, 
search and convert into useful formats. Markdown is not good at expressing more
complex structures such as metadata. Website generators that use markdown often
require a preable or _front matter_ in the markdown to provide any metadata. This
leaves your document head cluttered and less human readable.

Another approach is to include a parallel document with the metadata.  It occurred to me 
that an OPML file could easily hold that metadata. It can even hold Markdown content.
The trouble with OPML is that it is not quick to edit by hand.

    Is there a round trip semantic mapping between OPML and Markdown?


## Germination of an idea

Entering a web link in Fargo the link is URL encoded and saved in the _text_ attribute of the 
_outline_ element.

The source view of a web links in Fargo's _outline_ element looks like

```OPML
    <outline text="&gt; href=&quot;http://example.org&quot;&lt;My example.org&gt;/a&lt;" />
```

That _outline_ element might render in Markdown as

```
    + [My element.org](http://example.org)
```

The steps to create the Markdown view are simple

1. URL decode the _text_ attribute
2. Convert HTML to Markdown

Making a round trip could be done by

3. Convert Markdown into HTML
4. For each _li_ element covert to an _outline_ element URL encoding the inner HTML of the _li_

So far so good. What about something more complex?


Here's an _outline_ element example from http://hosting.opml.org/dave/spec/directory.opml 

```OPML
    <outline text="Scripting News sites" created="Sun, 16 Oct 2005 05:56:10 GMT" type="link" url="http://hosting.opml.org/dave/mySites.opml"/>
```

To me that should look like 

```
    + [Scripting News Sites](http://hosting.opml.org/dave/mySites.opml)
```

What about the _created_ attribute? We could render it as a sub list of anchors with data uri?

This suggest a rule struct like

+ if the _text_ attribute contains HTML markup
    + URL decode into HTML
    + Convert HTML to Markdown
+ else render additional attributes as sub-lists with data URI

This might work as follows. 

```OPML
    <outline text="Scripting News sites" 
        created="Sun, 16 Oct 2005 05:56:10 GMT" 
        type="link" 
        url="http://hosting.opml.org/dave/mySites.opml"/>
```

Would become 

```Markdown
    + [Scripting News Sites](http://hosting.opml.org/dave/mySites.opml) 
        + [type](data:text/plain;link)
        + [created](data:text/date;Sun, 16 Oct 2005 05:56:10 GMT)
```

In HTML this would look like

```HTML
    <li><a href="http://histing.opml.org/dave/mySites.opml">Scripting News Sites</a>
        <ul>
            <li><a href="data:text/plain;link">type></a></li>
            <li><a href="data:text/date;Sun, 16 Oct 2005 05:56:10 GMT">created</a></li>
        </ul></li>
        
```

### Markdown to OPML

Coming back to OPML from Markdown then becomes

+ Convert Markdown to HTML
+ For each _li_ element inspect for _type_ 
    + if _li_ contains a _anchor_ and _ul_ then convert to _outline_ element and attributes
    + else 
        + URL encode
        + save in _text_ attribute of _outline_ element

Is this viable? Does it have any advantages?

