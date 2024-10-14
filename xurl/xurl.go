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
	parsedURL, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		if rawUrl == "" {
			return false, xerr.ErrNotAbsolute // Return custom error for empty URLs
		}
		return false, fmt.Errorf("invalid URL: %v", err)
	}

	// check if scheme component of the URL is empty
	if !parsedURL.IsAbs() {
		return false, xerr.ErrNotAbsolute
	}

	// check if the scheme is neither http nor https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false, xerr.ErrWrongScheme
	}

	// if host is empty
	if parsedURL.Host == "" {
		return false, xerr.ErrEmptyHostName
	}

	// ensure host does not start with . or -
	if strings.HasPrefix(parsedURL.Host, ".") || strings.HasPrefix(parsedURL.Host, "-") {
		return false, xerr.ErrWrongHostFormat
	}

	// Check that the host contains at least one dot (valid domain format) or is localhost
	if !strings.Contains(parsedURL.Host, ".") && parsedURL.Host != "localhost" {
		return false, xerr.ErrInvalidDomainFormat
	}

	return true, nil
}

// tryFixScheme adds "http" or "https" if no scheme is provided, and retries validation
func TryFixScheme(rawUrl string) (bool, error) {
	// Try with "http" first
	fixedURL := "http://" + rawUrl
	parsedURL, err := url.ParseRequestURI(fixedURL)
	if err == nil && parsedURL.Scheme == "http" && parsedURL.Host != "" {
		return true, nil
	}

	// If http failed, try with "https"
	fixedURL = "https://" + rawUrl
	parsedURL, err = url.ParseRequestURI(fixedURL)
	if err == nil && parsedURL.Scheme == "https" && parsedURL.Host != "" {
		return true, nil
	}

	// If both fail, return an error
	return false, fmt.Errorf("could not fix URL: missing valid scheme")
}
