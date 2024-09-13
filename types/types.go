// Package types contains user defined types and constants
package types

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

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

// Arg represents the commandline arguments passed through the command line by the user
// example: $ go run . -O=file.txt -B https://www.example.com
// -0 will be a field in Arg that specifies the output file to save the resource
// -B will send the download process to background mode
// https://www.example.com is the link to where the resource resides
type Arg struct {
	OutputFile     string   // identified by the -O flag
	BackgroundMode bool     // identified by the -B flag
	Links          []string // identified by the regexp pattern (http|https)://\w+ ,specifies path to resources on the web
	SavePath       string   // identified by the -P flag, specifies the location where to save the resource
	InputFile      string   // identified by the -i flag, specifies a file contains url(s)
	Rejects        []string // identified by the -R or --reject flag contains a list of resources to reject
	Mirror         bool     // identified by the --mirror flag, indicates whether to download an entire website or not
	RateLimit      string   // identified by the --rate-limit flag, specifies the download speed when fetching a resource
	RateLimitValue int64    // if RateLimit is specified, RateLimitValue will be
	IsHelp         bool     // identified by the --help flag, if pared it will print our program manual
	ConvertLinks   bool     // identified by the --convert-links
	Exclude        []string // identified by the --exclude or -X, takes a comma separated list of paths(directory) to avoid when fetching a resource
}

// Download handles normal downloads based on the provided URLs and other flags in the Arg struct
func (a *Arg) Download() error {

	var wg sync.WaitGroup

	// if a download is successful, the link will be stored here.
	successfulDownloads := make([]string, 0)

	// channel to receive URLs that are successfully downloaded.
	successChannel := make(chan string)

	// channel to receive errors
	errorChannel := make(chan error)

	// start a go routine to handle successful downloads
	go func() {
		for url := range successChannel {
			successfulDownloads = append(successfulDownloads, url)
		}
	}()

	// start a go routine to handle errors
	go func() {
		for url := range errorChannel {
			_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", url)
		}
	}()

	for _, url := range a.Links {
		wg.Add(1)

		go func(url string) {

			defer wg.Done()

			// determine the output path and filename
			outputFilePath := a.determineOutputPath(url)

			// open the output file for writing, create if it does not exist, truncate if it does exist.
			outFile, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

			if err != nil {
				errorChannel <- fmt.Errorf("failed to open file for writing %w: %s", err, url)
				return
			}
			defer outFile.Close()

			// get the data
			resp, err := http.Get(url)
			fmt.Print("HTTP request sent, awaiting response...")
			if err != nil {
				errorChannel <- fmt.Errorf("failed to download %s: %w", url, err)
				return
			}

			if resp.StatusCode == 200 {
				a.LogDownloadInfo(resp)
			} // else log errors here

			defer resp.Body.Close()

			// copy the data to the outFile
			n, err := io.Copy(outFile, resp.Body)
			if err != nil {
				errorChannel <- fmt.Errorf("failed to save URL %s to file %w", url, err)
				return
			}
			fmt.Printf("downloaded %d bytes from %s\n", n, url)

			// successful download
			successChannel <- url
		}(url)
	}

	// wait for all downloading to complete
	wg.Wait()
	close(successChannel)
	close(errorChannel)

	//fmt.Println(">>>")
	for _, successURL := range successfulDownloads {
		successfulDownloads = append(successfulDownloads, successURL)
	}
	printUrls(successfulDownloads)
	return nil
}

func (a *Arg) LogDownloadInfo(response *http.Response) {
	code := response.StatusCode
	msg := fmt.Sprintf("\rHTTP request sent, awaiting response... %d OK", code)
	fmt.Println(msg)
}

// determineOutputPath determines the full path for the output file
func (a *Arg) determineOutputPath(url string) string {
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
	fmt.Printf("Download finished:	[%s]\n", res)
}

// IsEmpty function checks whether an iterable is empty, an iterable is a string,array or slice.
// it returns true if the `data` which is expected to be an iterable, is empty else return false
func (a *Arg) IsEmpty(data interface{}) bool {

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
