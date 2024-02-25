package main

import (
	"html/template"
	"log"
	"os"
	"path"
	"strings"
)

var (
	quoteFile      string = "quotes.dsv"
	templateFolder string = "templates"
	indexTemplate  string = path.Join(templateFolder, "index.tmpl")
	listTemplate   string = path.Join(templateFolder, "list.tmpl")
	outputFolder   string = "public"
)

type Quote struct {
	Date       string
	Quote      string
	Author     string
	BirthDeath string
	Profession string
}

// readCSV is a custom csv reader function because the default can't ignore " if they are the first character
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
		// join author and birth/death for pretty print
		quotes = append(quotes, Quote{
			Date:       fields[0],
			Quote:      fields[1],
			Author:     fields[2],
			BirthDeath: fields[3],
			Profession: fields[4],
		})
	}
	return quotes, nil
}

func main() {
	log.Println("Reading quotes dsv file")
	records, err := readCSV()
	if err != nil {
		panic(err)
	}

	log.Println("Creating index.html")
	// index template
	indexTmpl := template.Must(template.ParseFiles(indexTemplate))
	// index output file
	indexHtml, err := os.Create(path.Join(outputFolder, "index.html"))
	if err != nil {
		panic(err)
	}
	// execute/write index template
	err = indexTmpl.Execute(indexHtml, records[len(records)-1])
	if err != nil {
		panic(err)
	}

	log.Println("Creating list.html")
	listTmpl := template.Must(template.ParseFiles(listTemplate))
	// list output file
	listHtml, err := os.Create(path.Join(outputFolder, "list.html"))
	if err != nil {
		panic(err)
	}
	// execute template
	err = listTmpl.Execute(listHtml, records)
	if err != nil {
		panic(err)
	}
}
