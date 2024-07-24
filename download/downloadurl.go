package download

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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

// isValidURL checks if the given string is a valid URL
func IsValidURL(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	fmt.Println(parsedURL)
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
