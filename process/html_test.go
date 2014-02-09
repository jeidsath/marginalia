package process

import (
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


func TestParagraph(t *testing.T) {
        input := "Hello World"
        expected := "<p>Hello World</p>\n"
	output, err := Convert(input)
        if output != expected || err != nil {
		t.Fail()
	}
}
