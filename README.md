# Marginalia

Formating Ancient Greek texts for digital storage and presentation presents special problems. The population of readers for these documents is small, and solutions that depend on a critical mass of programmer talent are not applicable for Ancient Greek.

XML, used by the Perseus project, has many advantages. Its flexibility makes it the obvious choice for digital Greek texts. This flexibility becomes a major disadvantage, however, when it comes clear document standards. The traditional problem of XML is that "the documentation is the code." The problem tends to become acute in small programming communities &mdash; exactly the situation faced by Ancient Greek texts.

Marginalia presents a small community solution to the Ancient Greek digital document problem

1. All documents will be human-readable. In specific, Marginalia documents are
    - UTF-8 NFC Unicode
    - Use a Markdown-style syntax
2. Because composing new texts is not a requirement, syntax for Marginalia is more exact than for Markdown
3. Unlike Markdown, Marginalia documents must support footnotes and sidenotes &mdash; a difficulty for digital documents (see the discussion of flow at the tail of this document)
4. We take advantage of Unicode to improve the presentation of text in source documents

## Syntax

### Headings

    # Heading Level 1 #
    ## Heading Level 2 ##

### Emphasis

    *em (italic)*

    _em (italic)_

    **strong (bold)**

    __strong (bold)__

### Paragraphs

    Text surrounded by blank lines is a paragraph.
    Line breaks are ignored.

### Line breaks

    Manual line breaks can be indicated by  
    two spaces at the end of a line.

### Footnotes

    This is an example† of a footnote.
    While the text breaks to show the footnote
    at the bottom of the "page," it is not 
    interpreted as a paragraph break instead

    †Example footnote

    the paragraph continues unbroken. It is
    assumed that rejustifying tools can move
    the location of the page break.

The † is the unicode character, not the HTML entity. An actual paragraph break is indicated by another blank line.

### Sidenotes

    ˙*1*  This is an example of ˙left and right 
          sidnotes. In a text with left 
    sidenotes, all document text is indented by 
    X + 3 characters (which can be any reasonable
                    number of characters) for
    ˙Slightly       the left sidenote. ˙Left
     more complex   sidenotes are placed within
     sidenote       the first X characters if 
                    possible, but can extend 
    beyond this, if necessary, as long as 3 
        spaces remain in the middle in order to
    ˙3  separate the main text ˙from the
        sidenote. Left sidenotes are also marked
    in the text by a ˙ character.

    Right sidenotes are indicated similarly with 
    a ring˚ instead of a dot To demark the right    ˚Ring example
    right sidenote text. The right sidenote 
    channel should also be separated from the 
    text.

### Quotations

    Because we do not need a 'code' syntax we
    instead use that for quotations. Curly quotes
    are used to begin and end the quote. An
    optional cite line at the bottom is allowed.
    
        “Four score and seven years ago our 
        fathers brought forth on this continent 
        a new nation, conceived in liberty, 
        and dedicated to the proposition that 
        all men are created equal.

        Now we are engaged in a great civil 
        war, testing whether that nation, or 
        any nation so conceived and so 
        dedicated, can long endure. We are met 
        on a great battlefield of that war. We 
        have come to dedicate a portion of that 
        field, as a final resting place for 
        those who here gave their lives that 
        that nation might live. It is 
        altogether fitting and proper that we 
        should do this.”

        Abraham Lincoln, G. A. 

    Inline quotes are also allowed like with the
    same technique. “Curly quotes are used to
    begin and end the quote. ‖ Myself” A unicode
    double bar is used to cite, if necessary. 


## Rejustification

In general, the above paragraph justification is too complicated without tooling. Therefore marginalia will rejustify txt files after editing, as long as proper spacing between channels is maintained.

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
