package webdata

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// ---------------------------------------------------------------------
// Type definitions
// ---------------------------------------------------------------------

type Layout struct {
	AllColumns  []Column
	StatusCodes map[string]string
	RaceCodes   map[string]string
	EthnicCodes map[string]string
	CountyCodes map[int]string
}

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
// Constructor
// ---------------------------------------------------------------------

// NewLayout parses the file layouts from an io.Reader
func NewLayout() *Layout {
	layout := new(Layout)
	layout.AllColumns = make([]Column, 0)
	layout.StatusCodes = make(map[string]string)
	layout.RaceCodes = make(map[string]string)
	layout.EthnicCodes = make(map[string]string)
	layout.CountyCodes = make(map[int]string)
	return layout
}

// ---------------------------------------------------------------------
// Functions
// ---------------------------------------------------------------------

// DownloadLayout gets the latest layout data from the voters website
// and writes it to a file in the system temporary directory
func DownloadLayout(url string) (string, error) {

	const BUFSIZ = 65536

	// Get the page with the layout data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	statusCode := resp.StatusCode
	if statusCode != 200 {
		err := fmt.Errorf("expected HTTP status code 200, got %d", statusCode)
		return "", err
	}

	// Write the page to /tmp/voter_layout.txt
	filename := path.Join(os.TempDir(), "voter_layout.txt")
	fp, err := os.Create(filename)
	if err != nil {
		return filename, err
	}
	defer fp.Close()

	buffer := make([]byte, BUFSIZ)

readLoop:
	for {
		n, err := resp.Body.Read(buffer)
		switch {
		case err == io.EOF:
			if n > 0 {
				fp.Write(buffer[:n])
			}
			break readLoop
		case err != nil:
			return filename, err
		default:
			fp.Write(buffer[:n])
		}
	}

	return filename, nil

}

// ---------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------

// GetMetadataDDL returns the metadata DDL extracted from a layout
// object
func (layout *Layout) GetMetadataDDL() (string, error) {
	sb := strings.Builder{}
	sb.WriteString(CreateColumnsDDL(layout.AllColumns))
	sb.WriteString(CreateStatusCodesDDL(layout.StatusCodes))
	sb.WriteString(CreateRaceCodesDDL(layout.RaceCodes))
	sb.WriteString(CreateEthnicCodesDDL(layout.EthnicCodes))
	sb.WriteString(CreateCountyCodesDDL(layout.CountyCodes))
	ddl := sb.String()
	return ddl, nil
}
