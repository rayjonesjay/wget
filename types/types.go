// types contains user defined types and constants
package types

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	successfulDownloads := []string{}
	//failedDownloads := []string{}
	for _, url := range a.Links {
		// determine the output path and filename
		outputFilePath := a.determineOutputPath(url)

		// open the output file for writing
		outFile, err := os.Create(outputFilePath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer outFile.Close()

		// get the data
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to download %s: %w", url, err)
		}
		defer resp.Body.Close()

		// copy the data to the outFile
		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to save URL %s to file %w", url, err)
		}
		successfulDownloads = append(successfulDownloads, url)
	}
	// TODO Print all the urls after successful download
	if len(successfulDownloads) > 0 {
		printUrls(successfulDownloads)
	}
	return nil
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

func (a *Arg) downloadFile(url, outputFilePath string) error {

	// open the output file for writing
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

}

const (
	// Size of Data in KilobBytes
	KB = 1000 * 1

	// Size of Data in KibiBytes, same as 2^10
	KiB = 1 << (10 * 1)

	// Size of Data in MegaBytes
	MB = 1000 * KB

	// Size of Data in MebiBytes, same as 2^20
	MiB = 1 << 20

	// Size of Data in GigaBytes
	GB = 1000 * MB

	// Size of Data in GibiBytes, same as 2^30
	GiB = 1 << 30

	// Size of Data in TeraBytes
	TB = 1000 * GB

	// Size of Data in TebiBytes, same as 2^40
	TiB = 1 << 40
)

var (
	// Display usage message to user on how to use our program
	UsageMessage = `wget: missing URL
Usage: $ ./wget [OPTION] [URL]

Try './wget --help' for more options.`

	Manual = `wget 1.0.0, a non-interactive retriever.
Usage: ./wget [OPTION] [URL]

Supported options as of version 1.0.0:

FLAG:				 USAGE:						EXPLANATION:

--help				 ******						print this manual
-B 	  		        -B URL					        allow download in the background after startup
-v 			        ******					        display the current version of wget and exit
-O   				-O=FILE	 URL					log messages to FILE
--mirror  			--mirror URL					mirror a website
-P			        -P=PATH  URL				        specify path to save downloaded resource
--rate-limit=NS									set speed limit, N is a number, S is either (k or M)
-i  				-i=FILE						allow download of resource by reading links stored in a file


bug reports, questions, issues to https://github.com/rayjonesjay
`
	Version = "1.0.0"

	DefaultDownloadDirectory = "$HOME/Downloads"

	// default port used by http if not specified
	DefaultHTTPPort = "80"

	// default port used by https if not specified
	DefaultHTTPSPort = "443"
)
