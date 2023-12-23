package webdata

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
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
// Functions
// ---------------------------------------------------------------------

func GetMetadataDDL() string {

	// Get the up-to-date layout from the voters website
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Parse the file into the data of interest
	layout := NewLayout(resp.Body)

	// Start building the DDL string
	sb := strings.Builder{}
	sb.WriteString(CreateColumnsDDL(layout))
	sb.WriteString(CreateStatusCodesDDL(layout))
	sb.WriteString(CreateRaceCodesDDL(layout))
	sb.WriteString(CreateEthnicCodesDDL(layout))
	sb.WriteString(CreateCountyCodesDDL(layout))
	ddl := sb.String()

	return ddl
}

// Creates DDL to create and load data into the columns table
func CreateColumnsDDL(layout *Layout) string {
	parts := []string{}
	parts = append(parts, "BEGIN TRANSACTION;")
	parts = append(parts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", TABLE_COLUMNS))
	parts = append(parts, fmt.Sprintf("CREATE TABLE %s (", TABLE_COLUMNS))
	parts = append(parts, "  name           TEXT,")
	parts = append(parts, "  dataType       TEXT,")
	parts = append(parts, "  description    TEXT")
	parts = append(parts, `);`)
	for _, column := range layout.GetColumns() {
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

// Creates DDL to create and load data into the status_codes table
func CreateStatusCodesDDL(layout *Layout) string {
	parts := []string{}
	parts = append(parts, "BEGIN TRANSACTION;")
	parts = append(parts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", TABLE_STATUS_CODES))
	parts = append(parts, fmt.Sprintf("CREATE TABLE %s (", TABLE_STATUS_CODES))
	parts = append(parts, "  status         TEXT,")
	parts = append(parts, "  description    TEXT")
	parts = append(parts, `);`)

	// Sort the codes in alphabetical order
	keys := []string{}
	for key := range layout.StatusCodes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Write the insert statements
	for _, code := range keys {
		value := layout.StatusCodes[code]
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

// Creates DDL to create and load data into the race_codes table
func CreateRaceCodesDDL(layout *Layout) string {
	parts := []string{}
	parts = append(parts, "BEGIN TRANSACTION;")
	parts = append(parts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", TABLE_RACE_CODES))
	parts = append(parts, fmt.Sprintf("CREATE TABLE %s (", TABLE_RACE_CODES))
	parts = append(parts, "  race           TEXT,")
	parts = append(parts, "  description    TEXT")
	parts = append(parts, `);`)

	// Sort the codes in alphabetical order
	keys := []string{}
	for key := range layout.RaceCodes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Write the insert statements
	for _, code := range keys {
		value := layout.RaceCodes[code]
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

// Creates DDL to create and load data into the ethnic_codes table
func CreateEthnicCodesDDL(layout *Layout) string {
	parts := []string{}
	parts = append(parts, "BEGIN TRANSACTION;")
	parts = append(parts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", TABLE_ETHNIC_CODES))
	parts = append(parts, fmt.Sprintf("CREATE TABLE %s (", TABLE_ETHNIC_CODES))
	parts = append(parts, "  ethnicity      TEXT,")
	parts = append(parts, "  description    TEXT")
	parts = append(parts, `);`)

	// Sort the codes in alphabetical order
	keys := []string{}
	for key := range layout.EthnicCodes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Write the insert statements
	for _, code := range keys {
		value := layout.EthnicCodes[code]
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

// Creates DDL to create and load data into the county table
func CreateCountyCodesDDL(layout *Layout) string {
	parts := []string{}
	parts = append(parts, "BEGIN TRANSACTION;")
	parts = append(parts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", TABLE_COUNTY_CODES))
	parts = append(parts, fmt.Sprintf("CREATE TABLE %s (", TABLE_COUNTY_CODES))
	parts = append(parts, "  county_id      TEXT,")
	parts = append(parts, "  county         TEXT")
	parts = append(parts, `);`)

	// Sort the codes in alphabetical order
	keys := []int{}
	for key := range layout.CountyCodes {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	// Write the insert statements
	for _, code := range keys {
		value := layout.CountyCodes[code]
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
