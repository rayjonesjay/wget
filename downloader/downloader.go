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
	"time"
	"wget/ctx"
	"wget/fetch"
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

var (
	// DefaultDownloadDirectory is the default location where files retrieved will reside.
	DefaultDownloadDirectory = "$HOME/Downloads"

	// DefaultHTTPPort is the default port used by http if not specified.
	DefaultHTTPPort = "80"

	// DefaultHTTPSPort is the default port used by https if not specified.
	DefaultHTTPSPort = "443"
)

// Get downloads any files, website mirrors, or resources as defined by the provided download context
func Get(c ctx.Context) {
	a := arg{Context: &c}
	err := a.Get()
	if err != nil {
		// TODO: handle errors during downloads
		panic(err)
	}
}

var width = syscheck.GetTerminalWidth()
var mu sync.Mutex

// Download handles each download and prints progress across 6 lines.
func (a *arg) Download() error {
	var wg sync.WaitGroup

	successfulDownloads := make(chan string, len(a.Links))

	syscheck.ClearScreen()
	syscheck.HideCursor()
	defer syscheck.ShowCursor() // Ensure cursor is shown again when done

	for lineNumber, url := range a.Links {

		wg.Add(1)
		go func(url string, rowOffset int) {
			defer wg.Done()

			startTime := time.Now()

			outputFilePath := CheckIfFileExists(a.determineOutputPath(url))

			GetFile := func(downloadUrl string, header http.Header) (*os.File, error) {
				return os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			}

			_, err := fetch.URL(url, fetch.Config{
				GetFile: GetFile,
				Limit:   0,
				ProgressListener: func(downloaded, total int64) {
					progress := float64(downloaded) / float64(total)
					barLength := width / 3
					filled := int(progress * float64(barLength))
					notFilled := barLength - filled

					bar := fmt.Sprintf("[%s%s]", strings.Repeat("=", filled), strings.Repeat(" ", notFilled))
					percentage := (downloaded * 100) / total
					eta := CalculateETA(downloaded, total, time.Since(startTime))

					mu.Lock() // Avoid race conditions on terminal output
					PrintLines(rowOffset, []string{
						fetch.Status.Start,
						fetch.Status.Status + fmt.Sprintf("%d", fetch.Status.StatusCode),
						fetch.Status.ContentLength,
						fmt.Sprintf("saving file to: %s", outputFilePath),
						fmt.Sprintf("%s / %s %s %d%% %s", formatSize(downloaded), formatSize(total), bar, percentage, eta),
					})
					mu.Unlock()
				},
				RateListener:       func(rate int32) {},
				Body:               nil,
				Method:             "GET",
				AllowedStatusCodes: []int{http.StatusOK},
			})

			if err != nil {
				mu.Lock()
				PrintLines(rowOffset+7, []string{
					fmt.Sprintf("Error: %s", err.Error()),
				})
				mu.Unlock()
			} else {
				successfulDownloads <- url
				mu.Lock()
				PrintLines(rowOffset+5, []string{
					syscheck.GetCurrentTime(false),
				})
				mu.Unlock()
			}
		}(url, lineNumber*7) // Reserve 7 lines per download
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
	if len(successList) != 0 {
		// Print the successfully downloaded URLs
		printUrls(successList)
	}

	fmt.Print("\033[?25h")
	return nil
}

// PrintLines prints multiple lines of text starting from a specific row.
func PrintLines(baseRow int, lines []string) {
	for i, line := range lines {
		syscheck.MoveCursor(baseRow + i)
		fmt.Print("\r", line)
	}
}

// CalculateETA estimates the time left for a download.
func CalculateETA(downloaded, total int64, elapsed time.Duration) string {
	if downloaded == 0 {
		return "calculating..."
	}
	rate := float64(downloaded) / elapsed.Seconds()
	remaining := float64(total-downloaded) / rate
	return formattedTime(int64(remaining)) // Just return the formatted time.
}

// Format time into a human-readable string.
func formattedTime(seconds int64) string {
	if seconds >= 3600 {
		return fmt.Sprintf("%d h", seconds/3600)
	} else if seconds >= 60 {
		return fmt.Sprintf("%d m", seconds/60)
	} else {
		return fmt.Sprintf("%d s", seconds)
	}
}

// CheckIfFileExists will check if fname exists in the provided path if it exists it will add
// a number starting from 1 between the filename and the beginning of extension
// example: if file.txt exist CheckIfFileExist will generate a new name file1.txt. It does this iteratively.
func CheckIfFileExists(fname string) string {

	if strings.TrimSpace(fname) == "" {
		return ""
	}
	// get the file extension
	extension := filepath.Ext(fname)
	base := fname
	if extension != "" {
		base = fname[:len(fname)-len(extension)]
	}
	n := 1

	for {
		_, err := os.Stat(fname)
		if os.IsNotExist(err) {
			return fname
		}

		fname = fmt.Sprintf("%s%d%s", base, n, extension)
		n++
	}
}

/*
start at 2017-10-14 03:46:06
sending request, awaiting response... status 200 OK
content size: 56370 [~0.06MB]
saving file to: ./meme.jpg
 55.05 KiB / 55.05 KiB [================================================================================================================] 100.00% 1.24 MiB/s 0s

Downloaded [https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg]
finished at 2017-10-14 03:46:07
*/
// determineOutputPath determines the full path for the output file
func (a *arg) determineOutputPath(url string) string {
	var outputFilePath string

	if a.OutputFile != "" {
		outputFilePath = a.OutputFile
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

func printUrls(urls []string) {
	res := ""
	for i, url := range urls {
		if i != len(urls)-1 {
			res += url + " "
		} else {
			res += url
		}
	}
	fmt.Printf("\nDownloaded	[%s]\n", res)
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

func (a *arg) MirrorWeb() error {

	for _, link := range a.Links {
		err := mirror.Site(*a.Context, link)
		if err != nil {
			return err
		}
		//parsedUrl, err := url.Parse(link)
		//if err != nil {
		//	return err
		//}
		//
		//// by default the downloaded data will be saved to the name of domain if not specified
		//domain := parsedUrl.Host
		//directoryToSaveData := filepath.Join(a.SavePath, domain)
		//
		//err = os.MkdirAll(directoryToSaveData, 0755)
		//if err != nil {
		//	return err
		//}

		// Download and parse the HTML/CSS
		//err = a.downloadAndParseHTML(link, directoryToSaveData)
		//if err != nil {
		//	return err
		//}
	}
	return nil
}

func (a *arg) Get() (err error) {
	if a.Mirror {
		// run in mirror mode
		err = a.MirrorWeb()
	} else {
		// regular download
		err = a.Download()
	}
	return
}

// Helper function to format byte size
func formatSize(size int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GiB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MiB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KiB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
