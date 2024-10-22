// Package httpx creates convenience utility wrappers around the standard's http package
package httpx

import (
	"errors"
	"fmt"
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

// RoundOfSizeOfData  converts dataInBytes (size of file downloaded) in bytes to the nearest size
//
// Deprecated: Use [globals.FormatSize] instead
func RoundOfSizeOfData(dataInBytes int64) string {
	var size float64
	var unit string
	if dataInBytes >= GB {
		size = float64(dataInBytes) / GB
		unit = "GB"
	} else if dataInBytes >= KB {
		size = float64(dataInBytes) / MB
		unit = "MB"
	} else {
		size = float64(dataInBytes)
		unit = "KB"
	}
	return fmt.Sprintf("%.2f%s", size, unit)
}

const (
	// KB Size of Data in KiloBytes
	KB = 1000 * 1

	// KiB Size of Data in KibiBytes, same as 2^10
	KiB = 1 << (10 * 1)

	// MB Size of Data in MegaBytes
	MB = 1000 * KB

	// MiB Size of Data in MebiBytes, same as 2^20
	MiB = 1 << 20

	// GB Size of Data in GigaBytes
	GB = 1000 * MB

	// GiB Size of Data in GibiBytes, same as 2^30
	GiB = 1 << 30

	// TB Size of Data in TeraBytes
	TB = 1000 * GB

	// TiB Size of Data in TebiBytes, same as 2^40
	TiB = 1 << 40
)
