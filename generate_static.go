package main

import (
	"html/template"
	"log"
	"os"
	"path"
	"slices"
	"strings"
	"time"
)

var (
	quoteFile         string = "quotes.dsv"
	templateFolder    string = "templates"
	indexTemplate     string = path.Join(templateFolder, "index.html.tmpl")
	listTemplate      string = path.Join(templateFolder, "list.html.tmpl")
	duplicateTemplate string = path.Join(templateFolder, "duplicates.html.tmpl")
	outputFolder      string = "public"
)

type Quote struct {
	Date       string
	Quote      string
	Author     string
	BirthDeath string
	Profession string
}

type DuplicateQuotes struct {
	Dates  []string
	Quote  string
	Author string
}

// levenshteinRatio calculate the ratio of levenshtein distance
func levenshteinRatio(strA, strB string) float32 {
	n := len(strA)
	m := len(strB)
	lev := make([][]int, n+1)
	for i := range len(lev) {
		lev[i] = make([]int, m+1)
	}

	// populate matrix
	for i := range n + 1 {
		lev[i][0] = i
	}
	for i := range m + 1 {
		lev[0][i] = i
	}

	// calculate distance
	for i := 1; i < (n + 1); i++ {
		for j := 1; j < (m + 1); j++ {
			insertion := lev[i-1][j] + 1
			deletion := lev[i][j-1] + 1
			delta := 0
			if strA[i-1] != strB[j-1] {
				delta = 1
			}
			substitution := lev[i-1][j-1] + delta
			lev[i][j] = min(min(insertion, deletion), substitution)
		}
	}
	// calculate ration
	return float32(n+m-lev[n][m]) / float32(n+m)
}

// calculateSimilar check and group simular quote
func calculateSimilar(quotes []Quote, limit float32) []DuplicateQuotes {
	duplicates := make([]DuplicateQuotes, 0)

	for len(quotes) > 1 {
		dups := DuplicateQuotes{
			Dates:  []string{quotes[0].Date},
			Quote:  quotes[0].Quote,
			Author: quotes[0].Author,
		}
		for i := 1; i < len(quotes); i++ {
			ratio := levenshteinRatio(dups.Quote, quotes[i].Quote)
			if ratio > limit {
				dups.Dates = append(dups.Dates, quotes[i].Date)
				quotes = slices.Delete(quotes, i, i+1)
				i--
			}
		}
		if len(dups.Dates) > 1 {
			duplicates = append(duplicates, dups)
		}
		quotes = slices.Delete(quotes, 0, 1)
	}

	return duplicates
}

// readCSV is a custom csv reader function because the default can't ignore " if they are the first character
func readCSV() ([]Quote, error) {
	today := time.Now().UTC()
	quotes := make([]Quote, 0)
	// open file
	data, err := os.ReadFile(quoteFile)
	if err != nil {
		return nil, err
	}
	records := strings.Split(
		strings.TrimSpace(string(data)),
		"\n",
	)
	//process qoutes, skiping header row
	for _, record := range records[1:] {
		fields := strings.Split(record, "|")
		// check if is future quote
		quoteDay, err := time.Parse("2006-01-02", fields[0])
		if err != nil {
			return nil, err
		}
		if quoteDay.After(today) {
			continue
		}
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

	// Calculate Similar Quotes
	log.Println("Calculate Similar")
	similars := calculateSimilar(records, 0.8)

	log.Println("Creating duplicates.html")
	duplicateTmpl := template.Must(template.ParseFiles(duplicateTemplate))
	duplicateHtml, err := os.Create(path.Join(outputFolder, "duplicates.html"))
	if err != nil {
		panic(err)
	}
	// execute template
	err = duplicateTmpl.Execute(duplicateHtml, similars)
	if err != nil {
		panic(err)
	}
}
