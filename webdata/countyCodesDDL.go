package webdata

import (
	"fmt"
	"sort"
	"strings"
)

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
