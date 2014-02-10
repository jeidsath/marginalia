package process

import (
	"fmt"
	"testing"
)

func TestOutput(t *testing.T) {
	input := ""
	output, err := Convert(input)
	if output != "" || err != nil {
		t.Fail()
	}
}

//# Header #
//<h1>Header</h1>
func TestHeader(t *testing.T) {
	input := "# Header #"
	expected := "<h1>Header</h1>\n"
	output, err := Convert(input)
	if output != expected || err != nil {
		t.Fail()
	}
}

//### Header ###
//<h3>Header</h3>
func TestHeader3(t *testing.T) {
	input := "### Header ###"
	expected := "<h3>Header</h3>\n"
	output, err := Convert(input)
	if output != expected || err != nil {
		t.Fail()
	}
}

/*func TestParagraph(t *testing.T) {
        input := "Hello World"
        expected := "<p>Hello World</p>\n"
	output, err := Convert(input)
        if output != expected || err != nil {
		t.Fail()
	}
}*/

func TestEmphasis(t *testing.T) {
	ee := Emphasis{Text{"Hello"}, true, false}
	expected_txt := "_Hello_"
	expected_html := "<em>Hello</em>"

	if ee.ToText() != expected_txt || ee.ToHtml() != expected_html {
		t.Fail()
	}
}

func TestHeaderStruct(t *testing.T) {
	coll := Header{}
	coll.Level = 3
	coll.AddElement(&Text{"Hello World!"})

	expected_txt := "### Hello World! ###"
	expected_html := "<h3>Hello World!</h3>"

	if coll.ToHtml() != expected_html || coll.ToStrings()[0] != expected_txt {
		t.Fail()
	}
}

func compareStrings(lhs, rhs []string) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for ii := range lhs {
		if lhs[ii] != rhs[ii] {
			return false
		}
	}
	return true
}

func TestParagraphStruct(t *testing.T) {
	para := Paragraph{}
	para.AddElement(&Text{"This is a test paragraph"})
	para.AddElement(&Emphasis{Text{"that contains bold"}, false, true})
	para.AddElement(&Text{"text."})
	para.AddElement(&LineBreak{})
	para.AddElement(&Text{"And text after a line break."})

	expected_html := "<p>This is a test paragraph <strong>that contains bold</strong> text."
	expected_html += " </br> And text after a line break.</p>"

	expected_text := []string{"This is a test paragraph *that contains bold* text.",
		"And text after a line break."}

	if para.ToHtml() != expected_html || !compareStrings(expected_text, para.ToText()) {
		fmt.Println(para.ToHtml())
		fmt.Println(para.ToText())
		t.Fail()
	}
}

func TestFootNotes(t *testing.T) {
	para := Paragraph{}
	para.AddElement(&Text{"This is some"})
	foot := &Footnote{Text{"Only four words."}}
	para.AddElement(foot)
	para.AddElement(&Text{"text."})

	expected_html := "<p>This is some† "
	expected_html += "<span class=\"footnote\">†Only four words.</span> text.</p>"

	expected_text := []string{
		"This is some†",
		"†Only four words.",
		"text.",
	}

	if para.ToHtml() != expected_html {
		fmt.Println(para.ToHtml())
		t.Fail()
	}

	if !compareStrings(expected_text, para.ToText()) {
		for ii, ll := range para.ToText() {
			fmt.Printf("%d: %v\n", ii, ll)
		}
		t.Fail()
	}
}
