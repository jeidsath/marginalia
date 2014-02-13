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

	AddTextEm := func(coll *([]Element), text *string, strong bool, em bool) *string {
		ss := strings.Trim(*text, " \n")
		if len(ss) != 0 {
			if strong || em {
				*coll = append(*coll, &Emphasis{Text{ss}, em, strong})
			} else {
				*coll = append(*coll, &Text{ss})
			}
		}
		return new(string)
	}

	for _, ll := range input {
		if *text != "" {
			if (*text)[len(*text)-1] != ' ' {
				*text += " "
			}
		}
		switch ll.(type) {
		default:
			err = errors.New("makeText: Unexpected type")
			return output, err
		case *Footnote:
                        text = AddTextEm(&output, text, strong, em)
			output = append(output, ll.(*Footnote))
		case *Leftnote:
			output = append(output, ll.(*Leftnote))
		case *Rightnote:
			output = append(output, ll.(*Rightnote))
		case *InlineQuote:
			output = append(output, ll.(*InlineQuote))
		case string:
			newline := false
			ss := ll.(string)
			re := regexp.MustCompile("^.* {2}$")
			if re.MatchString(ss) {
				ss = ss[:len(ss)-2]
				newline = true
			}
			for _, letter := range ss {
				switch letter {
				case '*', '_':
					if len(ss) != 0 {
						text = AddTextEm(&output, text, strong, em)
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
				text = AddTextEm(&output, text, strong, em)
				output = append(output, &LineBreak{})
			}
		}
	}

	text = AddTextEm(&output, text, strong, em)
	return output, err
}

type consumeFootnoteState struct {
	lastBlank          bool
	footNoteInProgress bool
	footNoteCompleted  bool
	startLine          int
	endLine            int
	foot               []string
}

func (cfs *consumeFootnoteState) startFootnote(ii int) {
	cfs.startLine = ii - 1
	cfs.footNoteInProgress = true
}

func (cfs *consumeFootnoteState) endFootnote(ii int) {
	cfs.endLine = ii
	cfs.footNoteInProgress = false
	cfs.footNoteCompleted = true
}

func (cfs *consumeFootnoteState) nonString(ii int) {
	if cfs.footNoteInProgress {
		cfs.endFootnote(ii)
	} else {
		cfs.lastBlank = false
	}
}

func (cfs *consumeFootnoteState) stringEncountered(ii int, ss string) {
	if !cfs.footNoteInProgress && !cfs.footNoteCompleted {
		if ss == "" {
			cfs.lastBlank = true
		}
		if cfs.lastBlank && len(ss) >= 3 {
			if ss[:3] == dagger {
				cfs.startFootnote(ii)
				cfs.foot = append(cfs.foot, ss[3:])
			}
		}
                return
	}

	if cfs.footNoteInProgress {
                if ss == "" {
                        cfs.endFootnote(ii + 1)
                }
		cfs.foot = append(cfs.foot, ss)
	}
}

func makeFootnote(input []string) (*Footnote, error) {
	//when presented with strings,
	//makeText consumes them as Text, Emphasis
	inter := []intermediates{}
	for _, ss := range input {
		inter = append(inter, ss)
	}
	elements, err := makeText(inter)
	foot := &Footnote{}
	foot.Note.Elements = elements
	return foot, err
}

func consumeFootnotes(input *([]intermediates)) ([]intermediates, error) {

	var err error
	findFootnote := func(input *([]intermediates)) consumeFootnoteState {

		state := consumeFootnoteState{}

		for ii, ll := range *input {
			switch ll.(type) {
			default:
				state.nonString(ii)
			case string:
				state.stringEncountered(ii, ll.(string))
			}
		}

		state.nonString(len(*input))
		return state
	}

        putAtDagger := func(input *([]intermediates), foot *Footnote) ([]intermediates, error) {
                output := []intermediates{}
                var addedFootnote bool = false
                for _, ll := range *input {
                        if addedFootnote {
                                output = append(output, ll)
                                continue
                        }
                        switch ll.(type) {
                        default:
                                output = append(output, ll)
                        case string:
                                ss := ll.(string)
                                idx := strings.Index(ss, dagger)
                                if idx == -1 {
                                        output = append(output, ll)
                                } else {
                                        ss1 := strings.Trim(ss[0:idx], " ")
                                        output = append(output, ss1)
                                        output = append(output, foot)
                                        if idx + 3 < len(ss) {
                                                ss2 := strings.Trim(ss[idx + 3:], " ")
                                                output = append(output, ss2)
                                        }
                                        addedFootnote = true
                                }
                        }
                }
                if addedFootnote {
                        return output, nil
                } else {
                        return output, errors.New("No footnote added")
                }
        }

	replaceFootnote := func(input *([]intermediates)) ([]intermediates, bool, error) {
		var err error
		state := findFootnote(input)
		output := []intermediates{}
		if state.footNoteCompleted {
			for ii, ll := range *input {
				if ii < state.startLine || ii >= state.endLine {
					output = append(output, ll)
				}
			}
                        var foot *Footnote
                        foot, err = makeFootnote(state.foot)
                        if err != nil {
                                return output, state.footNoteCompleted, err
                        }
                        output, err = putAtDagger(&output, foot)
		        return output, state.footNoteCompleted, err
		} else {
                        return *input, state.footNoteCompleted, err
                }
	}

        inters := []intermediates{}
        var completed bool
        for inters, completed, err = replaceFootnote(input); completed; {
		if err != nil {
			return inters, err
		}
		inters, completed, err = replaceFootnote(&inters)
	}
        
	return inters, err
}

func printIntermediates(input []intermediates) {
        for ii, ll := range input {
                switch ll.(type) {
                default:
                        fmt.Printf("%d: %v\n", ii, ll.(Element).ToText())
                case string:
                        fmt.Printf("%d: %v\n", ii, ll.(string))
                }
        }
}

func printElements(input []Element) {
        for ii, ll := range input {
                fmt.Printf("%d: %v\n", ii, ll.ToText())
        }
}

func makeParagraph(input []intermediates) (*Paragraph, error) {
        //input will be a collection of strings, notes
	var err error

	para := &Paragraph{}
	elements, err := makeText(input)
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

	paraLines := []intermediates{}

	for _, ll := range *input {
		switch ll.(type) {
		default:
			output = append(output, ll)
			if len(paraLines) != 0 {
				para, err = makeParagraph(paraLines)
				output = append(output, para)
				paraLines = []intermediates{}
			}
                case *Footnote:
                       //TODO 
                       //Fix makeParagraphs to handle footnotes
                       //Fix consumeFootnotes to work correctly
                       paraLines = append(paraLines, ll)
		case string:
			re := regexp.MustCompile("^$")
			if res := re.FindStringSubmatch(ll.(string)); res != nil {
				if len(paraLines) != 0 {
					para, err = makeParagraph(paraLines)
					output = append(output, para)
					paraLines = []intermediates{}
				}
			} else {
				paraLines = append(paraLines, ll)
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

        intrColl, err = consumeFootnotes(&intrColl)
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

	return output, err
}

func Rejustify(input []string) (string, error) {
	return strings.Join(input, "\n"), nil
}
