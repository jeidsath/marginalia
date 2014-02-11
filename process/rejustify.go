package process

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type intermediates interface {
}

func consumeHeaders(input *([]intermediates)) ([]intermediates, error) {
	var output []intermediates
	re := regexp.MustCompile("^(#+)(.*?)(#+)$")
	for _, ll := range *input {
		switch ll.(type) {
		default:
			output = append(output, ll)
		case string:
			if res := re.FindStringSubmatch(ll.(string)); res != nil {
				if res[1] != res[3] {
					fmt.Println(res[1])
					fmt.Println(res[3])
					return *input, errors.New("Header levels not matched")
				}
				head := Header{}
				head.Level = len(res[1])
				head.AddElement(&Text{strings.Trim(res[2], " ")})
				output = append(output, &head)
			} else {
				output = append(output, ll)
			}
		}
	}
	return output, nil
}

func consumeQuotes(input *([]intermediates)) ([]intermediates, error) {
	var output []intermediates
	var err error
	var quote *BlockQuote

	quoteLines := []string{}

	for _, ll := range *input {
		switch ll.(type) {
		default:
			output = append(output, ll)
			if len(quoteLines) != 0 {
				quote, err = makeQuote(quoteLines)
				output = append(output, quote)
			}
			quoteLines = []string{}
		case string:
			re := regexp.MustCompile("^ {4}(.*)$")
			if res := re.FindStringSubmatch(ll.(string)); res != nil {
				quoteLines = append(quoteLines, ll.(string)[4:])
			} else {
				output = append(output, ll)
				if len(quoteLines) != 0 {
					quote, err = makeQuote(quoteLines)
					output = append(output, quote)
					quoteLines = []string{}
				}
			}
		}
	}

	return output, err
}

// Return paragraph elements
func makeText(input []intermediates) ([]Element, error) {
	// Emphasis is preserved across notes, linebreaks;
	// it is an error when it hits a quote
	// Linebreaks are ignored unless they follow two spaces

	output := []Element{}
	var err error

	text := new(string)
	strong := false
	em := false

        fn := func(coll *([]Element), text *string, strong bool, em bool) (*string) {
                if len(*text) != 0 {
                        if strong || em {
                                *coll = append(*coll, &Emphasis{Text{*text}, em, strong})
                        } else {
                                *coll = append(*coll, &Text{*text})
                        }
                }
                return new(string)
        }

	for _, ll := range input {
                if *text != "" {
                        if (*text)[len(*text) - 1] != ' ' {
                                *text += " "
                        }
                }
		switch ll.(type) {
		default:
			err = errors.New("makeText: Unexpected type")
			return output, err
		case *Footnote, *Leftnote, *Rightnote:
		case *InlineQuote:
		case string:
                        newline := false
			ss := ll.(string)
                        re := regexp.MustCompile("^.* {2}$")
                        if re.MatchString(ss) {
                                ss = ss[:len(ss) - 2]
                                newline = true
                        }
			for _, letter := range ss {
				switch letter {
				case '*', '_':
                                        if len(ss) != 0 {
                                                text = fn(&output, text, strong, em)
                                        }
                                        if letter == '*' {
                                                strong = !strong
                                        } else {
                                                em = !em
                                        }
				default:
					*text += string(letter)
				}
			}
                        if newline {
                                text = fn(&output, text, strong, em)
                                output = append(output, &LineBreak{})
                        }
		}
	}

        text = fn(&output, text, strong, em)
	return output, err
}

func makeParagraph(input []string) (*Paragraph, error) {
	// 1. Find footnotes
	//    Any line starting with a dagger is turned into a footnote
	//    and placed at the last dagger location
	// 2. Find sidenotes
	//    Same, but for rings and dots
	// 3. Inline quote
	//    Beginning quote to ending quote
	//    Emphasis running across a quote line is an error

        inters := []intermediates{}
        for _, ss := range input {
                inters = append(inters, ss)
        }

	para := &Paragraph{}
        elements, err := makeText(inters)
        fmt.Printf("Length elements: %d\n", len(elements))
        for _, ee := range elements {
                para.AddElement(ee)
        }

	return para, err
}

func makeQuote(input []string) (*BlockQuote, error) {
	return &BlockQuote{}, nil
}

func consumeParagraphs(input *([]intermediates)) ([]intermediates, error) {
	var output []intermediates
	var err error
	var para *Paragraph

	paraLines := []string{}

	for _, ll := range *input {
		switch ll.(type) {
		default:
			output = append(output, ll)
			if len(paraLines) != 0 {
				para, err = makeParagraph(paraLines)
				output = append(output, para)
				paraLines = []string{}
			}
		case string:
			re := regexp.MustCompile("^$")
			if res := re.FindStringSubmatch(ll.(string)); res != nil {
				if len(paraLines) != 0 {
					para, err = makeParagraph(paraLines)
					output = append(output, para)
					paraLines = []string{}
				}
			} else {
				paraLines = append(paraLines, ll.(string))
			}
		}
	}

	if len(paraLines) != 0 {
		para, err = makeParagraph(paraLines)
		output = append(output, para)
	}

	return output, err
}

func convertToCollection(input *([]intermediates)) ([]Collection, error) {
	coll := []Collection{}
	var err error
	for _, ll := range *input {
		switch ll.(type) {
		default:
			coll = append(coll, ll.(Collection))
		case string:
			err = errors.New("Non-consumed lines")
		}

	}
	return coll, err
}

func Import(input string) ([]Collection, error) {

	var intrColl []intermediates

	lines := strings.Split(input, "\n")

	for _, ll := range lines {
		intrColl = append(intrColl, ll)
	}

	intrColl, err := consumeHeaders(&intrColl)
	if err != nil {
		return []Collection{}, err
	}

	intrColl, err = consumeQuotes(&intrColl)
	if err != nil {
		return []Collection{}, err
	}

	intrColl, err = consumeParagraphs(&intrColl)
	if err != nil {
		return []Collection{}, err
	}

	output, err := convertToCollection(&intrColl)

	fmt.Println(intrColl)
	fmt.Println(output)

	return output, err
}

func Rejustify(input []string) (string, error) {
	return strings.Join(input, "\n"), nil
}
