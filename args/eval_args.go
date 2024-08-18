// args package deals with detecting arguments and evaluating the arguments
// and passing them to different functions.
package args

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"wget/types"
)

// EvalArgs takes a slice of arguments and redirects them to the specific functions.
func EvalArgs(arguments []string) {
	// detecting go run . http://example.com
	for _, arg := range arguments {
		if isHelp(arg) {
			fmt.Println(types.Manual)
			os.Exit(0)
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
