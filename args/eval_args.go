// Package args package deals with detecting arguments and evaluating the arguments
// and parsing them to the intended functions.
package args

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"wget/ctx"
	"wget/fileio"
	"wget/help"
	"wget/xerr"
	"wget/xurl"
)

// DownloadContext builds and returns the download context,
// as defined by (parsing and evaluating) the commandline arguments.
func DownloadContext(arguments []string) (Arguments ctx.Context) {
	for _, arg := range arguments {

		switch {
		case arg == "--help":
			xerr.WriteError(help.Manual, 0, true)

		case arg == "-B":
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
					xerr.WriteError(err, 2, false)
				}

				// if we read the file and find no urls(empty file)
				if len(slice) == 0 {
					xerr.WriteError(fmt.Sprintf("No URLs found in %v", path), 2, false)
				}
				Arguments.Links = append(Arguments.Links, slice...)
			}

		case arg == "--mirror":
			Arguments.Mirror = true

		case arg == "--convert-links":
			Arguments.ConvertLinks = true

		case strings.HasPrefix(arg, "-O="):
			if ok, file := IsOutputFlag(arg); ok && file != "" {
				Arguments.OutputFile = file
			}

		case strings.HasPrefix(arg, "--rate-limit="):
			Arguments.RateLimit = strings.TrimPrefix(arg, "--rate-limit=")
			Arguments.RateLimitValue = ToBytes(Arguments.RateLimit)

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
			Arguments.Links = append(Arguments.Links, arg)
		}
	}
	return
}

// ToBytes converts a rateLimit in string format to bytes, if no suffix is supplied then the value is considered in bytes
// the only suffixes allowed are (k == kilobytes) and (M == megabytes)
// example when user passes: 20k toBytes returns 20000
// example when user passes: 20M toBytes returns 20000000
func ToBytes(rateLimit string) (rateLimitBytes int64) {
	// 1k == 1000 bytes
	// 1M == 1_000_000 bytes

	index := func(rateLimit []rune) int {
		for i := len(rateLimit) - 1; i >= 0; i-- {
			ch := rateLimit[i]
			if unicode.IsDigit(ch) {
				return i
			}
		}
		return -1
	}
	indx := index([]rune(rateLimit))
	size, suffix := rateLimit[:indx+1], rateLimit[indx+1:]

	sizeN, err := strconv.Atoi(size)
	if err != nil {
		fmt.Printf("Failed to convert size rate limit %s defaulting to 0\n", rateLimit)
		return 0
	}

	if strings.TrimSpace(suffix) == "" {
		return int64(sizeN)
	}
	if suffix == "k" {
		return int64(sizeN * 1000)
	} else if suffix == "M" {
		return int64(sizeN * 1000000)
	} else if suffix != "M" && suffix != "k" {
		fmt.Printf("Failed to convert size rate limit %s defaulting to 0\n", rateLimit)
		return 0
	}
	return int64(sizeN)
}

// ReadUrlFromFile opens fpath to read the contents of the file (urls) and returns a slice of the urls
func ReadUrlFromFile(fpath string) (links []string, err error) {
	fd, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer fileio.Close(fd)

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		link := strings.TrimSpace(scanner.Text())
		lnk, ok, err := xurl.IsValidURL(link)
		if err != nil {
			xerr.WriteError(err, 1, true)
		}
		if ok {
			links = append(links, lnk)
		}
	}
	return links, nil
}

// IsOutputFlag checks if -O=<filename> flag has been parsed with a valid filename and returns true
// and filename if successful else  returns false and empty string
func IsOutputFlag(arg string) (bool, string) {
	if strings.HasPrefix(arg, "-O=") {
		filename := strings.TrimSpace(strings.TrimPrefix(arg, "-O="))
		if filename == "" || filename == "-" || filename == ".." || filename == "." || strings.HasPrefix(
			filename, "/",
		) {
			return false, ""
		} else {
			return true, filename
		}
	}
	return false, ""
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
		xerr.WriteError(fmt.Sprintf("%v %s\n", xerr.ErrWrongPath, matches[1]), 1, true)
	}
	return true, matches[1]
}
