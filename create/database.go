package create

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	goncvoters "github.com/philhanna/go-ncvoters"
	"github.com/philhanna/go-ncvoters/util"
)

// CreateDatabase is the mainline for creating a database from the zip file.
func CreateDatabase(zipFileName, entryName, dbFileName string, progressEvery int) error {
	log.Println("Creating database...")

	stime := time.Now()

	// Open the database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	// Create a table with selected columns
	query := CreateDDL()
	_, err = tx.Exec(query)
	if err != nil {
		log.Println(err)
		return err
	}

	// Create a prepared statement for inserting records into the voters
	// table.
	stmt, err := CreatePreparedStatement(tx)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	// Get the zip file entry for the embedded CSV file
	zipEntry, err := GetZipEntry(zipFileName, entryName)
	if err != nil {
		log.Println(err)
		return err
	}

	// Open the CSV file entry
	f, err := zipEntry.Open()
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()

	// Create a CSV csvReader over the zip file entry
	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	csvReader.FieldsPerRecord = -1 // Allow varying number of fields

	// Get the column names
	colNames, err := csvReader.Read()
	if err != nil {
		log.Println(err)
		return err
	}
	selectedNames := goncvoters.Configuration.GetColumnNames()
	selectedIndices := GetSelectedIndices(colNames, selectedNames)

	// Read from the CSV reader and insert records into the database
	progress := util.NewProgress()
	progress.Total = 8_300_367 // TODO get exact number from reading CSV file
	progress.SoFar = 0
	progress.LastPercent = 0

	for {

		// Read the next record
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return err
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

		// Insert a record into the database
		_, err = stmt.Exec(values...)
		if err != nil {
			log.Println(err)
			return err
		}

		progress.SoFar++
		percent := int(float64(progress.SoFar) / float64(progress.Total) * 100)
		if percent > progress.LastPercent {
			s := strings.Repeat("*", percent/2)
			for len(s) < 50 {
				s += "."
			}
			if percent > progress.LastPercent {
				fmt.Printf("Percent complete: %d%%, [%-s] %d records added in %v\r",
					percent, s, progress.SoFar, time.Since(stime))
			}
			progress.LastPercent = percent

		}
	}
	fmt.Print()

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return err
	}

	// Now copy to the real database on disk
	log.Println("Attaching physical database...")
	sql := fmt.Sprintf(`ATTACH DATABASE %q AS diskdb;`, dbFileName)
	_, err = db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Copying voters table...")
	sql = `CREATE TABLE diskdb.voters AS SELECT * FROM voters;`
	_, err = db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Detaching physical database...")
	sql = `DETACH DATABASE diskdb;`
	_, err = db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Database created successfully in %v\n", time.Since(stime))
	return nil
}
