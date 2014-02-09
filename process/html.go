package process

import (
        "bufio"
        "strings"
        "strconv"
        "regexp"
)

// Convert text -> html
func Convert(input string) (string, error) {
        ss := input

        ss = convertHeaders(ss)
        ss = paragraphs(ss)

        return ss, nil
}

// # Header # -> <h1>Header</h1>
// ##Header -> <h2>Header</h2>, etc.
func convertHeaders(input string) string {
        reader := strings.NewReader(input)
        scanner := bufio.NewScanner(reader)

        output := ""

        for scanner.Scan() {
                line := scanner.Text()

                for _, ii := range []int{6,5,4,3,2,1} {
                        hs := "#{" + strconv.Itoa(ii) + "}"
                        pattern := "^" + hs + " ?(.*?) ?#*$"
                        rr := regexp.MustCompile(pattern)

                        header := "h" + strconv.Itoa(ii)
                        repl := "<" + header + ">$1</" + header + ">"
                        line = string(rr.ReplaceAll([]byte(line), []byte(repl)))
                }

                output += line + "\n"
        }

        return output

}

func paragraphs(input string) string {
        return input
}
