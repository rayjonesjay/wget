package download

import (
	"fmt"
	"net/url"
	"strings"
)

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
