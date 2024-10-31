// Package downloader contains user defined downloader and constants
package downloader

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"wget/ctx"
	"wget/fetch"
	"wget/globals"
	"wget/mirror"
	"wget/syscheck"
)

// arg represents the commandline arguments passed through the command line by the user
// example: $ go run . -O=file.txt -B https://www.example.com
// -0 will be a field in arg that specifies the output file to save the resource
// -B will send the download process to background mode
// https://www.example.com is the link to where the resource resides
type arg struct {
	*ctx.Context
}

// Get downloads any files, website mirrors, or resources as defined by the provided download context
func Get(c ctx.Context) {
	a := arg{Context: &c}
	var err error
	var dType string
	if a.Mirror {
		// run in mirror mode
		err = a.MirrorWeb()
		dType = "mirror"
	} else {
		// regular download
		err = a.Download()
		dType = "download"
	}
	if err != nil {
		syscheck.ShowCursor()
		fmt.Printf("\n%s failed: %v\n", dType, err)
	}
}

// Download handles each download and prints progress across 6 lines.
func (a *arg) Download() error {
	var wg sync.WaitGroup
	successfulDownloads := make(chan string, len(a.Links))

	syscheck.MoveCursor(1)
	syscheck.ClearScreen()
	syscheck.HideCursor()
	defer syscheck.ShowCursor() // Ensure cursor is shown again when done

	// how many rows each download progress indicator for a given link is allocated
	rows := 9
	for i, url := range a.Links {
		lineNumber := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			outputFilePath := CheckIfFileExists(a.determineOutputPath(url))

			GetFile := func(downloadUrl string, header http.Header) (*os.File, error) {
				return os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE, 0644)
			}

			// define a download status listener for the current mirror URL
			status := fetch.DownloadStatus{}
			status.OnUpdate = func(status *fetch.DownloadStatus, hint int) {
				// whenever the status of this download has changed, we print the progress to its
				// respective position in the terminal
				globals.PrintLines((lineNumber*rows)+hint, []string{status.GetField(hint)})
			}
			globals.PrintLines(lineNumber*rows, globals.StringTimes("\033[38;5;208m···\033[0m", rows-1))

			// configure an Advanced Progress Listener for the GET request
			advancedProgressListener := *status.ProgressListener()
			_, err := fetch.URL(
				url,
				fetch.Config{
					GetFile:                  GetFile,
					Limit:                    int32(a.RateLimitValue),
					ProgressListener:         nil,
					RateListener:             nil,
					Body:                     nil,
					Method:                   "GET",
					AllowedStatusCodes:       []int{http.StatusOK},
					AdvancedProgressListener: advancedProgressListener,
				},
			)

			if err != nil {
				errString := fmt.Sprintf("error: %s", err.Error())
				errString = strings.Replace(errString, "\n", " : ", -1)
				globals.PrintLines((lineNumber*rows)+5, []string{errString})
			} else {
				if lineNumber != len(a.Links) {
					successfulDownloads <- url
				}
			}
		}()
	}

	// wait for all go routines to finish in order to close the channel
	go func() {
		wg.Wait()
		close(successfulDownloads)
	}()
	// Collect all successful downloads
	var successList []string
	for url := range successfulDownloads {
		successList = append(successList, url)
	}

	if len(a.Links) != 0 {
		// move the cursor to the terminal row after the last progress
		syscheck.MoveCursor(len(a.Links) * rows)
	}

	if len(successList) > 1 {
		// Print the successfully downloaded URLs
		fmt.Printf("\nDownloads finished:\t%v\n", successList)
	}

	return nil
}

// CheckIfFileExists checks if a file with the provided name exists. If it exists, it will add
// a number starting from 1 between the filename and the beginning of extension
// example: if file.txt exist CheckIfFileExist will generate a new name file1.txt. It does this iteratively.
func CheckIfFileExists(filename string) string {
	if strings.TrimSpace(filename) == "" {
		return ""
	}
	// get the file extension
	extension := filepath.Ext(filename)
	base := filename
	if extension != "" {
		base = filename[:len(filename)-len(extension)]
	}
	n := 1

	for {
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			return filename
		}

		filename = fmt.Sprintf("%s(%d)%s", base, n, extension)
		n++
	}
}

// determineOutputPath determines the full path for the output file
func (a *arg) determineOutputPath(url string) string {
	var outputFilePath string

	if a.OutputFile != "" {
		outputFilePath = filepath.Join(a.SavePath, a.OutputFile)
	} else {
		// get the filename from the url
		tokens := strings.Split(url, "/")
		filename := tokens[len(tokens)-1]

		// If SavePath is specified use it otherwise use the current directory
		if a.SavePath != "" {
			outputFilePath = filepath.Join(a.SavePath, filename)
		} else {
			outputFilePath = filename
		}
	}
	return outputFilePath
}

// IsEmpty function checks whether an iterable is empty, an iterable is a string,array or slice.
// it returns true if the `data` which is expected to be an iterable, is empty else return false
func (a *arg) IsEmpty(data interface{}) bool {

	// get the value of the interface
	val := reflect.ValueOf(data)

	// determine which kind of iterable object it is.

	switch val.Kind() {

	// if it's a string, check if its empty
	case reflect.String:
		return strings.TrimSpace(val.String()) == ""

	// if it's an array or a slice
	case reflect.Array, reflect.Slice:
		return val.Len() == 0

	// for any other type return false, meaning it's not empty or object is not an iterable
	default:
		return false
	}
}

func (a *arg) MirrorWeb() (gErr error) {
	for _, link := range a.Links {
		err := mirror.Site(*a.Context, link)
		if err != nil {
			gErr = err
			continue
		}
	}
	return
}
