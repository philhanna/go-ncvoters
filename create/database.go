package create

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/philhanna/commas"
	goncvoters "github.com/philhanna/go-ncvoters"
	"github.com/philhanna/go-ncvoters/util"
)

var (
	selectedIndices []int
	colNames        []string
	progress        *util.Progress
	stime           time.Time
	MAX_ENTRIES     int
)

// CreateDatabase is the mainline for creating a database from the zip file.
func CreateDatabase(zipFileName, entryName, dbFileName string, progressEvery int) error {
	if !quiet() {
		log.Println("Creating database...")
	}
	stime = time.Now()

	// Open the database
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()

	// Begin a transaction
	tx, _ := db.Begin()

	// Create a table with selected columns
	query := CreateDDL()
	tx.Exec(query)

	// Create a prepared statement for inserting records into the voters
	// table.
	stmt, _ := CreatePreparedStatement(tx)
	defer stmt.Close()

	// Get the zip file entry for the embedded CSV file
	zipEntry, err := GetZipEntry(zipFileName, entryName)
	if err != nil {
		return err
	}

	// Initialize the progress bar
	progress = util.NewProgress()
	progress.Total = estimatedNumberOfVoters(zipEntry.UncompressedSize64)
	progress.SoFar = 0
	progress.LastPercent = 0

	// Open the CSV file entry
	f, _ := zipEntry.Open()
	defer f.Close()

	// Create a CSV csvReader over the zip file entry
	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	csvReader.FieldsPerRecord = -1 // Allow varying number of fields

	// Get the column names
	colNames, _ = csvReader.Read()
	selectedNames := goncvoters.Configuration.GetColumnNames()
	selectedIndices = GetSelectedIndices(colNames, selectedNames)

	// Read from the CSV reader and insert records into the database
	entries := 0
	for values := range readFromCSV(csvReader) {
		entries++
		if MAX_ENTRIES > 0 && entries > MAX_ENTRIES {
			break
		}

		// Insert a record into the database
		stmt.Exec(values...)

		// Update the progress bar
		showProgress()
	}
	fmt.Println()

	// Commit the transaction
	tx.Commit()

	// Now copy to the real database on disk
	if util.FileExists(dbFileName) {
		if !quiet() {
			log.Printf("Deleting existing disk database %s\n", dbFileName)
		}
		os.Remove(dbFileName)
	}

	if !quiet() {
		log.Println("Attaching physical database...")
	}
	sql := fmt.Sprintf(`ATTACH DATABASE %q AS diskdb;`, dbFileName)
	db.Exec(sql)

	if !quiet() {
		log.Println("Copying voters table...")
	}
	sql = `CREATE TABLE diskdb.voters AS SELECT * FROM voters;`
	db.Exec(sql)

	if !quiet() {
		log.Println("Detaching physical database...")
	}
	sql = `DETACH DATABASE diskdb;`
	db.Exec(sql)

	if !quiet() {
		log.Printf("Database created successfully in %v\n", time.Since(stime))
	}

	// Return without error
	return nil
}

// estimatedNumberOfVoters returns the estimated number of voters based
// on a heuristic that employs a ratio of actual number of voters
// divided by compressed file size. This is only used for the progress bar.
// These constants should be updated from time to time.
func estimatedNumberOfVoters(size uint64) int64 {
	const (
		// Values from December 22, 2023 file
		NUMER = 8465201
		DENOM = 3911973311
	)
	ratio := float64(NUMER) / float64(DENOM)
	countf := float64(size) * ratio
	count := int64(math.Round(countf))
	return count
}

func readFromCSV(reader *csv.Reader) chan []any {
	ch := make(chan []any, 100)
	go func() {
		defer close(ch)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}

			// Choose just the columns we want
			values := make([]any, len(selectedIndices))
			for i, idx := range selectedIndices {
				colName := colNames[idx]
				if IsSanitizeCol(colName) {
					value := string(record[idx])
					values[i] = Sanitize(value)
				} else {
					values[i] = record[idx]
				}
			}
			ch <- values // Send the row of values down the channel
		}
	}()
	return ch
}

func showProgress() {
	progress.SoFar++
	percent := int(float64(progress.SoFar) / float64(progress.Total) * 100)
	if percent > progress.LastPercent {
		s := strings.Repeat("*", percent/2)
		for len(s) < 50 {
			s += "."
		}
		if percent > progress.LastPercent {
			countWithCommas := commas.Format(progress.SoFar)
			if !quiet() {
				fmt.Printf("Percent complete: %d%%, [%-s] %s records added in %v\r",
					percent, s, countWithCommas, time.Since(stime))
			}
		}
		progress.LastPercent = percent

	}
}
