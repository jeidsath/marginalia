package process

import (
	"log"
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
	AddElement(Element)
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

type Paragraph struct {
	Elements []Element
}

func (pp *Paragraph) AddElement(ee Element) {
	pp.Elements = append(pp.Elements, ee)
}

func (pp *Paragraph) ToHtml() string {
	output := "<p>"

	for ii, ee := range pp.Elements {
		if ii != 0 {
			output += " "
		}
		output += ee.ToHtml()
	}

	output += "</p>"
	return output
}

func (pp *Paragraph) ToText() []string {
	output := []string{""}
	line := &output[0]

	for _, ee := range pp.Elements {
		switch ee.(type) {
		default:
			log.Fatalln("Bad type in ToText")
		case *Text, *Emphasis:
			if *line != "" {
				*line += " "
			}
			*line += ee.ToText()
		case *LineBreak:
			output = append(output, "")
			line = &output[len(output)-1]
		}
	}
	return output
}
