package create

import (
	"archive/zip"
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

// CreateDatabase is the mainline for creating a database from the zip file.
func CreateDatabase(zipFileName, entryName, dbFileName string, progressEvery int) {

	// Internal function for consistent error handling
	handleError := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Creating database...")

	stime := time.Now()

	// Open the database
	db, err := sql.Open("sqlite3", ":memory:")
	handleError(err)
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	handleError(err)

	// Create a table with selected columns
	query := CreateDDL()
	_, err = tx.Exec(query)
	handleError(err)

	// Create a prepared statement for inserting records into the voters
	// table.
	stmt, err := CreatePreparedStatement(tx)
	handleError(err)
	defer stmt.Close()

	// Get the zip file entry for the embedded CSV file
	zipEntry, err := GetZipEntry(zipFileName, entryName)
	handleError(err)

	// Open the CSV file entry
	f, err := zipEntry.Open()
	handleError(err)
	defer f.Close()

	// Create a CSV csvReader over the zip file entry
	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	csvReader.FieldsPerRecord = -1 // Allow varying number of fields

	// Get the column names
	colNames, err := csvReader.Read()
	handleError(err)
	selectedNames := goncvoters.Configuration.GetColumnNames()
	selectedIndices := GetSelectedIndices(colNames, selectedNames)

	// Read from the CSV reader and insert records into the database
	progress := util.NewProgress()
	progress.Total = estimatedNumberOfVoters(zipEntry)
	progress.SoFar = 0
	progress.LastPercent = 0

	for {

		// Read the next record
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		handleError(err)

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

		// Insert a record into the database
		_, err = stmt.Exec(values...)
		handleError(err)

		progress.SoFar++
		percent := int(float64(progress.SoFar) / float64(progress.Total) * 100)
		if percent > progress.LastPercent {
			s := strings.Repeat("*", percent/2)
			for len(s) < 50 {
				s += "."
			}
			if percent > progress.LastPercent {
				countWithCommas := commas.Format(progress.SoFar)
				fmt.Printf("Percent complete: %d%%, [%-s] %s records added in %v\r",
					percent, s, countWithCommas, time.Since(stime))
			}
			progress.LastPercent = percent

		}
	}
	fmt.Print()

	// Commit the transaction
	err = tx.Commit()
	handleError(err)

	// Now copy to the real database on disk
	if util.FileExists(dbFileName) {
		log.Printf("Deleting existing disk database %s\n", dbFileName)
		os.Remove(dbFileName)
	}

	log.Println("Attaching physical database...")
	sql := fmt.Sprintf(`ATTACH DATABASE %q AS diskdb;`, dbFileName)
	_, err = db.Exec(sql)
	handleError(err)

	log.Println("Copying voters table...")
	sql = `CREATE TABLE diskdb.voters AS SELECT * FROM voters;`
	_, err = db.Exec(sql)
	handleError(err)

	log.Println("Detaching physical database...")
	sql = `DETACH DATABASE diskdb;`
	_, err = db.Exec(sql)
	handleError(err)

	log.Printf("Database created successfully in %v\n", time.Since(stime))
}

// estimatedNumberOfVoters returns the estimated number of voters based
// on a heuristic that employs a ratio of actual number of voters
// divided by compressed file size. These constants should be updated
// from time to tome.

func estimatedNumberOfVoters(file *zip.File) int64 {
	const (
		// Values from December 22, 2023 file
		NUMER = 8465201
		DENOM = 3911973311
	)
	ratio := float64(NUMER) / float64(DENOM)
	compressedSize := float64(file.FileHeader.CompressedSize64)
	countf := compressedSize * ratio
	count := int64(math.Round(countf))
	return count
}
