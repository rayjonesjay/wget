// Package httpx creates convenience utility wrappers around the standard's http package
package httpx

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ExtractMimeType extracts the MIME type of the response body, from HTTP response headers.
func ExtractMimeType(headers http.Header) string {
	contentType := headers.Get("Content-Type")
	if contentType == "" {
		return ""
	}

	// Content-Type can have parameters (e.g., "text/html; charset=utf-8")
	// We only want the main MIME type part.
	parts := strings.Split(contentType, ";")
	// since our sep argument to strings.Split is not empty, strings.Split must
	// return an array of at least length 1, thus we can safely access index 0
	mimeType := strings.TrimSpace(parts[0])

	return mimeType
}

// ExtractContentLength extracts the `Content-Length` from HTTP response headers as
// an int. Suppose, no `Content-Length` header exists in the response headers or
// the if the value doesn't make sense as an integer, then, -1 is returned
func ExtractContentLength(headers http.Header) int64 {
	contentLengthStr := headers.Get("Content-Length")
	if contentLengthStr == "" {
		return -1
	}

	contentLengthStr = strings.TrimSpace(contentLengthStr)
	contentLength, err := strconv.ParseInt(contentLengthStr, 10, 64)
	if err != nil {
		return -1
	}

	return contentLength
}

// FilenameFromContentDisposition returns the filename for the response contents, as dictated by
// HTTP `Content -Disposition` headers
func FilenameFromContentDisposition(headers http.Header) (string, error) {
	contentDisposition := headers.Get("Content-Disposition")
	if contentDisposition == "" {
		return "", errors.New("content-disposition header not found")
	}

	re := regexp.MustCompile(`(?i)filename\s*=\s*"?([^";]+)"?`)
	matches := re.FindStringSubmatch(contentDisposition)

	if len(matches) < 2 {
		return "", errors.New("filename parameter not found in content-disposition header")
	}

	return matches[1], nil
}
