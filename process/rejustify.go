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

func makeParagraph(input []string) (*Paragraph, error) {
        ss := strings.Join(input, " ")
        para := &Paragraph{}
        para.AddElement(&Text{ss})
        return para, nil
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
