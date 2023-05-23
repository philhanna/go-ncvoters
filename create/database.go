package create

import (
	"database/sql"
	"log"
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
	
	// Get the zip file entry for the embedded CSV file
	zipEntry, err := GetZipEntry(zipFileName, csvFileName)
	if err != nil {
		log.Println(err)
		return err
	}

	// Open that CSV file entry
	csvReader, err := zipEntry.Open()
	if err != nil {
		log.Println(err)
		return err
	}
	defer csvReader.Close()
	
	/*
	
	zipEntry, err = os.Open(csvFileName)
	if err != nil {
		log.Println(err)
		return err
	}
	defer zipEntry.Close()
	
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
