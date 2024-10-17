// Package downloader contains user defined downloader and constants
package downloader

import (
	"fmt"
	"io"
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
	"wget/xerr"
)

// arg represents the commandline arguments passed through the command line by the user
// example: $ go run . -O=file.txt -B https://www.example.com
// -0 will be a field in arg that specifies the output file to save the resource
// -B will send the download process to background mode
// https://www.example.com is the link to where the resource resides
type arg struct {
	*ctx.Context
}

const (
	// KB Size of Data in KiloBytes
	KB = 1000 * 1

	// KiB Size of Data in KibiBytes, same as 2^10
	KiB = 1 << (10 * 1)

	// MB Size of Data in MegaBytes
	MB = 1000 * KB

	// MiB Size of Data in MebiBytes, same as 2^20
	MiB = 1 << 20

	// GB Size of Data in GigaBytes
	GB = 1000 * MB

	// GiB Size of Data in GibiBytes, same as 2^30
	GiB = 1 << 30

	// TB Size of Data in TeraBytes
	TB = 1000 * GB

	// TiB Size of Data in TebiBytes, same as 2^40
	TiB = 1 << 40
)

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

// Download handles normal downloads based on the provided URLs and other flags in the arg struct
func (a *arg) Download() error {
	var m sync.Mutex
	successfulDownloads := make(chan string, len(a.Links))

	var wg sync.WaitGroup
	syscheck.ClearScreen()      //clear screen before download begins
	defer syscheck.ShowCursor() // show the cursor if the download is finished
	for lineNumber, url := range a.Links {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			GetFile := func(downloadUrl string, header http.Header) (*os.File, error) {
				outputFilePath := CheckIfFileExists(a.determineOutputPath(url))
				return os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			}

			_, err := fetch.URL(url, fetch.Config{
				GetFile: GetFile,
				Limit:   0,
				ProgressListener: func(downloaded, total int64) {
					// the size of the progress bar will be determined at runtime depending on terminal size, but default is 112 according to project requirements

					var barLength float64 = float64(width) / 2

					// filled tells us how many equal signs to print
					progress := float64(downloaded) / float64(total)

					filled := (progress * barLength)
					// not filled will represent the empty part thats not filled with equal signs
					notFilled := barLength - filled

					notFilledString := ""
					if notFilled > 0 {
						notFilledString = strings.Repeat(" ", int(notFilled))
					}

					bar := fmt.Sprintf("[%s%s]", strings.Repeat("=", int(filled)), notFilledString)
					percentage := (downloaded / total) * 100

					// hide cusor visibility
					fmt.Print("\033[?25l")

					a := fmt.Sprintf("\r%.2f / %.2f %s %d%%  ", float64(downloaded), float64(total), bar, percentage)
					if width >= 65 {
						m.Lock()                            // Since each routine is accessing the terminal we need to lock it to prevent race conditions
						syscheck.MoveCursor(lineNumber + 5) // each download gets a size of 5 lines to print its output
						fmt.Print(a)
						m.Unlock()
					}
				},
				RateListener: func(rate int32) {
				},
				Body:               nil,
				Method:             "GET",
				AllowedStatusCodes: []int{http.StatusOK},
			})
			if err != nil {
				fmt.Println(err)
			} else {
				successfulDownloads <- url
			}
		}(url)
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

	// enable cusor visibility
	fmt.Print("\033[?25h")
	return nil
}

// ETA calculates the downloaded data and current internet speed, the estimated time download finishes
func ETA(total, downloaded, speed int64) int64 {

	remainingTime := (total / downloaded) / speed

	return remainingTime
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

// GetResource takes an url to a resource and attempts to fetch the specified resource, if an error occurs err is returned
func (a *arg) GetResource(url string) (err error) {
	outputFilePath := CheckIfFileExists(a.determineOutputPath(url))
	// open the output file for writing, create if it does not exist, truncate if it does exist.
	outFile, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return fmt.Errorf("failed to open file for writing %w: %s", err, url)
	}
	defer outFile.Close()

	fmt.Println(GetCurrentTime(true))

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	fmt.Print("\rsending request, awaiting response...")
	resp, err := client.Get(url)
	if err != nil {
		xerr.WriteError(err, 2, false)
	}

	fmt.Printf("\rsending request, awaiting response... %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))

	if err != nil {
		return fmt.Errorf("failed to download %s: %w", url, err)
	}

	defer resp.Body.Close()

	n, err := io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save URL %s to file %w", url, err)
	}
	r, _ := client.Head(url)
	contentSizeMB := RoundOfSizeOfData(r.ContentLength)

	fmt.Printf("\rcontent size: %d [~%s]\n", n, contentSizeMB)

	fmt.Println(GetCurrentTime(false))
	return nil
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

// GetCurrentTime prints a textual representation of the current time, it takes isStart
// if isStart is true it indicates the download process has started, if false means end of the current download
func GetCurrentTime(isStart bool) string {
	// Get the current time
	currentTime := time.Now()

	// Format time to print up to seconds
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	if isStart {
		return fmt.Sprintf("start at %s", formattedTime)
	}
	return fmt.Sprintf("finished at %s", formattedTime)
}

// RoundOfSizeOfData  converts dataInBytes (size of file downloaded) in bytes to the nearest size
func RoundOfSizeOfData(dataInBytes int64) string {
	var size float64
	var unit string
	if dataInBytes >= GB {
		size = float64(dataInBytes) / GB
		unit = "GB"
	} else if dataInBytes >= KB {
		size = float64(dataInBytes) / MB
		unit = "MB"
	} else {
		size = float64(dataInBytes)
		unit = "KB"
	}
	return fmt.Sprintf("%.2f%s", size, unit)
}
