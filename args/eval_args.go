// Package args package deals with detecting arguments and evaluating the arguments
// and parsing them to the intended functions.
package args

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
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
				slice, err := ReadUrlFromFile(path)
				if err != nil {
					errorss.WriteError(err, 2, false)
				}
				Arguments.Links = append(Arguments.Links, slice...)
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
			isValid, err := types.IsValidURL(arg)
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

// ReadUrlFromFile opens fpath to read the contents of the file (urls) and returns a slice of the urls
func ReadUrlFromFile(fpath string) (links []string, err error) {
	fd, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		link := strings.TrimSpace(scanner.Text())
		ok, err := types.IsValidURL(link)
		if err != nil {
			errorss.WriteError(err, 1, true)
		}
		if ok {
			links = append(links, link)
		}
	}
	return links, nil
}

// IsOutputFlag checks if -O=<filename> flag has been parsed with a valid filename and returns true
// and filename if successful else  returns false and empty string
func IsOutputFlag(arg string) (bool, string) {
	if strings.HasPrefix(arg, "-O=") {
		filename := strings.TrimSpace(strings.TrimPrefix(arg, "-O="))
		if filename == "" || filename == "-" || filename == ".." || filename == "." || strings.HasPrefix(filename, "/") {
			return false, ""
		} else {
			return true, filename
		}
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
		filename := matches[2]
		if filename == "." || filename == ".." || strings.HasPrefix(filename, "/") {
			return false, ""
		}
		return true, filename
	}
	return false, ""
}

// IsPathFlag returns true if -P flag has been used and a valid path has been specified
func IsPathFlag(s string) (bool, string) {
	pattern := `^-P=(.+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(s)
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
