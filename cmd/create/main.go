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
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/philhanna/go-ncvoters/create"
	"github.com/philhanna/go-ncvoters/download"
	"github.com/philhanna/go-ncvoters/util"
)

func main() {

	const (
		zipURL        = "https://s3.amazonaws.com/dl.ncsbe.gov/data/ncvoter_Statewide.zip"
		entryName     = "ncvoter_Statewide.txt"
		progressEvery = 100_000 // Log progress every n records
	)

	var (
		zipFileName = filepath.Join(os.TempDir(), "voter_data.zip")
		dbFileName  = filepath.Join(os.TempDir(), "voter_data.db")
	)

	var err error

	// Set the logging flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Starting voter database creation")
	
	// Download or reuse the voter zip file
	reuse := false
	if util.FileExists(zipFileName) && util.IsGoodZipFile(zipFileName) {
		reuse = true
	}

	if reuse {
		log.Println("Reusing existing zip file")
	} else {
		err = download.DownloadFile(zipURL, zipFileName)
		if err != nil {
			log.Fatal("Failed to download the zip file:", err)
		}
	}

	// Create the database
	err = create.CreateDatabase(zipFileName, entryName, dbFileName, progressEvery)
	if err != nil {
		log.Fatal("Failed to create the database:", err)
	}

	log.Println("Process completed successfully!")
}
