// Package xerr defines common error types, error messages and functions for handing errors
package xerr

import (
	"errors"
	"fmt"
	"os"
)

// common error types are defined here with their respective messages
var (
	// ErrWrongScheme use when url scheme is not http or https
	ErrWrongScheme = errors.New("wrong url scheme")

	// ErrEmptyHostName use when host is empty
	ErrEmptyHostName = errors.New("host missing")

	// ErrNotAbsolute use when url is absolute
	ErrNotAbsolute = errors.New("url is not absolute")

	// ErrInvalidDomainFormat use when domain format is invalid
	ErrInvalidDomainFormat = errors.New("invalid domain format")

	// ErrWrongPath is returned when a given path isn't a valid location in the filesystem
	ErrWrongPath = errors.New("invalid path")

	ErrRelativeURL = errors.New("relative path")
)

// WriteError takes errorMessage of any type and statusCode
// then writes errorMessage to stderr and exits with the given statusCode if shouldExit is set to true
func WriteError(errorMessage interface{}, statusCode int, shouldExit bool) {
	writeError(errorMessage, statusCode, shouldExit, os.Exit)
}

func writeError(errorMessage interface{}, statusCode int, shouldExit bool, exit func(int)) {
	_, _ = os.Stderr.WriteString(fmt.Sprintf("%v\n", errorMessage))
	if shouldExit {
		exit(statusCode)
	}
}
