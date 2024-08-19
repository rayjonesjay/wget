// args package deals with detecting arguments and evaluating the arguments
// and parsing them to the intended functions.
package args

import (
	"regexp"
	"strings"

	"wget/download"
	"wget/errorss"
	"wget/types"
)

// EvalArgs takes a slice of arguments and redirects them to the specific functions.
func EvalArgs(arguments []string) {
	// if one argument is passed:
	if len(arguments) == 1 {
		arg := arguments[0]

		// display the man page of our program and exit with 0 status code
		if isHelp(arg) {
			errorss.WriteError(types.Manual, 0)
		}

		// if the argument passed is a valid url, then call the downldoad function
		isValid, err := download.IsValidURL(arg)
		if isValid {
			download.DownloadUrl(arg)
		} else {
			errorss.WriteError(err, 1)
		}
	}
}

// isHelp detects if any of the flags parsed has the --help format which displays the how to use the program.
func isHelp(argument string) bool {
	argument = strings.ToLower(strings.TrimSpace(argument))

	helpPattern := `^--help$`

	re := regexp.MustCompile(helpPattern)

	return re.MatchString(argument)
}
