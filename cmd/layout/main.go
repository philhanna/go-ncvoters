package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/philhanna/go-ncvoters/webdata"
)

const TABLE_ALL_COLUMNS = "all_columns"
const URL = "https://s3.amazonaws.com/dl.ncsbe.gov/data/layout_ncvoter.txt"
const USAGE = `usage: layout [OPTION]

Creates the layout tables for the NCVoter database, using the latest
layout from the Board of Elections website.

options:
  -h, --help               Displays this help text and exits
  -o, --output             Output file name. If not specified, uses stdout
`

var outputFile string

func main() {

	flag.StringVar(&outputFile, "o", "", "Output file name")
	flag.StringVar(&outputFile, "output", "", "Output file name")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, USAGE)
	}
	flag.Parse()

	// Get the up-to-date layout from the voters website
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Parse the file into the data of interest
	p := webdata.NewLayout(resp.Body)

	// Open the output fp
	fp := os.Stdout
	if outputFile != "" {
		fp, err = os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()
	}

	createAllColumns(p, fp)
	createStatusCodes(p, fp)
	createRaceCodes(p, fp)
	createEthnicCodes(p, fp)
	createCountyCodes(p, fp)
}

func createAllColumns(p *webdata.Layout, fp *os.File) {
	parts := []string{}
	parts = append(parts, "BEGIN TRANSATION;")
	parts = append(parts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", TABLE_ALL_COLUMNS))
	parts = append(parts, fmt.Sprintf("CREATE TABLE %s (", TABLE_ALL_COLUMNS))
	parts = append(parts, "  colname TEXT")
	parts = append(parts, `);`)
	for _, colName := range p.GetColumns() {
		stmt := fmt.Sprintf(`INSERT INTO %s VALUES('%s')`, TABLE_ALL_COLUMNS, colName)
		parts = append(parts, stmt)
	}
	parts = append(parts, "COMMIT;")
	s := strings.Join(parts, "\n")
	fmt.Fprintln(fp, s)
}
func createStatusCodes(p *webdata.Layout, fp *os.File) {

}
func createRaceCodes(p *webdata.Layout, fp *os.File) {

}
func createEthnicCodes(p *webdata.Layout, fp *os.File) {

}
func createCountyCodes(p *webdata.Layout, fp *os.File) {

}
