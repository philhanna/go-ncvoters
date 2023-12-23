package webdata

import (
	"log"
	"net/http"
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
	sb.WriteString(CreateColumnsDDL(layout.GetColumns()))
	sb.WriteString(CreateStatusCodesDDL(layout.GetStatusCodes()))
	sb.WriteString(CreateRaceCodesDDL(layout.GetRaceCodes()))
	sb.WriteString(CreateEthnicCodesDDL(layout.GetEthnicCodes()))
	sb.WriteString(CreateCountyCodesDDL(layout))
	ddl := sb.String()

	return ddl
}
