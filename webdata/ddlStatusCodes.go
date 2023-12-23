package webdata

import (
	"fmt"
	"sort"
	"strings"
)

// Creates DDL to create and load data into the status_codes table
func CreateStatusCodesDDL(statusCodes map[string]string) string {
	parts := []string{}
	parts = append(parts, "BEGIN TRANSACTION;")
	parts = append(parts, fmt.Sprintf("DROP TABLE IF EXISTS %s;", TABLE_STATUS_CODES))
	parts = append(parts, fmt.Sprintf("CREATE TABLE %s (", TABLE_STATUS_CODES))
	parts = append(parts, "  status         TEXT,")
	parts = append(parts, "  description    TEXT")
	parts = append(parts, `);`)

	// Sort the codes in alphabetical order
	keys := []string{}
	for key := range statusCodes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Write the insert statements
	for _, code := range keys {
		value := statusCodes[code]
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
