package download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/philhanna/go-ncvoters/util"
)

// ---------------------------------------------------------------------
// Functions
// ---------------------------------------------------------------------

// DownloadFile gets the data from the specified url and writes it to a
// file.
func DownloadFile(url, fileName string) error {
	const (
		MEGABYTE   = int64(1024 * 1024)
		BLOCK_SIZE = MEGABYTE
	)

	var err error

	length, err := GetContentLength(url)
	if err != nil {
		return err
	}
	progress := util.NewProgress()
	progress.Total = length
	progress.SoFar = 0
	progress.LastPercent = 0

	mb := float64(progress.Total) / float64(MEGABYTE)
	log.Printf("Downloading file (%.2fMB)...\n", mb)

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

	// Create a byte buffer with a size of one megabyte
	buffer := make([]byte, MEGABYTE)

	// Read from the response body and write to the file using the byte buffer
	stime := time.Now()
	for {

		// Read bytes from the response body into the buffer
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n <= 0 {
			break
		}
		progress.SoFar += int64(n)
		percent := int(float64(progress.SoFar) / float64(progress.Total) * 100)
		if percent != progress.LastPercent {
			s := strings.Repeat("*", percent/2)
			for len(s) < 50 {
				s += "."
			}
			if percent > progress.LastPercent {
				mb := float64(progress.SoFar) / float64(MEGABYTE)
				fmt.Printf("Percent complete: %d%%, [%-s] %.2fMB in %v\r",
					percent, s, mb, time.Since(stime))
			}
			progress.LastPercent = percent
		}

		// Write the bytes from the buffer to the file
		_, err = file.Write(buffer[:n])
		if err != nil {
			panic(err)
		}

		// Stop reading when we reach the end of the response body
		if err == io.EOF {
			break
		}
	}

	fmt.Println()
	log.Println("File downloaded successfully!")
	return nil
}

func NewProgress() {
	panic("unimplemented")
}
