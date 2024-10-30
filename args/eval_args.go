// Package args package deals with detecting arguments and evaluating the arguments
// and parsing them to the intended functions.
package args

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unicode"
	"wget/ctx"
	"wget/downloader"
	"wget/fileio"
	"wget/help"
	"wget/info"
	"wget/xerr"
	"wget/xurl"
)

// DownloadContext builds and returns the download context,
// as defined by (parsing and evaluating) the commandline arguments.
func DownloadContext(arguments []string) (Arguments ctx.Context) {

	length := len(arguments)

	// This closure function aids in preventing recursive launching of processes,
	// we only need one background option,
	// if you pass more than one, the function will consider the first background option
	withoutDuplicateBackground := func(args []string) []string {
		seen := make(map[string]bool)
		unique := []string{}

		for _, arg := range args {
			if arg == "-B" || arg == "--background" {
				if seen["background"] {
					continue // Skip duplicate -B or --background
				}
				seen["background"] = true
			}
			unique = append(unique, arg)
		}
		return unique
	}
	arguments = withoutDuplicateBackground(arguments)
	var backgroundLocation int

	for i, arg := range arguments {

		switch {
		case arg == "--help" || arg == "-h":
			xerr.WriteError(help.PrintManPage(), 0, true)

		case arg == "-B" || arg == "--background":
			Arguments.BackgroundMode = true
			if len(arguments) == 1 {
				xerr.WriteError(help.UsageMessage, 1, true)
			}
			backgroundLocation = i

		case strings.HasPrefix(arg, "-P="):
			isParsed, path := IsPathFlag(arg)
			if length == 1 {
				xerr.WriteError(help.UsageMessage, 1, true)
			}
			if isParsed {
				Arguments.SavePath = CreateDirFromPath(path)
			}

		case strings.HasPrefix(arg, "-i="):
			isParsed, path, err := InputFile(arg)
			if isParsed && err == nil {
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
			if length == 1 {
				xerr.WriteError(help.UsageMessage, 1, true)
			}
			if ok, file := IsOutputFlag(arg); ok && file != "" {
				Arguments.OutputFile = file
			}

		case strings.HasPrefix(arg, "--rate-limit="):
			if length == 1 {
				xerr.WriteError(help.UsageMessage, 1, true)
			}
			Arguments.RateLimit = strings.TrimPrefix(arg, "--rate-limit=")
			Arguments.RateLimitValue = ToBytes(Arguments.RateLimit)

		case strings.HasPrefix(arg, "-R="):
			if length == 1 {
				xerr.WriteError(help.UsageMessage, 1, true)
			}
			rejects := strings.Split(strings.TrimPrefix(arg, "-R="), ",")
			Arguments.Rejects = append(Arguments.Rejects, rejects...)

		case strings.HasPrefix(arg, "--reject="):
			if length == 1 {
				xerr.WriteError(help.UsageMessage, 1, true)
			}
			rejects := strings.Split(strings.TrimPrefix(arg, "--reject="), ",")
			Arguments.Rejects = append(Arguments.Rejects, rejects...)

		case strings.HasPrefix(arg, "-X="):
			if length == 1 {
				xerr.WriteError(help.UsageMessage, 1, true)
			}
			excludes := strings.Split(strings.TrimPrefix(arg, "-X="), ",")
			Arguments.Exclude = append(Arguments.Exclude, excludes...)

		case strings.HasPrefix(arg, "--exclude="):
			if length == 1 {
				xerr.WriteError(help.UsageMessage, 1, true)
			}
			excludes := strings.Split(strings.TrimPrefix(arg, "--exclude="), ",")
			Arguments.Exclude = append(Arguments.Exclude, excludes...)

		case arg == "--version" || arg == "-v":
			xerr.WriteError(info.VersionText(), 0, true)

		default:
			// if arg gets here and fails it wont be added to the urls flag
			url, isValid, err := xurl.IsValidURL(arg)
			if !isValid {
				xerr.WriteError(err, 1, false)
			}
			if isValid {
				Arguments.Links = append(Arguments.Links, url)
			}
		}
	}

	launchInBackground := func(args []string, index int, fd *os.File) {
		// remove the background flag to prevent recursion
		copy(args[index:], args[index+1:])
		args = args[:len(args)-1] // remove the last element since its will be a duplicate (same as len(args)-2)

		// detach the current process from its parent
		// run in the background
		executable, err := os.Executable()
		if err != nil {
			xerr.WriteError(fmt.Errorf("failed to get executable path:%v Defaulting to normal", err), 1, false)
		}

		// if stdin or stdout or stderr is nill, the process will read from os.DevNull
		cmd := exec.Command(executable, args...)

		// send the output of the file to specified file descriptor
		cmd.Stdout = fd
		os.Stdout = fd
		cmd.Stderr = fd
		cmd.Stdin = nil

		// set the command/program to run in a new session (child process)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true, // create a new session
		}

		// start the command and dont wait for it to finish
		// this allows us to get the terminal prompt back
		if err = cmd.Start(); err != nil {
			xerr.WriteError(fmt.Errorf("failed to start background process: %v Defaulting to normal", err), 1, false)
		}

		// This is critical, dont remove it
		os.Exit(0)
	}
	if Arguments.BackgroundMode {

		logFile := downloader.CheckIfFileExists("wget-log")
		fd, err := os.Create(logFile)
		if err != nil {
			xerr.WriteError(fmt.Errorf("failed to create %q defaulting to stdout", logFile), 2, false)
		}
		fmt.Printf("Output will be written to \"%s\".\n", logFile)
		launchInBackground(arguments, backgroundLocation, fd)
	}

	return
}

// CreateDirFromPath creates a directory from given dirPath and returns a clean path
// if dirPath does not exist it is created
func CreateDirFromPath(dirPath string) string {

	// fmt.Println(">>>", dirPath)
	user := os.ExpandEnv("$USER")
	// check for tilde '~` symbol and replace with
	if strings.HasPrefix(dirPath, "~") {
		dirPath = strings.Replace(dirPath, "~", "/home/"+user, 1)
	}
	// clean the path by:
	//removing more than one consecutive backslashes
	//replacing '.' with nothing
	//replacing '..' and replace all non-.. preceding it with nothing
	dirPath = path.Clean(dirPath)
	// fmt.Println(">>>", dirPath)
	absolutePath, _ := filepath.Abs(dirPath)
	// fmt.Println(absolutePath)
	err := os.MkdirAll(absolutePath, 0755)
	if err != nil {
		return "./"
	}

	return absolutePath
}

// ToBytes converts a rateLimit in (decimal or float) to bytes, if no suffix is supplied then the value is assumed to be bytes
// the only suffixes allowed are (k == kilobytes) and (M == megabytes)
// example when user passes: 20k ToBytes returns 20000
// example when user passes: 20M ToBytes returns 20000000
// example when user passes: 12.2 ToBytes returns 12
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

	sizeFloat, err := strconv.ParseFloat(size, 64)

	// Round of to nearest int less than sizeFloat
	sizeN := math.Floor(sizeFloat)
	if err != nil {
		log.Printf("Failed to convert size rate limit %s defaulting to 0\n", rateLimit)
		return 0
	}

	if strings.TrimSpace(suffix) == "" {
		return int64(sizeN)
	}
	if suffix == "k" {
		return int64(sizeN * 1000)
	} else if suffix == "M" {
		return int64(sizeN * 1000000)
	} else if suffix != "" {
		log.Printf("Unrecognized rate limit suffix %q defaulting to 0\n", rateLimit)
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
func InputFile(s string) (bool, string, error) {
	pattern := `^(-i=)(.+)`
	re := regexp.MustCompile(pattern)
	if re.MatchString(s) {
		matches := re.FindStringSubmatch(s)

		containsAll := func(s string, c string) bool {
			s = strings.TrimSpace(s)
			return strings.Count(s, c) == len(s)
		}
		filename := matches[2]

		if containsAll(filename, "/") {
			return false, "", errors.New("path is a directory")
		}
		if filename == "." || filename == ".." {
			return false, "", errors.New("path is a directory")
		}
		return true, filename, nil
	}
	return false, "", errors.New("path might be empty")
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
