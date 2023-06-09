package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/philhanna/go-ncvoters/webdata"
	_ "github.com/mattn/go-sqlite3"
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

var dbName string

// ---------------------------------------------------------------------
// Functions
// ---------------------------------------------------------------------

func main() {

	const USAGE = `usage: metadata [OPTION]

Creates the metadata tables for the North Carolina voter registration
database, using the latest layout from the Board of Elections website.

options:
  -h, --help               Displays this help text and exits.
  -o, --output             Output database name. If not specified,
                           uses "metadata.db" in the temp directory.
`

	// Handle the command line options
	flag.StringVar(&dbName, "o", "", "Output file name")
	flag.StringVar(&dbName, "output", "", "Output file name")
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
	if dbName == "" {
		dbName = filepath.Join(os.TempDir(), "metadata.db")
	}

	// Start building the DDL string
	sb := strings.Builder{}
	sb.WriteString(createColumns(p))
	sb.WriteString(createStatusCodes(p))
	sb.WriteString(createRaceCodes(p))
	sb.WriteString(createEthnicCodes(p))
	sb.WriteString(createCountyCodes(p))
	ddl := sb.String()

	// Open a connection to the database file
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Execute the DDL statement
	_, err = db.Exec(ddl)
	if err != nil {
		log.Println("Error creating table:", err)
		return
	}

	log.Printf("Metadata databse created in %s\n", dbName)
}

// Creates DDL to create and load data into the columns table
func createColumns(p *webdata.Layout) string {
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
	return strings.Join(parts, "\n") + "\n"
}

func createStatusCodes(p *webdata.Layout) string {
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

	// Create the whole string
	return strings.Join(parts, "\n") + "\n"
}

func createRaceCodes(p *webdata.Layout) string {
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

	// Create the whole string
	return strings.Join(parts, "\n") + "\n"
}

func createEthnicCodes(p *webdata.Layout) string {
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

	// Create the whole string
	return strings.Join(parts, "\n") + "\n"
}

func createCountyCodes(p *webdata.Layout) string {
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
	return strings.Join(parts, "\n") + "\n"
}
