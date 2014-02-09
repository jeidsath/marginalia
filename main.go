package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"./process"
)

func toString(scanner *bufio.Scanner) (string, error) {
	output := ""
	for scanner.Scan() {
		output += scanner.Text() + "\n"
	}
	return output, scanner.Err()
}

func main() {

	var reformat bool
	var fileName string

	flag.BoolVar(&reformat, "reformat", false, "reformat margins")
	flag.StringVar(&fileName, "file", "", "filename to convert (default: stdin)")
	flag.Parse()

	text := ""

	if fileName == "" {
		scanner := bufio.NewScanner(os.Stdin)
		var err error
		text, err = toString(scanner)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var err error
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(file)
		text, err = toString(scanner)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !reformat {
		var err error
		output, err := process.Convert(text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(output)
	} else {
	}
}
