// Package download package contains functionalities for downloading a file, and also downloading and saving to a specified file name
package download

import (
	"strings"
)

// IsStringEmpty checks if string is empty
func IsStringEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
