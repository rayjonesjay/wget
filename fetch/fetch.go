// Package fetch contains the URL utility, to fetch a resource from a given URL
// with configurable options, including rate limiting the download.
package fetch

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"sync"
	"wget/fileio"
	"wget/httpx"
	"wget/limitedio"
	"wget/syscheck"
)

const KiB = 1024

var client = &http.Client{}

// headers our http client will send by default. Adapted from Chrome,
// as some web servers will deny requests without a valid user agent
var headers = map[string]string{
	"connection":                "keep-alive",
	"sec-ch-ua":                 "\"Chromium\";v=\"128\", \"Not;A=Brand\";v=\"24\", \"Google Chrome\";v=\"128\"",
	"sec-ch-ua-mobile":          "?0",
	"sec-ch-ua-platform":        "\"Linux\"",
	"dnt":                       "1",
	"upgrade-insecure-requests": "1",
	"user-agent": "Mozilla/5.0 (X11; Linux x86_64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
	"accept": "text/html,application/xhtml+xml,application/xml;" +
		"q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	"sec-fetch-site":  "none",
	"sec-fetch-mode":  "navigate",
	"sec-fetch-user":  "?1",
	"sec-fetch-dest":  "document",
	"accept-encoding": "identity",
	"accept-language": "en-US,en;q=0.9,la;q=0.8",
}

// FileInfo contains information from the server about the downloaded file from URL
type FileInfo struct {
	Name       string
	Headers    http.Header
	StatusCode int
}

// Config contains configuration options for URL
type Config struct {
	// GetFile will be called to return a valid file, with write access, to hold the resource from the
	// given URL. The function is provided the headers received from the request to
	// the target url
	GetFile func(url string, header http.Header) (*os.File, error)
	// Limit the download speed to a maximum of Limit bytes/second. A Limit <= 0 infers no rate limiting
	Limit int32
	// ProgressListener will be called every time some buffered read occurs,
	// depending on the underlying buffer size or the Limit. It reports how much of
	// the resource has been downloaded, and how much is the expected total.
	//Do not do long-running tasks directly in this callback's goroutine, create a separate goroutine to do that.
	//
	//Note, that, some servers may not send the content-length header, thus,
	//the reported `total` length will be set to -1
	ProgressListener func(downloaded, total int64)
	// RateListener will be called, every second, from when the download starts, to
	// when the download ends, to report the current bandwidth being utilized by the
	// download, i.e., the `rate` of the download in bytes/second
	RateListener func(rate int32)
	// The request body
	Body io.Reader
	// Method holds the HTTP request method to use for the request, will default to GET if undefined
	Method string
	// AllowedStatusCodes keeps a list of all the status codes that are allowed for the given request.
	//Any other status code will be considered an error
	AllowedStatusCodes []int
}

type DownloadStatus struct {
	StatusCode    int
	Status        string
	Start         string
	End           string
	ContentLength string
}

var m sync.Mutex
var Status = &DownloadStatus{} // am having problems accessing this struct concurrently

// URL downloads the file from the given url, and saves it to the given file,
// respecting the given speed limit; i.e., the download speed never exceeds `limit` bytes/second
func URL(url string, config Config) (info FileInfo, err error) {

	m.Lock()
	Status.Start = syscheck.GetCurrentTime(true)

	m.Unlock()
	{ // sanity checks on the configuration
		if config.GetFile == nil {
			err = errors.New("bad config: function `GetFile` is required")
			return
		}

		if config.RateListener == nil {
			config.RateListener = func(int32) {}
		}

		if config.ProgressListener == nil {
			config.ProgressListener = func(int64, int64) {}
		}

		if config.Method == "" {
			config.Method = "GET"
		}
	}

	req, err := http.NewRequest(config.Method, url, config.Body)
	if err != nil {
		return info, fmt.Errorf("failed to create new request: %v", err)
	}

	// Set the default client headers, including user agent
	setClientHeaders(&req.Header)

	m.Lock()
	// Send the request
	Status.Status = "\rsending request, awaiting response..."
	m.Unlock()

	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to download file: %v", err)
		return
	}
	defer fileio.Close(resp.Body)

	m.Lock()
	Status.StatusCode = resp.StatusCode
	Status.Status = fmt.Sprintf("\rsending request, awaiting response... %d %s\n", Status.StatusCode, Status.Status)
	m.Unlock()
	if config.AllowedStatusCodes != nil && !slices.Contains(config.AllowedStatusCodes, resp.StatusCode) {
		err = fmt.Errorf("wrong status code: %v", resp.StatusCode)
		return
	}

	info.Headers = resp.Header

	m.Lock()
	info.StatusCode = resp.StatusCode
	m.Unlock()
	// m.Lock()
	contentLength := httpx.ExtractContentLength(resp.Header)
	Status.ContentLength = fmt.Sprintf("\rcontent size: %d [~%s]\n", contentLength, httpx.RoundOfSizeOfData(contentLength))
	// m.Unlock()
	// Create the output file
	file, err := config.GetFile(url, resp.Header)
	if err != nil {
		err = fmt.Errorf("failed to get writable file: %v", err)
		return
	}
	defer fileio.Close(file)
	info.Name = file.Name()

	// Create a buffer to store the downloaded bytes
	// Many clients use a default buffer size of 8KiB, we follow that standard
	buffer := make([]byte, 8*KiB)

	// Use a speed governed reader to limit reads from the response body
	body := limitedio.NewSGReader(config.Limit, &resp.Body)
	defer fileio.Close(body)
	body.SetRateListener(config.RateListener)

	// keeps track of how many bytes have been downloaded
	downloadedBytes := int64(0)

	// ReadAll bytes from the speed governed response body in chunks of 8KiBs.
	//See io.ReadAll for more details
	for {
		var n int
		// Read a chunk of bytes from the response body
		n, err = body.Read(buffer)
		if err != nil {
			if err == io.EOF && n == 0 {
				// Reached the end of the file, shouldn't be reported as an error
				err = nil
				break
			} else if err != io.EOF {
				err = fmt.Errorf("failed to read response body: %v", err)
				return
			}
			// read some n bytes, before reaching the end of the file,
			//shouldn't be reported as an error; will then proceed to save them to file
			err = nil
		}
		downloadedBytes += int64(n)

		// Write the chunk of bytes to the output file
		_, err = file.Write(buffer[:n])
		if err != nil {
			err = fmt.Errorf("failed to write to download file: %v", err)
			return
		}

		config.ProgressListener(downloadedBytes, contentLength)
	}
	m.Lock()
	Status.End = syscheck.GetCurrentTime(false)
	m.Unlock()
	return info, nil
}

// setClientHeaders updates the pointed request http.Header with the default client headers
func setClientHeaders(reqHeaders *http.Header) {
	for key, value := range headers {
		reqHeaders.Set(key, value)
	}
}
