// Package types contains user defined types and constants
package types

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
	"wget/errorss"
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
	for _, url := range a.Links {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := a.GetResource(url)
			if err != nil {
				errorss.WriteError(err, 2, false)
			}
		}()
	}

	wg.Wait()
	return nil
}

func (a *Arg) GetResource(url string) (err error) {
	outputFilePath := a.determineOutputPath(url)
	// open the output file for writing, create if it does not exist, truncate if it does exist.
	outFile, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

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
		errorss.WriteError(err, 2, false)
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
	fmt.Printf("Downloaded	[%s]\n", res)
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

func (a *Arg) MirrorWeb() error {

	for _, link := range a.Links {
		parsedUrl, err := url.Parse(link)
		if err != nil {
			return err
		}

		// by default the downloaded data will be saved to the name of domain if not specified
		domain := parsedUrl.Host
		directoryToSaveData := filepath.Join(a.SavePath, domain)

		err = os.MkdirAll(directoryToSaveData, 0755)
		if err != nil {
			return err
		}

		// Download and parse the HTML/CSS
		//err = a.downloadAndParseHTML(link, directoryToSaveData)
		//if err != nil {
		//	return err
		//}
	}
	return nil
}

//func (a *Arg) downloadAndParseHTML(link, saveDir) error {
//	{
//	}
//}

// rejectFileBasedOnExtension method checks if a file should be downloaded based on
// its file extension
func (a *Arg) rejectFileBasedOnExtension(url string) bool {
	for _, ext := range a.Rejects {
		if strings.HasSuffix(url, ext) {
			return true
		}
	}
	return false
}

// shouldExcludeDir method checks if an url belongs to an excluded directory
func (a *Arg) shouldExcludeDir(url string) bool {
	for _, directory := range a.Exclude {
		if strings.Contains(url, directory) {
			return true
		}
	}
	return false
}

func (a *Arg) Run() {
	if a.Mirror {
		// run in mirror mode
		_ = a.MirrorWeb()
		//if err != nil {
		//	return
		//}
	} else {
		// regular download
		_ = a.Download()
		//if err != nil {
		//	return
		//}
	}
}

//func (a *Arg) ConvertLinksForOfflineView(htmlContent, saveDir string) string {
//	// replace href/src urls with local file paths
//	// example: href="http://example.com/style.css" -> href="./style.css"
//	return htmlCOncent
//}

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

// IsValidURL checks if the given string is a valid URL
func IsValidURL(urlStr string) (bool, error) {
	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false, fmt.Errorf("invalid URL: %v", err)
	}

	// check if scheme component of the URL is empty
	if !parsedURL.IsAbs() {
		return false, errorss.ErrNotAbsolute
	}

	// check if the scheme is neither http nor https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false, errorss.ErrWrongScheme
	}

	// if host is empty
	if parsedURL.Host == "" {
		return false, errorss.ErrEmptyHostName
	}

	// ensure host does not start with . or -
	if strings.HasPrefix(parsedURL.Host, ".") || strings.HasPrefix(parsedURL.Host, "-") {
		return false, fmt.Errorf("wrong host format %q", parsedURL.Host)
	}

	// Check that the host contains at least one dot (valid domain format) or is localhost
	if !strings.Contains(parsedURL.Host, ".") && parsedURL.Host != "localhost" {
		return false, errorss.ErrInvalidDomainFormat
	}

	return true, nil
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
