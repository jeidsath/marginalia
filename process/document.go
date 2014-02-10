package process

import (
	"log"
        "errors"
	"strconv"
)

// Document consists of

// Headers
//   - Text
//   - Emphasis

// Paragraphs
//   - Text
//   - Emphasis
//   - Line breaks
//   - Footnotes, sidenotes
//   - Inline quote (with citation)

// Quotations
//   - Text
//   - Emphasis
//   - Line breaks
//   - Footnotes, sidenotes
//   - Citation

type Element interface {
	ToHtml() string
	ToText() string
}

type Collection interface {
	ToStrings() []string
	ToHtml() string
}

type Text struct {
	content string
}

func (tt *Text) ToHtml() string {
	return tt.content
}

func (tt *Text) ToText() string {
	return tt.content
}

type Emphasis struct {
	Text
	Em     bool
	Strong bool
}

func (ee *Emphasis) ToHtml() string {
	output := ""

	if ee.Strong {
		output += "<strong>"
	}
	if ee.Em {
		output += "<em>"
	}

	output += ee.Text.ToHtml()

	if ee.Em {
		output += "</em>"
	}
	if ee.Strong {
		output += "</strong>"
	}

	return output
}

func (ee *Emphasis) ToText() string {
	output := ""

	if ee.Strong {
		output += "*"
	}
	if ee.Em {
		output += "_"
	}

	output += ee.Text.ToHtml()

	if ee.Em {
		output += "_"
	}
	if ee.Strong {
		output += "*"
	}

	return output
}

type LineBreak struct {
}

func (*LineBreak) ToHtml() string {
	return "</br>"
}

func (*LineBreak) ToText() string {
	return "\n"
}

type Header struct {
	Content Element
	Level   int
}

func (hh *Header) AddElement(ee Element) {
	hh.Content = ee
}

func (hh *Header) ToStrings() []string {
	hashes := ""
	for ii := 0; ii < hh.Level; ii++ {
		hashes += "#"
	}
	output := hashes + " " + hh.Content.ToText() + " " + hashes
	return []string{output}
}

func (hh *Header) ToHtml() string {
	hlevel := "h" + strconv.Itoa(hh.Level)
	inner_html := hh.Content.ToHtml()
	return "<" + hlevel + ">" + inner_html + "</" + hlevel + ">"
}

type Block struct {
	Elements []Element
}

type Note struct {
	Elements []Element
}

func (nn *Note) ToHtml() string {
	output := ""

	for ii, ee := range nn.Elements {
		switch ee.(type) {
		default:
			log.Fatalln("Note contains non-Text or non-Emphasis type")
		case *Text, *Emphasis:
			if ii != 0 {
				output += " "
			}
			output += ee.ToHtml()
		}
	}

	output += ""
	return output
}

func (nn *Note) ToText() string {
	output := ""
	for ii, ee := range nn.Elements {
		switch ee.(type) {
		default:
			log.Fatalln("Note contains non-Text or non-Emphasis type")
		case *Text, *Emphasis:
			if ii != 0 {
				output += " "
			}
			output += ee.ToText()
		}
	}
	return output
}

func (nn *Note) AddElement(ee Element) error {
        switch ee.(type) {
        default:
                return errors.New("Bad type for Note")
        case *Text, *Emphasis:
                nn.Elements = append(nn.Elements, ee)
                return nil
        }
}

func (pp *Block) AddElement(ee Element) {
	pp.Elements = append(pp.Elements, ee)
}

func (pp *Block) ToHtml() string {
	output := ""
	spaceNeeded := false

	for _, ee := range pp.Elements {
		switch ee.(type) {
		default:
			if spaceNeeded {
				output += " "
			}
			output += ee.ToHtml()
			spaceNeeded = true
		case *Footnote:
			output += "† "
			output += ee.ToHtml()
			spaceNeeded = true
		case *Leftnote:
			output += " "
			output += ee.ToHtml()
			output += " ˙"
			spaceNeeded = false
		case *Rightnote:
			output += "˚ "
			output += ee.ToHtml()
			spaceNeeded = true
		}

	}

	return output
}

func (pp *Block) ToStrings() []string {
	output := []string{""}
	line := &output[0]

	spaceNeeded := false

	for _, ee := range pp.Elements {
		switch ee.(type) {
		default:
			log.Fatalln("Bad type in ToText")
		case *Text, *Emphasis:
			if spaceNeeded {
				*line += " "
			}
			*line += ee.ToText()
			spaceNeeded = true
		case *LineBreak:
			output = append(output, "")
			line = &output[len(output)-1]
			spaceNeeded = false
		case *Footnote:
			*line += "†"
			note := "†" + ee.ToText()
			output = append(output, note)
			output = append(output, "")
			line = &output[len(output)-1]
			spaceNeeded = false
		case *Leftnote:
			note := "˙" + ee.ToText()
			output = append(output, note)
			output = append(output, "˙")
			line = &output[len(output)-1]
			spaceNeeded = false
		case *Rightnote:
			*line += "˚"
			note := "˚" + ee.ToText()
			output = append(output, note)
			output = append(output, "")
			line = &output[len(output)-1]
			spaceNeeded = false
		}
	}
	return output
}

type Paragraph struct {
	Block
}

func (pp *Paragraph) ToHtml() string {
	return "<p>" + pp.Block.ToHtml() + "</p>"
}

// These footnotes can't handle collections
// This can be changed in v2 this at some point
type Footnote struct {
        Note
}

func (ff *Footnote) ToHtml() string {
	output := "<span class=\"footnote\">†"

        output += ff.Note.ToHtml()

	output += "</span>"
	return output
}

type Leftnote struct {
	Note
}

func (ll *Leftnote) ToHtml() string {
	output := "<span class=\"leftnote\">˙"
        output += ll.Note.ToHtml()
	output += "</span>"
	return output
}

type Rightnote struct {
	Note
}

func (rr *Rightnote) ToHtml() string {
	output := "<span class=\"rightnote\">˚"
	output += rr.Note.ToHtml()
	output += "</span>"
	return output
}

type BlockQuote struct {
	Paragraphs []Paragraph
	Citation   string
}

func (bb *BlockQuote) AddParagraph(para Paragraph) {
	bb.Paragraphs = append(bb.Paragraphs, para)
}

func (bb *BlockQuote) ToHtml() string {
	inner_html := ""
	for ii, pp := range bb.Paragraphs {
		if ii != 0 {
			inner_html += "\n"
		}
		inner_html += pp.ToHtml()
	}
	cite := ""
	if bb.Citation != "" {
		cite = " cite=\"" + bb.Citation + "\""
	}
	return "<blockquote" + cite + ">\n" + inner_html + "\n</blockquote>"
}

func (bb *BlockQuote) ToStrings() []string {
	var strings []string
	for ii, pp := range bb.Paragraphs {
		if ii != 0 {
			strings = append(strings, "")
		}
		for _, ss := range pp.ToStrings() {
			strings = append(strings, ss)
		}
	}

	output := []string{}
	for _, ss := range strings {
		//indent 4
		output = append(output, "    "+ss)
	}

	if bb.Citation != "" {
		output = append(output, "    ")
		output = append(output, "    "+bb.Citation)
	}

	return output
}
