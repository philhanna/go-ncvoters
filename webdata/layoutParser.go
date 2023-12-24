package webdata

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Parse creates a layout object from the contents of a file
func Parse(filename string) (*Layout, error) {
	layout := new(Layout)
	return layout, nil
}

// ParseColumns get the list of columns from a file
func ParseColumns(filename string) ([]Column, error) {
	columns := make([]Column, 0)

	// Open the layout_ncvoter.txt file
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	const (
		INIT uint8 = iota
		LOOKING_FOR_BOUNDARY
		READING_COLUMNS
		DONE
	)
	state := INIT

	// Read the layout file looking for columns
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		switch state {
		case INIT:
			if strings.HasPrefix(line, "-- File layout") {
				state = LOOKING_FOR_BOUNDARY
			}
		case LOOKING_FOR_BOUNDARY:
			if strings.HasPrefix(line, "-------") {
				state = READING_COLUMNS
			}
		case READING_COLUMNS:
			if strings.HasPrefix(line, "-------") {
				state = DONE
			} else {
				column := NewColumn(line)
				columns = append(columns, column)
			}
		}
	}

	// Verify that columns section was closed
	if state != DONE {
		err := fmt.Errorf("did not find closing boundary for columns. State=%v", state)
		return columns, err
	}
	
	return columns, nil
}
