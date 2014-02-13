package process

import (
	"fmt"
	"strings"
	"testing"
	//"reflect"
)

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
	expected_html += "</br>\nAnd text after a line break.</p>"

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
                "",
		"†Only *four* words.",
                "",
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

func collectionString(coll []Collection) string {
	output := ""
	for _, cc := range coll {
		output += strings.Join(cc.ToStrings(), "\n") + "\n\n"
	}
	return output
}

func collectionHtml(coll []Collection) string {
	output := ""
	for _, cc := range coll {
		output += cc.ToHtml() + "\n"
	}
	return output
}

func TestImport(t *testing.T) {
	document := "# Title #\n"
	document += "\n"
	document += "A short *paragraph* that doesn't say very\n"
	document += "much and is wrapped by _hand_ to make sure\n"
	document += "that we are smart enought to pick up the\n"
	document += "entire paragraph.\n"
	document += "\n"
	document += "## Header2 ##\n"
	document += "\n"
	document += "A short paragraph that doesn't say anything.  \n"
	document += "But Roses are Red  \n"
	document += "And Violets aren't  \n"
	document += "Poetry is great!\n"
	document += "Isn't it?\n"

	coll, err := Import(document)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if len(coll) != 4 {
		fmt.Printf("Collection size: %d\n", len(coll))
		t.Fail()
	}

	expected := "# Title #\n"
	expected += "\n"
	expected += "A short *paragraph* that doesn't say very "
	expected += "much and is wrapped by _hand_ to make sure "
	expected += "that we are smart enought to pick up the "
	expected += "entire paragraph.\n"
	expected += "\n"
	expected += "## Header2 ##\n"
	expected += "\n"
	expected += "A short paragraph that doesn't say anything.\n"
	expected += "But Roses are Red\n"
	expected += "And Violets aren't\n"
	expected += "Poetry is great! Isn't it?\n\n"

	ss := collectionString(coll)
	if expected != ss {
		fmt.Println(ss)
		t.Fail()
	}

	expected_html := "<h1>Title</h1>\n"
	expected_html += "<p>A short <strong>paragraph</strong> "
	expected_html += "that doesn't say very much and is wrapped "
	expected_html += "by <em>hand</em> to make sure that we are "
	expected_html += "smart enought to pick up the entire "
	expected_html += "paragraph.</p>\n"
	expected_html += "<h2>Header2</h2>\n"
	expected_html += "<p>A short paragraph that doesn't say anything.</br>\n"
	expected_html += "But Roses are Red</br>\n"
	expected_html += "And Violets aren't</br>\n"
	expected_html += "Poetry is great! Isn't it?</p>\n"

	ss = collectionHtml(coll)
	if expected_html != ss {
		fmt.Println(ss)
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

func TestUnicode(t *testing.T) {
	ss := "† And some text"
	if ss[:3] != "†" {
		fmt.Println(len("†"))
		t.Fail()
	}
}

func printComparedStrings(ss1, ss2 string) {
        fmt.Printf("||%v||\n", ss1)
        fmt.Printf("||%v||\n", ss2)

        longer := ss1
        shorter := ss2
        if len(ss2) > len(ss1) {
                longer = ss2
                shorter = ss1
        }

        output := []rune{}
        for ii, cc := range longer {
                if ii >= len(shorter) {
                        output = append(output, 'X')
                        continue
                }
                if longer[ii] != shorter[ii] {
                        output = append(output, 'X')
                } else {
                        output = append(output, cc)
                }
        }
        fmt.Print("||")
        fmt.Print(string(output))
        fmt.Print("||\n")
}

func TestImportFootnotes(t *testing.T) {

	document := "# Title #\n"
	document += "\n"
	document += "This is a test paragraph line one.\n"
        document += "This is a test† paragraph line two (with footnote after test).\n"
	document += "This is a test paragraph line three.\n"
        document += "\n"
        document += "†Hello world\n"
        document += "\n"
	document += "This is a test paragraph continued after the footnote.\n"

        expected_html := "<h1>Title</h1>\n"
        expected_html += "<p>This is a test paragraph line one. This is a test† "
        expected_html += "<span class=\"footnote\">†Hello world</span> paragraph line two "
        expected_html += "(with footnote after test). This is a test paragraph line three. "
        expected_html += "This is a test paragraph continued after the footnote.</p>\n"

        expected_text := "# Title #\n"
        expected_text += "\n"
        expected_text += "This is a test paragraph line one. This is a test†\n"
        expected_text += "\n"
        expected_text += "†Hello world\n"
        expected_text += "\n"
        expected_text += "paragraph line two (with footnote after test). "
        expected_text += "This is a test paragraph line three. This is a test "
        expected_text += "paragraph continued after the footnote.\n"
        expected_text += "\n"

	coll, err := Import(document)
	if err != nil {
                fmt.Printf("Error: %v", err)
		t.Fail()
	}

        if collectionHtml(coll) != expected_html {
                printComparedStrings(collectionHtml(coll), expected_html)
                t.Fail()
        }

        if collectionString(coll) != expected_text {
                printComparedStrings(collectionString(coll), expected_text)
                t.Fail()
        }
}
