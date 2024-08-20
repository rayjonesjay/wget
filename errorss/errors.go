// errors package contains all error types, error messages and functions for handing errors
package errorss

import (
	"fmt"
	"os"
)

// all error types are defined here and their respective messages, because repetition is boring
var (
	// use when url scheme is not http or https
	ErrWrongSheme = fmt.Errorf("wrong url scheme")

	// use when host is empty
	ErrEmptyHostName = fmt.Errorf("host missing")

	// use when url is absolute
	ErrNotAbsolute = fmt.Errorf("url is not absolute")

	// use when domain format is invalid
	ErrInvalidDomainFormat = fmt.Errorf("invalid domain format")
)

// WriteError takes errorMessage of any type and statusCode and writes errorMessage to stdout and exits with statusCode
func WriteError(errorMessage interface{}, statusCode int) {
	os.Stdout.WriteString(fmt.Sprintf("%v\n", errorMessage))
	os.Exit(statusCode)
}
