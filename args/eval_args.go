// Package args package deals with detecting arguments and evaluating the arguments
// and parsing them to the intended functions.
package args

import (
	"fmt"
	"regexp"
	"strings"
	"wget/download"

	"wget/errorss"
	"wget/help"
	"wget/types"
)

// EvalArgs takes a slice of strings parsed at the command line and adds them to types.Arg
// as field values and returns types.Arg
func EvalArgs(arguments []string) (Arguments types.Arg) {
	for _, arg := range arguments {
		switch {
		case IsHelpFlag(arg):
			errorss.WriteError(help.UsageMessage, 0, true)

		case IsBackgroundFlag(arg):
			Arguments.BackgroundMode = true

		case strings.HasPrefix(arg, "-P="):
			isParsed, path := IsPathFlag(arg)
			if isParsed {
				Arguments.SavePath = path
			}

		case strings.HasPrefix(arg, "-i="):
			isParsed, path := InputFile(arg)
			if isParsed {
				Arguments.InputFile = path
			}

		case arg == "--mirror":
			Arguments.Mirror = true

		case IsConvertLinksOn(arg):
			Arguments.ConvertLinks = true

		case strings.HasPrefix(arg, "-O="):
			if ok, file := IsOutputFlag(arg); ok && file != "" {
				Arguments.OutputFile = file
			}

		case strings.HasPrefix(arg, "--rate-limit="):
			Arguments.RateLimit = strings.TrimPrefix(arg, "--rate-limit=")

		case strings.HasPrefix(arg, "-R="):
			rejects := strings.Split(strings.TrimPrefix(arg, "-R="), ",")
			Arguments.Rejects = append(Arguments.Rejects, rejects...)

		case strings.HasPrefix(arg, "--reject="):
			rejects := strings.Split(strings.TrimPrefix(arg, "--reject="), ",")
			Arguments.Rejects = append(Arguments.Rejects, rejects...)

		case strings.HasPrefix(arg, "-X="):
			excludes := strings.Split(strings.TrimPrefix(arg, "-X="), ",")
			Arguments.Exclude = append(Arguments.Exclude, excludes...)

		case strings.HasPrefix(arg, "--exclude="):
			excludes := strings.Split(strings.TrimPrefix(arg, "--exclude="), ",")
			Arguments.Exclude = append(Arguments.Exclude, excludes...)

		default:
			isValid, err := download.IsValidURL(arg)
			if err != nil {
				errorss.WriteError(err, 1, true)
			}
			if isValid {
				Arguments.Links = append(Arguments.Links, arg)
			}
		}
	}
	return
}

// IsOutputFlag checks if -O=<filename> flag has been parsed with a valid filename and returns true
// and filename if successful else  returns false and empty string
func IsOutputFlag(arg string) (bool, string) {
	if strings.HasPrefix(arg, "-O=") {
		filename := strings.TrimPrefix(arg, "-O=")
		return true, filename
	}
	return false, ""
}

// IsConvertLinksOn checks if --convert-links argument is parsed, in order to determine whether
// the links will be converted for local viewing
func IsConvertLinksOn(arg string) bool {
	return strings.HasPrefix(arg, "--") && strings.Contains(arg, "convert-links")
}

// InputFile recognizes if -i=<filename> has been parsed together with a valid filename that exist
// if it does not exist or is empty returns false and empty string
func InputFile(s string) (bool, string) {
	pattern := `^(-i=)(.+)`
	re := regexp.MustCompile(pattern)
	if re.MatchString(s) {
		matches := re.FindStringSubmatch(s)
		filename := matches[1]
		return true, filename
	}
	return false, ""
}

// IsPathFlag returns true if -P flag has been used and a valid path has been specified
func IsPathFlag(s string) (bool, string) {
	pattern := `^-P=(.+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(s, 1)
	if matches == nil {
		return false, ""
	}
	if matches[1] == "." || matches[1] == ".." {
		errorss.WriteError(fmt.Sprintf("%v %s\n", errorss.ErrWrongPath, matches[1]), 1, true)
	}
	return true, matches[1]
}

// IsBackgroundFlag returns true if -B has been parsed in the command line,else false
func IsBackgroundFlag(s string) bool {
	s = strings.ToUpper(s)
	pattern := `^-B`
	re := regexp.MustCompile(pattern)
	return re.MatchString(s)
}

// IsHelpFlag detects if any of the flags parsed has the --help
// format which displays how to use the program.
func IsHelpFlag(argument string) bool {
	argument = strings.ToLower(strings.TrimSpace(argument))

	helpPattern := `^--help$`

	re := regexp.MustCompile(helpPattern)

	return re.MatchString(argument)
}
