package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func DownloadUrl(url string) error {
	if !IsValidURL(url) {
		return fmt.Errorf("invalid-url")
	}

	// get the filename
	files := strings.Split(url, "/")
	filePath := "./" + files[len(files)-1]

	// send http request
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching file: %v", err)
	}

	defer response.Body.Close()

	// check if server returned status code 200
	if response.StatusCode == http.StatusOK {
		err = nil
	} else {
		return fmt.Errorf("%v", response.Status)
	}

	// create file to save content
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// save content to file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}
	return nil
}
