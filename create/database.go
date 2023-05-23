package create

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
)

// CreateDatabase is the mainline for creating a database from the zip file.
func CreateDatabase(zipFileName, csvFileName, dbFileName string, progressEvery int) error {
	log.Println("Creating database...")

	// Open the database
	db, err := sql.Open("sqlite3", dbFileName)
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
	zipEntry, err := GetZipEntry(zipFileName, csvFileName)
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
 
	// Get the column names.
	colNames, err := csvReader.Read()
	if err != nil {
		log.Println(err)
		return err
	}
	selectedIndices := GetSelectedIndices(colNames, selectedCols)

	// Read from the CSV reader and insert records into the database
	count := 0
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
		values := make([]interface{}, len(selectedIndices))
		for i, idx := range selectedIndices {
			values[i] = record[idx]
		}

		// Insert a record into the database
		_, err = stmt.Exec(values...)
		if err != nil {
			log.Println(err)
			return err
		}

		count++
		if count%progressEvery == 0 {
			fmt.Printf("%d records added\r", count)
		}
	}
	fmt.Print()

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Database created successfully!")
	return nil
}
