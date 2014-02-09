# Marginalia

Formating Ancient Greek texts for digital storage and presentation presents special problems. The population of readers for these documents is small, and solutions that depend on a critical mass of programmer talent are not applicable for Ancient Greek.

XML, used by the Perseus project, has many advantages. Its flexibility makes it the obvious choice for digital Greek texts. This flexibility becomes a major disadvantage, however, when it comes clear document standards. The traditional problem of XML is that "the documentation is the code." The problem tends to become acute in small programming communities &mdash; exactly the situation faced by Ancient Greek texts.

Marginalia presents a small community solution to the Ancient Greek digital document problem

1. All documents will be human-readable. In specific, Marginalia documents are
    - UTF-8 NFC Unicode
    - Use a Markdown-style syntax
2. Because composing new texts is not a requirement, syntax for Marginalia is more exact than for Markdown
3. Unlike Markdown, Marginalia documents must support footnotes and sidenotes &mdash; a difficulty for digital documents (see flow)
4. We take advantage of Unicode to improve the presentation of text in source documents

## Syntax

### Headings

The `# Heading #` syntax is supported.

### Paragraphs

Text surrounded by blank lines is a paragraph.

### Line breaks

Manual line breaks are indicated by two spaces at the end of a line.

### Footnotes

    This is an example† of a footnote.
    While the text breaks to show the footnote
    at the bottom of the "page," it is not 
    interpreted as a paragraph break instead

    †Example footnote

    the paragraph continues unbroken. It is
    assumed that rejustifying tools can move
    the location of the page break.

The † is the unicode character, not the HTML entity.

## Flow

A document for reading has a "flow" channel. A user of this document is expected to be able to follow that flow without distraction. There are two types of document elements:

1. Flow elements (presentation elements, paragraphing, text styles)
2. Non-flow elements (marginalia, endnotes)

Digital documents are entirely "flow." Non-flow elements are generally represented with a tagging syntax. However, it is useful to have non-flow elements associated with the document.

**Example 1**

```
Flow

This is an example paragraph that 
the reader is expected to read 
straight through. There is no 
special presentation element 
beyond punctuation. We wrap the
text by specifying newlines.
```

**Example 2**

```
Non-Flow

5  This is an example paragraph that    paragraph: Gk term
   the reader is expected to read       "to write along" 
   straight through. There is no         
   special presentation* element         
6  beyond punctuation. We wrap the      punctuation: Latin? 
   text by specifying newlines.         

   * Footnotes are the great
   invention of early humanist 
   scholarship
```

**Example 3**
```
<title>Tagged non-flow</title>

<p><milestone n="5"/>This is an example 
paragraph<sidenote text="Gk term 
\"to write along\""> that the reader
is expected to read straight 
through. There is no special 
presentation<footnote 
text="Footnotes are the great 
invention of early humanist 
scholarship"> element beyond 
punctuation<sidenote text="Latin?">. 
We wrap the <milestone n="6"/>text 
by specifying newlines.</p>
```

The flow document presents no special challenges to men or computers. The non-flow document creates complexity issues for copy/paste, document modification, and scrolling. The tagged document is readable to computers, but only together with a full specification, otherwise substantial reverse engineering is required. Syntax errors represent a special challenge. It is unreadable to human begins without specialized presentation software.

## Fixing flow

### Fixing Example 3 

What about "readable" formats that currently replace XML-style solutions, such as JSON or CSV. Both formats replace tagging by replacing descriptive elements with structure. Although the formats are easier to program with, they are less expressive than full tagging, and are not remotely human readble.

### Fixing Example 2

Markdown has become the solution to the problems of creating readable documents capable of stylized presentation. It creates a human-readable and editable document that presents cleanly. 

While Markdown does not have footnote or sidenote elements, these can be envisioned. More problematic are document modification concerns. In Markdown, human beings are not forced to do their own paragraph justification. We could break flow to mark these elements, but that would destroy readability.

Instead, we are forced to look at process. The current Markdown process is as follows:

```
txt -> |Markdown markup| -> html
```

A modification to solve document editing issues could be as follows:

```
txt -> |Markdown rejustify| -> txt 
txt -> |Markdown markup| -> html
```

The rejustification step complicates editing, but not unduly. It is easy to envision inputs to that command for special cases (column width, page height, generate a "copy/paste" format).
