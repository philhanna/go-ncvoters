package util

import (
	"archive/zip"
)

func IsGoodZipFile(filename string) bool {

	// Try to open the Zip file. If it fails, the file is corrupt or
	// non-existent
	reader, err := zip.OpenReader(filename)
	if err != nil {
		return false
	}
	defer reader.Close()

	// Iterate over the files in the Zip file to see if they can be
	// opened.
	for _, file := range reader.File {
		entry, err := file.Open()
		if err != nil {
			return false
		}
		defer entry.Close()
	}

	// Welp, looks good
	return true
}
