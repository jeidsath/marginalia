package process

import (
        "errors"
	"regexp"
	"strings"
        "fmt"
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
                                output = append(output, head)
			} else {
                                output = append(output, ll)
                        }
		}
	}
	return output, nil
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

        fmt.Println(intrColl)

	return []Collection{&Header{}, &Paragraph{}}, nil
}

func Rejustify(input []string) (string, error) {
	return strings.Join(input, "\n"), nil
}
