package main

import (
	"flag"
	"fmt"
	"strings"
)

var fileName, search, fields, format string

func init() {
	flag.StringVar(&fileName, "file", "", "MARC file to process. Required.")
	flag.StringVar(&search, "m", "", "Only records that match the string passed (case insensitive).")
	flag.StringVar(&fields, "fields", "", "Comma delimited list of fields to output.")
	flag.StringVar(&format, "format", "MRK", "Output format (defaults to MRK).")
	flag.Parse()
}

func main() {

	if fileName == "" {
		fmt.Printf("marcli parameters:\r\n")
		flag.PrintDefaults()
		return
	}

	file, err := NewMarcFile(fileName)
	if err != nil {
		panic(err)
	}

	processor := ConsoleProcessor{
		Filters:     NewFieldFilters(fields),
		SearchValue: strings.ToLower(search),
		Format:      format,
	}

	err = file.ReadAll(&processor)
	if err != nil {
		panic(err)
	}
}
