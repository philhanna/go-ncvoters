package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/philhanna/go-ncvoters/webdata"
)

// ---------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------

const (
	URL                = "https://s3.amazonaws.com/dl.ncsbe.gov/data/layout_ncvoter.txt"
	TABLE_COLUMNS      = "columns"
	TABLE_STATUS_CODES = "status_codes"
	TABLE_RACE_CODES   = "race_codes"
	TABLE_ETHNIC_CODES = "ethnic_codes"
	TABLE_COUNTY_CODES = "county_codes"
)

// ---------------------------------------------------------------------
// Variables
// ---------------------------------------------------------------------

var outputFile string

// ---------------------------------------------------------------------
// Functions
// ---------------------------------------------------------------------

func main() {

	const USAGE = `usage: layout [OPTION]

Creates the layout tables for the NCVoter database, using the latest
layout from the Board of Elections website.

options:
  -h, --help               Displays this help text and exits
  -o, --output             Output file name. If not specified, uses stdout
`

	// Handle the command line options
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
	parts = append(parts, "BEGIN TRANSACTION;")
	parts = append(parts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", TABLE_COLUMNS))
	parts = append(parts, fmt.Sprintf("CREATE TABLE %s (", TABLE_COLUMNS))
	parts = append(parts, "  name           TEXT,")
	parts = append(parts, "  dataType       TEXT,")
	parts = append(parts, "  description    TEXT")
	parts = append(parts, `);`)
	for _, column := range p.GetColumns() {
		stmt := fmt.Sprintf(`INSERT INTO %s VALUES('%s','%s','%s');`,
			TABLE_COLUMNS,
			column.Name,
			column.DataType,
			column.Description,
		)
		parts = append(parts, stmt)
	}
	parts = append(parts, "COMMIT;")
	s := strings.Join(parts, "\n") + "\n"
	fmt.Fprintln(fp, s)
}

func createStatusCodes(p *webdata.Layout, fp *os.File) {
	parts := []string{}
	parts = append(parts, "BEGIN TRANSACTION;")
	parts = append(parts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", TABLE_STATUS_CODES))
	parts = append(parts, fmt.Sprintf("CREATE TABLE %s (", TABLE_STATUS_CODES))
	parts = append(parts, "  status         TEXT,")
	parts = append(parts, "  description    TEXT")
	parts = append(parts, `);`)

	// Sort the codes in alphabetical order
	keys := []string{}
	for key := range p.StatusCodes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Write the insert statements
	for _, code := range keys {
		value := p.StatusCodes[code]
		stmt := fmt.Sprintf(`INSERT INTO %s VALUES('%s','%s');`,
			TABLE_STATUS_CODES,
			code,
			value,
		)
		parts = append(parts, stmt)
	}

	// Write the commit
	parts = append(parts, "COMMIT;")

	// Create the whole string and print it
	s := strings.Join(parts, "\n") + "\n"
	fmt.Fprintln(fp, s)
}

func createRaceCodes(p *webdata.Layout, fp *os.File) {
	parts := []string{}
	parts = append(parts, "BEGIN TRANSACTION;")
	parts = append(parts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", TABLE_RACE_CODES))
	parts = append(parts, fmt.Sprintf("CREATE TABLE %s (", TABLE_RACE_CODES))
	parts = append(parts, "  race           TEXT,")
	parts = append(parts, "  description    TEXT")
	parts = append(parts, `);`)

	// Sort the codes in alphabetical order
	keys := []string{}
	for key := range p.RaceCodes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Write the insert statements
	for _, code := range keys {
		value := p.RaceCodes[code]
		stmt := fmt.Sprintf(`INSERT INTO %s VALUES('%s','%s');`,
			TABLE_RACE_CODES,
			code,
			value,
		)
		parts = append(parts, stmt)
	}

	// Write the commit
	parts = append(parts, "COMMIT;")

	// Create the whole string and print it
	s := strings.Join(parts, "\n") + "\n"
	fmt.Fprintln(fp, s)
}

func createEthnicCodes(p *webdata.Layout, fp *os.File) {
	parts := []string{}
	parts = append(parts, "BEGIN TRANSACTION;")
	parts = append(parts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", TABLE_ETHNIC_CODES))
	parts = append(parts, fmt.Sprintf("CREATE TABLE %s (", TABLE_ETHNIC_CODES))
	parts = append(parts, "  ethnicity      TEXT,")
	parts = append(parts, "  description    TEXT")
	parts = append(parts, `);`)

	// Sort the codes in alphabetical order
	keys := []string{}
	for key := range p.EthnicCodes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Write the insert statements
	for _, code := range keys {
		value := p.EthnicCodes[code]
		stmt := fmt.Sprintf(`INSERT INTO %s VALUES('%s','%s');`,
			TABLE_ETHNIC_CODES,
			code,
			value,
		)
		parts = append(parts, stmt)
	}

	// Write the commit
	parts = append(parts, "COMMIT;")

	// Create the whole string and print it
	s := strings.Join(parts, "\n") + "\n"
	fmt.Fprintln(fp, s)
}

func createCountyCodes(p *webdata.Layout, fp *os.File) {
	parts := []string{}
	parts = append(parts, "BEGIN TRANSACTION;")
	parts = append(parts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", TABLE_COUNTY_CODES))
	parts = append(parts, fmt.Sprintf("CREATE TABLE %s (", TABLE_COUNTY_CODES))
	parts = append(parts, "  county_id      TEXT,")
	parts = append(parts, "  county         TEXT")
	parts = append(parts, `);`)

	// Sort the codes in alphabetical order
	keys := []int{}
	for key := range p.CountyCodes {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	// Write the insert statements
	for _, code := range keys {
		value := p.CountyCodes[code]
		stmt := fmt.Sprintf(`INSERT INTO %s VALUES(%d,'%s');`,
			TABLE_COUNTY_CODES,
			code,
			value,
		)
		parts = append(parts, stmt)
	}

	// Write the commit
	parts = append(parts, "COMMIT;")

	// Create the whole string and print it
	s := strings.Join(parts, "\n") + "\n"
	fmt.Fprintln(fp, s)
}
