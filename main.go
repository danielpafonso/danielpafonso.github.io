package main

import (
	"fmt"
	"os"
	"strings"
)

var (
	quoteFile string = "quotes.dsv"
	//outputFolder string = "public"
)

type Quote struct {
	date       string
	quote      string
	author     string
	birthDeath string
	profession string
}

// readCSV is a custom csv reader function because the default can't ignore " if they are the first caracther
func readCSV() ([]Quote, error) {
	quotes := make([]Quote, 0)
	// open file
	data, err := os.ReadFile(quoteFile)
	if err != nil {
		return nil, err
	}
	records := strings.Split(string(data), "\n")
	//process qoutes, skiping header row
	for _, record := range records[1:] {
		fields := strings.Split(record, "|")
		quotes = append(quotes, Quote{
			date:       fields[0],
			quote:      fields[1],
			author:     fields[2],
			birthDeath: fields[3],
			profession: fields[4],
		})
	}
	return quotes, nil
}

func main() {
	records, err := readCSV()
	if err != nil {
		panic(err)
	}
	for _, quote := range records {
		fmt.Println(quote)
	}
}
