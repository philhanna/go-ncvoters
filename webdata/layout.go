package webdata

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/philhanna/fsm/v2"
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
// Functions
// ---------------------------------------------------------------------

const (
	INIT fsm.State = iota
	LOOKING_FOR_COLUMNS
	READING_COLUMNS
	DONE_WITH_COLUMNS
	LOOKING_FOR_STATUS
	READING_STATUS
	DONE_WITH_STATUS
	LOOKING_FOR_RACE_CODES
	READING_RACE_CODES
	DONE_WITH_RACE_CODES
	LOOKING_FOR_ETHNIC_CODES
	READING_ETHNIC_CODES
	DONE_WITH_ETHNIC_CODES
	LOOKING_FOR_COUNTY_CODES
	READING_COUNTY_CODES
	DONE_WITH_COUNTY_CODES
)

// NewLayout parses the file layouts from an io.Reader
func NewLayout(reader io.Reader) (*Layout, error) {
	layout := new(Layout)
	layout.AllColumns = make([]Column, 0)
	layout.StatusCodes = make(map[string]string)
	layout.RaceCodes = make(map[string]string)
	layout.EthnicCodes = make(map[string]string)
	layout.CountyCodes = make(map[int]string)

	machine := fsm.FSM[string]{
		States: []fsm.State{
			INIT,
			LOOKING_FOR_COLUMNS,
			READING_COLUMNS,
			DONE_WITH_COLUMNS,
			LOOKING_FOR_STATUS,
			READING_STATUS,
			DONE_WITH_STATUS,
			LOOKING_FOR_RACE_CODES,
			READING_RACE_CODES,
			DONE_WITH_RACE_CODES,
			LOOKING_FOR_ETHNIC_CODES,
			READING_ETHNIC_CODES,
			DONE_WITH_ETHNIC_CODES,
			LOOKING_FOR_COUNTY_CODES,
			READING_COUNTY_CODES,
			DONE_WITH_COUNTY_CODES,
		},
		InitialState: INIT,
		TransitionMap: map[fsm.State]fsm.Transition[string]{
			INIT:                     layout.LookingForColumnStart,
			LOOKING_FOR_COLUMNS:      layout.LookingForColumns,
			READING_COLUMNS:          layout.ReadingColumns,
			DONE_WITH_COLUMNS:        layout.DoneWithColumns,
			LOOKING_FOR_STATUS:       layout.LookingForStatus,
			READING_STATUS:           layout.ReadingStatus,
			DONE_WITH_STATUS:         layout.DoneWithStatus,
			LOOKING_FOR_RACE_CODES:   layout.LookingForRaceCodes,
			READING_RACE_CODES:       layout.ReadingRaceCodes,
			DONE_WITH_RACE_CODES:     layout.DoneWithRaceCodes,
			LOOKING_FOR_ETHNIC_CODES: layout.LookingForEthnicCodes,
			READING_ETHNIC_CODES:     layout.ReadingEthnicCodes,
			DONE_WITH_ETHNIC_CODES:   layout.DoneWithEthnicCodes,
			LOOKING_FOR_COUNTY_CODES: layout.LookingForCountyCodes,
			READING_COUNTY_CODES:     layout.ReadingCountyCodes,
			DONE_WITH_COUNTY_CODES:   layout.DoneWithCountyCodes,
		},
		Trace: false,
	}

	// Define the input and output channels
	inch := make(chan fsm.Event[string])
	defer close(inch)
	ouch := machine.Run(inch)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		inch <- fsm.Event[string](line)
		<-ouch
	}
	return layout, nil
}

// ---------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------

// GetColumns returns a list of all columns in the web site's database
func (pl *Layout) GetColumns() []Column {
	return pl.AllColumns
}

// GetStatusCodes returns a map of status codes to their descriptions
func (pl *Layout) GetStatusCodes() map[string]string {
	return pl.StatusCodes
}

// GetRaceCodes returns a map of race codes to their descriptions
func (pl *Layout) GetRaceCodes() map[string]string {
	return pl.RaceCodes
}

// GetEthnicCodes returns a map of ethnic codes to their descriptions
func (pl *Layout) GetEthnicCodes() map[string]string {
	return pl.EthnicCodes
}

// GetCountyCodes returns a map of county numbers to county names
func (pl *Layout) GetCountyCodes() map[int]string {
	return pl.CountyCodes
}

// ---------------------------------------------------------------------
// Transition functions for the finite state machine
// ---------------------------------------------------------------------

func (pl *Layout) LookingForColumnStart(event fsm.Event[string]) fsm.State {
	re := regexp.MustCompile(`name\s+data type\s+description`)
	line := event.(string)
	if re.MatchString(line) {
		return LOOKING_FOR_COLUMNS
	}
	return INIT
}

func (pl *Layout) LookingForColumns(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "-----") {
		return READING_COLUMNS
	}
	return LOOKING_FOR_COLUMNS
}

func (pl *Layout) ReadingColumns(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "-----") {
		return DONE_WITH_COLUMNS
	}
	column := NewColumn(line)
	pl.AllColumns = append(pl.AllColumns, column)
	return READING_COLUMNS
}

func (pl *Layout) DoneWithColumns(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "Status codes") {
		return LOOKING_FOR_STATUS
	}
	return DONE_WITH_COLUMNS
}

func (pl *Layout) LookingForStatus(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "*****") {
		return READING_STATUS
	}
	return LOOKING_FOR_STATUS
}

func (pl *Layout) ReadingStatus(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "*****") {
		return DONE_WITH_STATUS
	}
	re := regexp.MustCompile(`\s+`)
	tokens := re.Split(line, 2)
	pl.StatusCodes[tokens[0]] = tokens[1]
	return READING_STATUS
}

func (pl *Layout) DoneWithStatus(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "Race codes") {
		return LOOKING_FOR_RACE_CODES
	}
	return DONE_WITH_STATUS
}

func (pl *Layout) LookingForRaceCodes(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "*****") {
		return READING_RACE_CODES
	}
	return LOOKING_FOR_RACE_CODES
}

func (pl *Layout) ReadingRaceCodes(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "*****") {
		return DONE_WITH_RACE_CODES
	}
	re := regexp.MustCompile(`\s+`)
	tokens := re.Split(line, 2)
	pl.RaceCodes[tokens[0]] = tokens[1]
	return READING_RACE_CODES
}

func (pl *Layout) DoneWithRaceCodes(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "Ethnic codes") {
		return LOOKING_FOR_ETHNIC_CODES
	}
	return DONE_WITH_RACE_CODES
}

func (pl *Layout) LookingForEthnicCodes(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "*****") {
		return READING_ETHNIC_CODES
	}
	return LOOKING_FOR_ETHNIC_CODES
}

func (pl *Layout) ReadingEthnicCodes(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "*****") {
		return DONE_WITH_ETHNIC_CODES
	}
	re := regexp.MustCompile(`\s+`)
	tokens := re.Split(line, 2)
	pl.EthnicCodes[tokens[0]] = tokens[1]
	return READING_ETHNIC_CODES
}

func (pl *Layout) DoneWithEthnicCodes(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "County identification number codes") {
		return LOOKING_FOR_COUNTY_CODES
	}
	return DONE_WITH_ETHNIC_CODES
}

func (pl *Layout) LookingForCountyCodes(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "*****") {
		return READING_COUNTY_CODES
	}
	return LOOKING_FOR_COUNTY_CODES
}

func (pl *Layout) ReadingCountyCodes(event fsm.Event[string]) fsm.State {
	line := event.(string)
	if strings.HasPrefix(line, "*****") {
		return DONE_WITH_COUNTY_CODES
	}
	re := regexp.MustCompile(`\s+`)
	tokens := re.Split(line, 2)
	countyName := tokens[0]
	countyNumber, _ := strconv.Atoi(tokens[1])
	pl.CountyCodes[countyNumber] = countyName
	return READING_COUNTY_CODES
}

func (pl *Layout) DoneWithCountyCodes(event fsm.Event[string]) fsm.State {
	return DONE_WITH_COUNTY_CODES
}
