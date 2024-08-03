package download

import (
	"fmt"
	"os"
	"path/filepath"
)

// ToFolder takes url, outputDir, and filename as arguments and saves a download to a specified directory
func ToFolder(url, outputDir, filename string) error {
	// Ensure the output directory exists
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	// Generate the full output path
	outputPath := filepath.Join(outputDir, filename)

	// Call ToFile function to save the file
	err = ToFile(url, outputPath)
	if err != nil {
		return fmt.Errorf("error saving file to folder: %v", err)
	}

	return nil
}
