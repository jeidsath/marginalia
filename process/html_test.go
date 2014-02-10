package process

import (
	"fmt"
	"testing"
//        "reflect"
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

	if para.ToHtml() != expected_html || !compareStrings(expected_text, para.ToStrings()) {
		fmt.Println(para.ToHtml())
		fmt.Println(para.ToStrings())
		t.Fail()
	}
}

func TestFootnotes(t *testing.T) {
	para := Paragraph{}
	para.AddElement(&Text{"This is some"})
	foot := &Footnote{}

	foot.AddElement(&Text{"Only"})
	foot.AddElement(&Emphasis{Text{"four"}, false, true})
	foot.AddElement(&Text{"words."})

	para.AddElement(foot)
	para.AddElement(&Text{"text."})

	expected_html := "<p>This is some† "
	expected_html += "<span class=\"footnote\">†Only <strong>four</strong> words.</span> text.</p>"

	expected_text := []string{
		"This is some†",
		"†Only *four* words.",
		"text.",
	}

	if para.ToHtml() != expected_html {
		fmt.Println(para.ToHtml())
		t.Fail()
	}

	if !compareStrings(expected_text, para.ToStrings()) {
		for ii, ll := range para.ToStrings() {
			fmt.Printf("%d: %v\n", ii, ll)
		}
		t.Fail()
	}
}

func TestSideNotes(t *testing.T) {
	para := Paragraph{}
	para.AddElement(&Text{"This is a sentence with a"})
	left := &Leftnote{}
	left.AddElement(&Text{"L"})
	para.AddElement(left)
	para.AddElement(&Text{"leftnote. And this is a rightnote"})
	right := &Rightnote{}
	right.AddElement(&Text{"Rightnote text"})
	para.AddElement(right)
	para.AddElement(&Text{"sentence."})

	expected_html := "<p>This is a sentence with a <span class=\"leftnote\">˙L</span> "
	expected_html += "˙leftnote. And this is a rightnote˚ <span class=\"rightnote\">"
	expected_html += "˚Rightnote text</span> sentence.</p>"

	expected_text := []string{
		"This is a sentence with a",
		"˙L",
		"˙leftnote. And this is a rightnote˚",
		"˚Rightnote text",
		"sentence.",
	}

	if para.ToHtml() != expected_html {
		fmt.Println(para.ToHtml())
		t.Fail()
	}

	if !compareStrings(expected_text, para.ToStrings()) {
		for ii, ll := range para.ToStrings() {
			fmt.Printf("%d: %v\n", ii, ll)
		}
		t.Fail()
	}
}

func TestBlockQuote(t *testing.T) {
	quote := BlockQuote{}

	para := Paragraph{}
	para.AddElement(&Text{"This is a sentence with a"})
	left := &Leftnote{}
	left.AddElement(&Text{"L"})
	para.AddElement(left)
	para.AddElement(&Text{"leftnote. And this is a rightnote"})
	right := &Rightnote{}
	right.AddElement(&Text{"Rightnote text"})
	para.AddElement(right)
	para.AddElement(&Text{"sentence."})

	quote.AddParagraph(para)
	quote.AddParagraph(para)

	quote.Citation = "Joel"

	expected_html := "<blockquote cite=\"Joel\">\n"
	expected_html += "<p>This is a sentence with a "
	expected_html += "<span class=\"leftnote\">˙L</span> "
	expected_html += "˙leftnote. And this is a rightnote˚ <span class=\"rightnote\">"
	expected_html += "˚Rightnote text</span> sentence.</p>\n"
	expected_html += "<p>This is a sentence with a "
	expected_html += "<span class=\"leftnote\">˙L</span> "
	expected_html += "˙leftnote. And this is a rightnote˚ <span class=\"rightnote\">"
	expected_html += "˚Rightnote text</span> sentence.</p>\n"
	expected_html += "</blockquote>"

	expected_text := []string{
		"    This is a sentence with a",
		"    ˙L",
		"    ˙leftnote. And this is a rightnote˚",
		"    ˚Rightnote text",
		"    sentence.",
		"    ",
		"    This is a sentence with a",
		"    ˙L",
		"    ˙leftnote. And this is a rightnote˚",
		"    ˚Rightnote text",
		"    sentence.",
		"    ",
		"    Joel",
	}

	if quote.ToHtml() != expected_html {
		fmt.Println(quote.ToHtml())
		t.Fail()
	}

	if !compareStrings(expected_text, quote.ToStrings()) {
		for ii, ll := range quote.ToStrings() {
			if ii < len(expected_text) && ll != expected_text[ii] {
				fmt.Printf("(SHOULD BE %d: %v)\n", len(expected_text[ii]),
					expected_text[ii])
			}
			fmt.Printf("%d: %v\n", len(ll), ll)
		}
		t.Fail()
	}

}

func TestInlineQuote(t *testing.T) {
	para := Paragraph{}

	para.AddElement(&Text{"This is a sentence with an"})
	quote := InlineQuote{}
	quote.AddElement(&Text{"inline quote as part of the"})
	quote.Citation = "Joel"
	para.AddElement(&quote)
	para.AddElement(&Text{"sentence."})

	expected_html := "<p>This is a sentence with an <q cite=\"Joel\">"
	expected_html += "inline quote as part of the</q> sentence.</p>"

	expected_text := []string{
		"This is a sentence with an “inline quote as part of the ‖ Joel” sentence.",
	}

	if para.ToHtml() != expected_html {
		fmt.Println(para.ToHtml())
		t.Fail()
	}

	if !compareStrings(expected_text, para.ToStrings()) {
		for ii, ll := range para.ToStrings() {
			if ii < len(expected_text) && ll != expected_text[ii] {
				fmt.Printf("(SHOULD BE %d: %v)\n", len(expected_text[ii]),
					expected_text[ii])
			}
			fmt.Printf("%d: %v\n", len(ll), ll)
		}
		t.Fail()
	}

}

func TestImport(t *testing.T) {
        document := "# Title #\n"
        document += "\n"
        document += "A short paragraph that doesn't say very\n"
        document += "much and is wrapped by hand to make sure\n"
        document += "that we are smart enought to pick up the\n"
        document += "entire paragraph.\n"

        coll, err := Import(document)

        if err != nil {
                fmt.Println(err)
                t.Fail()
        }

        /*if len(coll) > 1 && reflect.TypeOf(coll[0]) != reflect.TypeOf(&Header{}) {
                fmt.Printf("Collection[0]: %v\n", reflect.TypeOf(coll[0]))
                t.Fail()
        } else {
                exp := []string{"Title"}
                if head := coll[0].(*Header) ; head.Level != 1 || 
                compareStrings(coll[0].ToStrings(), exp) {
                        fmt.Println("Bad Header")
                        t.Fail()
                }
        }

        if len(coll) > 2 && reflect.TypeOf(coll[1]) != reflect.TypeOf(&Paragraph{}) {
                fmt.Printf("Collection[0]: %v\n", reflect.TypeOf(coll[1]))
                t.Fail()
        }*/

        if len(coll) != 2 {
                fmt.Printf("Collection size: %d\n", len(coll))
                t.Fail()
        }
}

func TestRejustify(t *testing.T) {

	lines := []string{
		"This is a test",
		"document.",
	}

	expected_text := "This is a test\ndocument."

	output, err := Rejustify(lines)

	if err != nil || output != expected_text {
		fmt.Println(Rejustify(lines))
		t.Fail()
	}
}
