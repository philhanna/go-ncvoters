package ncvoters

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CreateDatabase(csvFileName, dbFileName string, progressEvery int) error {
	log.Println("Creating database...")

	// Open the database
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return err
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Create a table with selected columns
	query := CreateDDL()
	_, err = tx.Exec(query)
	if err != nil {
		return err
	}

	// Create a prepared statement for inserting records into the voters
	// table
	stmt, err := CreatePreparedStatement(tx)
	if err != nil {
		return err
	}
	defer stmt.Close()

	csvFile, err := os.Open(csvFileName)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1 // Allow varying number of fields

	columns, err := reader.Read()
	if err != nil {
		return err
	}

	selectedIndices := make([]int, 0)
	for i, col := range columns {
		found := false
		for _, element := range selectedCols {
			if strings.Contains(element, col) {
				found = true
				break
			}
		}
		if found {
			selectedIndices = append(selectedIndices, i)
		}
	}

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		values := make([]interface{}, len(selectedIndices))
		for i, idx := range selectedIndices {
			values[i] = record[idx]
		}

		/*
			_, err = stmt.Exec(values...)
			if err != nil {
				return err
			}
		*/

		count++
		if count%progressEvery == 0 {
			log.Printf("%d records inserted", count)
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("Database created successfully!")
	return nil
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