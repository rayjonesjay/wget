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
func IsValidURL(rawUrl string) (string, bool, error) {
	if !strings.HasPrefix(rawUrl, "http://") && !strings.HasPrefix(rawUrl, "https://") {
		rawUrl = "https://" + rawUrl
	}
	parsedURL, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return "", false, fmt.Errorf("invalid URL: %v", err)
	}

	// check if scheme component of the URL is empty
	if !parsedURL.IsAbs() {
		return "", false, xerr.ErrNotAbsolute
	}

	// check if the scheme is neither http nor https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", false, xerr.ErrWrongScheme
	}

	// if host is empty
	if parsedURL.Host == "" {
		return "", false, xerr.ErrEmptyHostName
	}

	// ensure host does not start with . or -
	if strings.HasPrefix(parsedURL.Host, ".") || strings.HasPrefix(parsedURL.Host, "-") {
		return "", false, fmt.Errorf("wrong host format %q", parsedURL.Host)
	}

	// Check that the host contains at least one dot (valid domain format) or is localhost
	if !strings.Contains(parsedURL.Host, ".") && parsedURL.Host != "localhost" {
		return "", false, xerr.ErrInvalidDomainFormat
	}

	return rawUrl, true, nil
}
