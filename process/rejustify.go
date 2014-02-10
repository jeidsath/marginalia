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
	re := regexp.MustCompile("^(#+)(.*)(#+)$")
	for _, ll := range *input {
		switch ll.(type) {
		default:
			output = append(output, ll)
		case string:
			if res := re.FindStringSubmatch(ll.(string)); res != nil {
				if res[1] != res[3] {
					return *input, errors.New("Header levels not matched")
				}
				head := Header{}
				head.Level = len(res[1])
				head.AddElement(&Text{res[2]})
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

	for _, ll := range *input {
		switch ll.(type) {
		default:
			output = append(output, ll)
		case string:
			re := regexp.MustCompile("^ {4}(.*)$")
			if res := re.FindStringSubmatch(ll.(string)); res != nil {
				//TODO
			} else {
				output = append(output, ll)
			}
		}
	}

	return output, err
}

func consumeParagraphs(input *([]intermediates)) ([]intermediates, error) {
	var output []intermediates
	var err error

	para := &Paragraph{}

	for _, ll := range *input {
		switch ll.(type) {
		default:
			output = append(output, ll)
			if !para.Empty() {
				output = append(output, para)
			}
			para = &Paragraph{}
		case string:
			re := regexp.MustCompile("^$")
			if res := re.FindStringSubmatch(ll.(string)); res != nil {
				if !para.Empty() {
					output = append(output, para)
				}
				para = &Paragraph{}
			} else {
				if para.Empty() {
					para.AddElement(&Text{ll.(string)})
				} else {
					le := para.Elements[len(para.Elements)-1]
                                        le.(*Text).content += " " + ll.(string)
				}
			}
		}
	}

	if !para.Empty() {
		output = append(output, para)
	}

	return output, err
}

func consumeBlankLines(input *([]intermediates)) ([]Collection, error) {
	coll := []Collection{}
	var err error
	for _, ll := range *input {
		switch ll.(type) {
		default:
			coll = append(coll, ll.(Collection))
		case string:
			if ll.(string) != "" {
				err = errors.New("Non-blank lines")
			}
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

	output, err := consumeBlankLines(&intrColl)

	fmt.Println(intrColl)
	fmt.Println(output)

	return output, err
}

func Rejustify(input []string) (string, error) {
	return strings.Join(input, "\n"), nil
}
