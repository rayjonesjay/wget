package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

/* RenamingFile function accepts url and filename as arguments and saves downloaded file using the provided name */
func RenamingFile(url, filename string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error downloading file: %v\n", err)
	}

	defer response.Body.Close()

	output, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Error creating file: %v\n", err)
	}

	defer output.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return fmt.Errorf("Error saving file: %v\n", err)
	}
	return nil
}
