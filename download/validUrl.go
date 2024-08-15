package download

import (
	"net/url"
	"strings"
)

// isValidURL checks if the given string is a valid URL
func IsValidURL(urlStr string) bool {
	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false
	}

	// Ensure the URL is absolute and has a valid scheme (http or https)
	if !parsedURL.IsAbs() || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return false
	}

	// Ensure the host is not empty, and doesn't start with a dot or hyphen
	if strings.HasPrefix(parsedURL.Host, ".") || strings.HasPrefix(parsedURL.Host, "-") || parsedURL.Host == "" {
		return false
	}

	// Check that the host contains at least one dot (valid domain format)
	if !strings.Contains(parsedURL.Host, ".") {
		return false
	}

	return true
}
