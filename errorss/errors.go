// Package errorss errors package contains all error types, error messages and functions for handing errors
package errorss

import (
	"errors"
	"fmt"
	"os"
)

// all error types are defined here and their respective messages, because repetition is boring
var (
	// ErrWrongScheme use when url scheme is not http or https
	ErrWrongScheme = errors.New("wrong url scheme")

	// ErrEmptyHostName use when host is empty
	ErrEmptyHostName = errors.New("host missing")

	// ErrNotAbsolute use when url is absolute
	ErrNotAbsolute = errors.New("url is not absolute")

	// ErrInvalidDomainFormat use when domain format is invalid
	ErrInvalidDomainFormat = errors.New("invalid domain format")

	ErrWrongPath = errors.New("invalid path")
)

// WriteError takes errorMessage of any type and statusCode
// then writes errorMessage to stderr and exits with the given statusCode if shouldExit is set to true
func WriteError(errorMessage interface{}, statusCode int, shouldExit bool) {
	_, err := os.Stderr.WriteString(fmt.Sprintf("%v\n", errorMessage))
	if err != nil {
		return
	}
	if shouldExit {
		os.Exit(statusCode)
	}
	return
}
