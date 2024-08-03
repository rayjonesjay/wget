// download package contains functionalites for downloading a file, and also downloading and saving to a specified file name
package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// ToFile accepts url and filename as arguments and saves downloaded file using the provided name
func ToFile(url, filename string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading file: %v", err)
	}

	defer response.Body.Close()

	output, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	defer output.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}
	return nil
}
