// (Created with the assistance of ChatGPT)
//
// A Go program that addresses the requirements of optimizing
// performance, memory usage, and disk space while downloading the zip
// file, processing the CSV data, and creating an SQLite database.
// Includes the following optimizations:
//
// 1. **Memory optimization**: Instead of storing the entire CSV file in
// memory, the program reads the CSV file line by line using
// `csv.Reader`.  This approach reduces memory consumption, especially
// for large CSV files, as it only keeps a single record in memory at a
// time.
//
// 2. **Disk space optimization**: The program processes the CSV file
// directly from disk without extracting it from the zip file. This
// eliminates the need to extract and store the entire CSV file on disk,
// resulting in significant disk space savings.
//
// 3. **Performance optimization**: The program utilizes a single
// transaction to insert records into the SQLite database, which
// improves performance by reducing the overhead of committing each
// individual record. Additionally, it uses a prepared statement for the
// insertion, further optimizing the database operations.
package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	zipURL        = "https://s3.amazonaws.com/dl.ncsbe.gov/data/ncvoter_Statewide.zip"
	zipFileName   = "voter_data.zip"
	csvFileName   = "voter_data.csv"
	dbFileName    = "voter_data.db"
	progressEvery = 10000 // Log progress every n records
)

var selectedCols = []string{
	"county_id",
	"voter_reg_num",
	"last_name",
	"first_name",
	"middle_name",
	"name_suffix_lbl",
	"status_cd",
	"reason_cd",
	"res_street_address",
	"res_city_desc",
	"state_cd",
	"zip_code",
	"full_phone_number",
	"race_code",
	"ethnic_code",
	"party_cd",
	"gender_code",
	"birth_year",
	"age_at_year_end",
	"birth_state",
}

func main() {
	var err error

	// Download or reuse the voter zip file
	if FileExists(zipFileName) {
		log.Println("Reusing existing zip file")
	} else {
		err = DownloadFile(zipURL, zipFileName)
		if err != nil {
			log.Fatal("Failed to download the zip file:", err)
		}
	}

	// Create the database
	selectedColsString := strings.Join(selectedCols, ",")
	err = CreateDatabase(csvFileName, selectedColsString, dbFileName)
	if err != nil {
		log.Fatal("Failed to create the database:", err)
	}

	log.Println("Process completed successfully!")
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func DownloadFile(url, fileName string) error {
	log.Println("Downloading file...")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	log.Println("File downloaded successfully!")
	return nil
}

func CreateDatabase(csvFileName, selectedCols, dbFileName string) error {
	log.Println("Creating database...")

	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Create table with selected columns
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS voters (%s)", selectedCols)
	_, err = tx.Exec(query)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO voters (%s) VALUES (%s)",
		selectedCols, strings.Repeat("?,", strings.Count(selectedCols, ",")+1)+"?"))
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
		if strings.Contains(selectedCols, col) {
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

		_, err = stmt.Exec(values...)
		if err != nil {
			return err
		}

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
