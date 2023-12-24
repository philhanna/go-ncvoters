package webdata

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------------
// Type Definitions
// ---------------------------------------------------------------------

type State uint8

// ---------------------------------------------------------------------
// Constants for the parsing machine
// ---------------------------------------------------------------------

const (
	INIT State = iota
	LOOKING_FOR_COLUMNS_START
	READING_COLUMNS
	LOOKING_FOR_CODE_BLOCK
	LOOKING_FOR_CODE_BLOCK_NAME
	LOOKING_FOR_CODE_BLOCK_START
	READING_CODE_BLOCK
)

// ---------------------------------------------------------------------
// Functions
// ---------------------------------------------------------------------

// ParseLayoutFile creates a layout object from the contents of a file
func ParseLayoutFile(filename string) (*Layout, error) {

	// Create an empty layout object
	layout := NewLayout()

	// Open the layout_ncvoter.txt file
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	// Create a map of code types to lines of the codes
	codeBlocks := make(map[string][]string)

	var cbName string
	var cb []string

	// Read it line by line and parse it with a finite state machine
	state := INIT

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		switch state {

		case INIT:
			// Looking for start of columns list
			if strings.HasPrefix(line, "-- File layout") {
				state = LOOKING_FOR_COLUMNS_START
			}

		case LOOKING_FOR_COLUMNS_START:
			// Looking for a row of hyphens
			if strings.HasPrefix(line, "-------") {
				state = READING_COLUMNS
			}

		case READING_COLUMNS:
			// Looking for the trailing row of hyphens
			if strings.HasPrefix(line, "-------") {
				state = LOOKING_FOR_CODE_BLOCK
			} else {
				column := NewColumn(line)
				layout.AllColumns = append(layout.AllColumns, column)
			}

		case LOOKING_FOR_CODE_BLOCK:
			if strings.HasPrefix(line, "/* ****") {
				state = LOOKING_FOR_CODE_BLOCK_NAME
			}

		case LOOKING_FOR_CODE_BLOCK_NAME:
			cbName = line
			cb = make([]string, 0)
			state = LOOKING_FOR_CODE_BLOCK_START

		case LOOKING_FOR_CODE_BLOCK_START:
			if strings.HasPrefix(line, "*******") {
				state = READING_CODE_BLOCK
			}

		case READING_CODE_BLOCK:
			// Check for end of code block
			if strings.HasPrefix(line, "*******") {
				codeBlocks[cbName] = cb
				state = LOOKING_FOR_CODE_BLOCK
			} else {
				cb = append(cb, line)
			}
		}
	}

	// Now put the code blocks into the proper slots in the Layout object
	for cbName, cb = range codeBlocks {
		switch cbName {
		case "Status codes":
			re := regexp.MustCompile(`\s+`)
			for _, line := range cb {
				tokens := re.Split(line, 2)
				layout.StatusCodes[tokens[0]] = tokens[1]
			}
		case "Race codes":
			re := regexp.MustCompile(`\s+`)
			for _, line := range cb {
				tokens := re.Split(line, 2)
				layout.RaceCodes[tokens[0]] = tokens[1]
			}
		case "Ethnic codes":
			re := regexp.MustCompile(`\s+`)
			for _, line := range cb {
				tokens := re.Split(line, 2)
				layout.EthnicCodes[tokens[0]] = tokens[1]
			}
		case "County identification number codes":
			re := regexp.MustCompile(`(\S+)\s+(\d+)`)
			for _, line := range cb {
				m := re.FindStringSubmatch(line)
				if m != nil {
					name := m[1]
					id, _ := strconv.Atoi(m[2])
					layout.CountyCodes[id] = name
				}
			}

		}
	}

	// Done
	return layout, nil
}
