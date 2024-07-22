package savetofolder

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// ToFolder function takes url, outPutDir and filename as arguments and saves a download to a specified directory
func ToFolder(url, outputDir, filename string) error {
	// Check if the outputDir file exists
	err0 := os.MkdirAll(outputDir, os.ModePerm)
	if err0 != nil {
		return fmt.Errorf("error creating directory: %v", err0)
	}

	outputPath := filepath.Join(outputDir, filename)

	response, err1 := http.Get(url)

	if err1 != nil {
		return fmt.Errorf("error fetching data:%v", err1)
	}

	defer response.Body.Close()

	outputFile, err2 := os.Create(outputPath)

	if err2 != nil {
		return fmt.Errorf("error creating directory: %v", err2)
	}

	defer outputFile.Close()

	_, err3 := io.Copy(outputFile, response.Body)

	if err3 != nil {
		return fmt.Errorf("error saving data to the provided path: %v", err3)
	}
	return nil
}
