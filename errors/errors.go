// errors package contains all error types, error messages and functions for handing errors
package errors

import "os"

func WriteError(errorMessage string, statusCode int) {
	os.Stdout.WriteString(errorMessage + "\n")
	os.Exit(statusCode)
}
