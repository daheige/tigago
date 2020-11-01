package gutils

import (
	"log"
	"os"
	"testing"
)

// TestFgetcsv test csv read.
func TestFgetcsv(t *testing.T) {
	file, _ := os.Open("test.csv")
	fileInfo, _ := file.Stat()
	log.Println("file size: ", fileInfo.Size())
	res, err := Fgetcsv(file, ',')
	log.Println(res, err)
}
