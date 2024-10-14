// Package xurl defines various utility functions that handle checking, parsing, and validating raw URLs.
// This package is named `xurl` to distinguish it from the go standard's `url` package that has a similar objective.
package xurl

import (
	"fmt"
	"net/url"
	"strings"

	"wget/xerr"
)

// IsValidURL checks if the given string is a valid URL
// IsValidURL checks if the given string is a valid URL
func IsValidURL(rawUrl string) (bool, error) {
	// Check if the scheme is missing
	if !strings.HasPrefix(rawUrl, "http://") && !strings.HasPrefix(rawUrl, "https://") {
		// If missing, prepend "http://" to the URL
		rawUrl = "http://" + rawUrl
	}

	// Attempt to parse the URL
	parsedURL, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return false, fmt.Errorf("invalid URL format: %v", err)
	}

	// Validate the scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false, xerr.ErrWrongScheme
	}

	// Check if the host is empty
	if parsedURL.Host == "" {
		return false, xerr.ErrEmptyHostName
	}

	// Ensure host does not start with . or -
	if strings.HasPrefix(parsedURL.Host, ".") || strings.HasPrefix(parsedURL.Host, "-") {
		return false, xerr.ErrWrongHostFormat
	}

	// Check for valid domain format (contains a dot or is localhost)
	if !strings.Contains(parsedURL.Host, ".") && parsedURL.Host != "localhost" {
		return false, xerr.ErrInvalidDomainFormat
	}

	// If everything passes, the URL is valid
	return true, nil
}
