// download package contains functionalites for downloading a file, and also downloading and saving to a specified file name
package download

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"wget/types"
)

// RenamingFile function accepts url and filename as arguments and saves downloaded file using the provided name
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

func DownloadUrl(url string) error {
	if !IsValidURL(url) {
		return fmt.Errorf("invalid-url")
	}

	Feedback(url)

	// get the filename
	files := strings.Split(url, "/")
	filePath := "./" + files[len(files)-1]

	// send http request
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching file: %v", err)
	}

	defer response.Body.Close()

	reqMessage := "sending request, awaiting response..."

	if response.StatusCode == http.StatusOK {

		fmt.Printf("\r"+reqMessage+" status %d OK\n", response.StatusCode)
		err = nil

	} else if response.StatusCode == http.StatusNotFound {

		fmt.Printf(reqMessage+" %d: Not Found\n", response.StatusCode)
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	contentLength, err := io.Copy(file, response.Body)
	roundOfSize := RoundOfSizeOfData(contentLength)
	fmt.Printf("content size: %d %s\n", contentLength, roundOfSize)
	if err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}

	fmt.Printf("saving file to %s\n", file.Name())
	// time.Sleep(12 * time.Second)
	fmt.Printf("finished at %s\n", GetCurrentTime())
	return nil
}

// isValidURL checks if the given string is a valid URL
func IsValidURL(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check if the scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// Check if the host is not empty
	if parsedURL.Host == "" {
		return false
	}

	// Ensure host is not empty and does not start with a dot or hyphen
	if parsedURL.Host == "" || strings.HasPrefix(parsedURL.Host, ".") || strings.HasPrefix(parsedURL.Host, "-") {
		return false
	}

	// Check if the host contains at least one dot
	if !strings.Contains(parsedURL.Host, ".") {
		return false
	}
	return true
}

// RoundOfSizeData converts dataInBytes (size of file downloaded) in bytes to the nearest size
func RoundOfSizeOfData(dataInBytes int64) string {
	var size float64
	var unit string

	if dataInBytes >= types.KB {
		size = float64(dataInBytes) / types.MB
		unit = "MB"
	} else {
		size = float64(dataInBytes)
		unit = "KB"
	}

	return fmt.Sprintf("%.2f%s", size, unit)
}
