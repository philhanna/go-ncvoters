package ncvoters

import (
	"archive/zip"
	"database/sql"

	// "encoding/csv"
	"fmt"
	// "io"
	"log"
	// "os"
	"strings"
)

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
	// table
	stmt, err := CreatePreparedStatement(tx)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	// Open the zip file
	zipFile, err := zip.OpenReader(zipFileName)
	if err != nil {
		log.Println(err)
		return err
	}
	defer zipFile.Close()

	// Find the CSV file in the zip archive
	var csvFile *zip.File
	for _, file := range zipFile.File {
		if file.Name == csvFileName {
			csvFile = file
			break
		}
	}

	// If the CSV file is not found, exit with an error
	if csvFile == nil {
		err = fmt.Errorf("file %q not found in zip archive", csvFileName)
		log.Println(err)
		return err
	}

	// Open the CSV file inside the archive
	csvReader, err := csvFile.Open()
	if err != nil {
		log.Println(err)
		return err
	}
	defer csvReader.Close()
	/*
		// Read the CSV file into a byte buffer





		csvFile, err := os.Open(csvFileName)
		if err != nil {
			log.Println(err)
			return err
		}
		defer csvFile.Close()

		reader := csv.NewReader(csvFile)
		reader.FieldsPerRecord = -1 // Allow varying number of fields

		columns, err := reader.Read()
		if err != nil {
			log.Println(err)
			return err
		}

		selectedIndices := GetSelectedIndices(columns, selectedCols)

		count := 0
		for {
			record, err := reader.Read()
			if err == io.EOF {
				log.Println(err)
				break
			} else if err != nil {
				log.Println(err)
				return err
			}

			values := make([]interface{}, len(selectedIndices))
			for i, idx := range selectedIndices {
				values[i] = record[idx]
			}

			_, err = stmt.Exec(values...)
			if err != nil {
				log.Println(err)
				return err
			}

			count++
			if count%progressEvery == 0 {
				log.Printf("%d records inserted", count)
			}
		}

		err = tx.Commit()
		if err != nil {
				log.Println(err)
			return err
		}
	*/
	log.Println("Database created successfully!")
	return nil
}

// CreateInsertSQL creates an SQL string that can be used for inserting
// records into the voters table.
func CreateInsertSQL(cols []string) string {

	// Create a string with a comma-separated list of column names.
	colNames := strings.Join(cols, ",")

	// Create a string with a comma-separated list of question marks for
	// the "VALUES" part of the SQL.
	qMarks := strings.Repeat("?,", len(cols)-1) + "?"

	// Create the SQL text of the INSERT statement using the two
	// substrings created above.
	sqlString := fmt.Sprintf("INSERT INTO voters (%s) VALUES (%s)", colNames, qMarks)

	return sqlString
}

// CreatePreparedStatement creates an SQL statement for inserting
// records into the voters table.
func CreatePreparedStatement(tx *sql.Tx) (*sql.Stmt, error) {
	sqlString := CreateInsertSQL(selectedCols)
	stmt, err := tx.Prepare(sqlString)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

// GetSelectedIndices returns the indices of selected columns from a
// given list of columns.
//
// It takes two parameters: 'columns' - a slice of strings representing
// all available columns, and 'selectedCols' - a slice of strings
// representing the columns that are selected.
//
// The function iterates over the 'columns' slice, checks if each column
// exists in the 'selectedCols' slice, and if found, appends the index
// of the column to the 'selectedIndices' slice.  Finally, it returns
// the 'selectedIndices' slice containing the indices of selected
// columns.
func GetSelectedIndices(columns, selectedCols []string) []int {

	// Create an empty slice to store the selected column indices.
	selectedIndices := make([]int, 0)

	// Iterate over each column in the 'columns' slice.
	for i, col := range columns {

		// Initialize a flag variable to track if the column is found in 'selectedCols'.
		found := false

		// Iterate over each element in 'selectedCols'.
		for _, element := range selectedCols {
			if element == col {
				found = true
				break
			}
		}
		if found {
			selectedIndices = append(selectedIndices, i)
		}
	}

	// Return the slice containing the indices of selected columns.
	return selectedIndices
}
