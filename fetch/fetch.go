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
	"strings"
	"sync"
	"time"

	"wget/fileio"
	"wget/globals"
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

// AdvancedProgressListener registers some callbacks that will be called when
// certain download events occur during the download
type AdvancedProgressListener struct {
	// OnStart will be called with the time just before the download starts
	OnStart func(time time.Time)
	// OnStatus will be called with the status code and message received
	OnStatus func(status string, code int)
	// OnContentLength will be called with the expected size of the download. This
	// may be -1 if the content length is unknown, e.g, when the server is streaming
	// the response
	OnContentLength func(length int64)
	// OnGetFile will be called with the filename the contents will be downloaded to
	OnGetFile func(filename string)
	// OnProgress will be called with stats on the size of the content that has been
	// downloaded against the total content length, alongside the rate of the
	// download in bytes/second
	OnProgress func(downloaded, total int64, rate int32)
	// OnDownloadFinished will be called with the time the whole content body was downloaded
	OnDownloadFinished func(url string, time time.Time)
}

// init initializes the receiver progress listener, in place, with the default no-op status listeners
func (from *AdvancedProgressListener) init() {
	l := from
	if from.OnStart == nil {
		l.OnStart = func(time time.Time) {}
	}
	if from.OnStatus == nil {
		l.OnStatus = func(status string, code int) {}
	}
	if from.OnContentLength == nil {
		l.OnContentLength = func(length int64) {}
	}
	if from.OnGetFile == nil {
		l.OnGetFile = func(filename string) {}
	}
	if from.OnProgress == nil {
		l.OnProgress = func(downloaded, total int64, rate int32) {}
	}
	if from.OnDownloadFinished == nil {
		l.OnDownloadFinished = func(url string, time time.Time) {}
	}
}

// Config contains configuration options for URL
type Config struct {
	// GetFile will be called to return a valid file, with write access, to hold the resource from the
	// given URL. The function is provided the headers received from the request to
	// the target url
	GetFile func(url string, header http.Header) (*os.File, error)
	// ShouldDownload will be called to validate whether the file from the given url
	// should be downloaded, based on the given headers as retrieved from the server
	ShouldDownload func(url string, header http.Header) bool
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
	// Any other status code will be considered an error
	AllowedStatusCodes []int
	AdvancedProgressListener
}

// URL downloads the file from the given url, and saves it to the given file,
// respecting the given speed limit; i.e., the download speed never exceeds `limit` bytes/second
func URL(url string, config Config) (info FileInfo, err error) {
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

		config.AdvancedProgressListener.init()
	}

	defer func() {
		if err != nil {
			config.AdvancedProgressListener.OnDownloadFinished("", time.Now())
		} else {
			config.AdvancedProgressListener.OnDownloadFinished(url, time.Now())
		}
	}()

	req, err := http.NewRequest(config.Method, url, config.Body)
	if err != nil {
		return info, fmt.Errorf("failed to create new request: %v", err)
	}

	// Set the default client headers, including user agent
	setClientHeaders(&req.Header)

	// Send the request
	config.AdvancedProgressListener.OnStart(time.Now())
	config.AdvancedProgressListener.OnStatus("", -1)
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to download file: %v", err)
		return
	}
	defer fileio.Close(resp.Body)

	config.AdvancedProgressListener.OnStatus(resp.Status, resp.StatusCode)
	if config.ShouldDownload != nil && !config.ShouldDownload(url, resp.Header) {
		err = fmt.Errorf("skipping download of url: %q", url)
		return
	}

	if config.AllowedStatusCodes != nil && !slices.Contains(config.AllowedStatusCodes, resp.StatusCode) {
		err = fmt.Errorf("bad status code: %v", resp.Status)
		return
	}

	info.Headers = resp.Header
	info.StatusCode = resp.StatusCode

	contentLength := httpx.ExtractContentLength(resp.Header)
	config.AdvancedProgressListener.OnContentLength(contentLength)

	// Create the output file
	file, err := config.GetFile(url, resp.Header)
	if err != nil {
		err = fmt.Errorf("failed to get writable file: %v", err)
		return
	}
	defer fileio.Close(file)
	info.Name = file.Name()
	config.AdvancedProgressListener.OnGetFile(file.Name())

	// Create a buffer to store the downloaded bytes
	// Many clients use a default buffer size of 8KiB, we follow that standard
	buffer := make([]byte, 8*KiB)

	// Use a speed governed reader to limit reads from the response body
	body := limitedio.NewSGReader(config.Limit, &resp.Body)
	defer fileio.Close(body)
	body.SetRateListener(
		func(rate int32) {
			config.AdvancedProgressListener.OnProgress(-1, -1, rate)
			config.RateListener(rate)
		},
	)

	// keeps track of how many bytes have been downloaded
	downloadedBytes := int64(0)

	// ReadAll bytes from the speed governed response body in chunks of 8KiBs.
	// See io.ReadAll for more details
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
			// shouldn't be reported as an error; will then proceed to save them to file
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
		config.AdvancedProgressListener.OnProgress(downloadedBytes, contentLength, -1)
	}

	return info, nil
}

// DownloadStatus holds printable download status of a Get request at any one instance during the download
type DownloadStatus struct {
	StatusCode    int
	Status        string
	Start         string
	SavePath      string
	ContentLength string
	Progress      string
	Finished      string
	End           string

	// StartTime records the time the download started
	StartTime time.Time

	Downloaded, Total int64
	Rate              int32
	// OnUpdate will be called whenever the download status (of this struct) changes.
	// A reference to this struct is provided for convenience
	OnUpdate func(status *DownloadStatus, hint int)
	// m mutex locker. Some status update operations, such as Progress, may be
	// changed by different threads; this locker will manage race conditions that may
	// arise
	m *sync.Mutex
	// n is the number os statistical observations, of the Rate, to be used in
	// calculating the moving average of the Rate
	n int
	// avgRate is the average Rate of the download
	avgRate int32
}

// GetField returns the field of the Task based on the provided index.
func (s *DownloadStatus) GetField(index int) string {
	switch index {
	case 0:
		return s.Start
	case 1:
		return s.Status
	case 2:
		return s.ContentLength
	case 3:
		return s.SavePath
	case 4:
		return s.Progress
	case 5:
		return s.Finished
	case 6:
		return s.End
	default:
		panic("Index out of range")
	}
}

// ProgressListener builds a compliant advanced download progress listener, that
// reports its download status to the receiver DownloadStatus
func (s *DownloadStatus) ProgressListener() *AdvancedProgressListener {
	l := &AdvancedProgressListener{}
	if s.m == nil {
		s.m = new(sync.Mutex)
	}

	if s.OnUpdate == nil {
		panic("You must initialize [OnUpdate] before this call")
	}

	m := s.m
	l.OnStart = func(t time.Time) {
		s.Start = fmt.Sprintf("start at %s", format(t))
		s.StartTime = t
		s.OnUpdate(s, 0)
	}

	l.OnStatus = func(status string, _ int) {
		if status != "" {
			status = fmt.Sprintf("status %s", status)
		}
		s.Status = fmt.Sprintf("\rsending request, awaiting response... %s\n", status)
		s.OnUpdate(s, 1)
	}

	l.OnContentLength = func(length int64) {
		if length < 0 {
			s.ContentLength = fmt.Sprintf("\rcontent size: unspecified [~%s]\n", globals.FormatSize(length))
		} else {
			s.ContentLength = fmt.Sprintf("\rcontent size: %d [~%s]\n", length, globals.RoundBytes(length))
		}
		s.OnUpdate(s, 2)
	}

	l.OnGetFile = func(filename string) {
		filename = fileio.AliasUserDir(filename)
		filename = fileio.PrependDotIfInCurrentDir(filename)
		s.SavePath = fmt.Sprintf("saving file to: %s", filename)
		s.OnUpdate(s, 3)
	}

	l.OnProgress = func(downloaded, total int64, rate int32) {
		m.Lock()
		defer m.Unlock()

		if rate >= 0 {
			s.Rate = rate
		} else {
			s.Downloaded = downloaded
			s.Total = total
		}

		// update the moving average of the download rate
		if rate > 0 {
			s.avgRate = int32(updateAverage(float64(s.avgRate), float64(rate), s.n))
			s.n++
		}

		s.Progress = onProgress(s.Downloaded, s.Total, s.Rate, "")
		s.OnUpdate(s, 4)
	}

	l.OnDownloadFinished = func(url string, t time.Time) {
		if url == "" {
			// this download failed
			s.Finished = fmt.Sprintf("\u001B[0;31m\nfailed at %s\u001B[0m", format(t))
			s.OnUpdate(s, 5)
			return
		}

		s.Finished = fmt.Sprintf("\nDownloaded [%s]\u001B[0;32m\nfinished at %s\u001B[0m", url, format(t))
		s.OnUpdate(s, 5)
		duration := t.Sub(s.StartTime).Truncate(time.Second)
		if s.Total == -1 {
			s.Total = s.Downloaded
			l.OnContentLength(s.Downloaded)
			s.OnUpdate(s, 2)
		}
		if s.avgRate == 0 {
			// This file was downloaded in less than a second, the average speed is relative
			// to the size of the file
			s.avgRate = int32(s.Downloaded)
		}
		s.Progress = onProgress(s.Downloaded, s.Total, s.avgRate, duration.String())
		s.OnUpdate(s, 4)
	}

	return l
}

func onProgress(downloaded, total int64, rate int32, eta string) string {
	width := syscheck.GetTerminalWidth()
	barLength := width / 3
	percentageString := "--.-%"

	bar := ""
	progress := float64(0)
	percentage := float64(0)
	if total <= 0 {
		bar = fmt.Sprintf("[<=>%s]", strings.Repeat(" ", barLength-3))
	} else {
		progress = float64(downloaded) / float64(total)
		percentage = float64(downloaded*100) / float64(total)

		filled := int(progress * float64(barLength))
		notFilled := barLength - filled
		repeat := func(s string, count int) string {
			if count == 0 && barLength > 1 {
				return ">"
			} else if count > 0 {
				return strings.Repeat(s, count-1) + ">"
			}
			return strings.Repeat(s, count)
		}

		bar = fmt.Sprintf("[%s%s]", repeat("=", filled), strings.Repeat(" ", notFilled))
		percentageString = fmt.Sprintf("%.2f%%", percentage)
	}

	if eta == "" {
		// get the estimated time to finish the download
		eta = CalculateETA(downloaded, total, rate)
	} else {
		// ignore the download duration, instead default to eta as in the question
		// in this case, the download is finished, so the eta is simply zero seconds
		eta = "0s"
	}

	return fmt.Sprintf(
		" %s / %s %s %s %s %s",
		globals.FormatSize(downloaded), globals.FormatSize(total), bar, percentageString,
		globals.FormatSize(int64(rate))+"/s", eta,
	)
}

// setClientHeaders updates the pointed request http.Header with the default client headers
func setClientHeaders(reqHeaders *http.Header) {
	for key, value := range headers {
		reqHeaders.Set(key, value)
	}
}

func format(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// CalculateETA estimates the time left for a download.
func CalculateETA(downloaded, total int64, rate int32) string {
	if downloaded == 0 || total < 0 || rate == 0 { // to avoid division by zero
		return "🕛"
	}
	remaining := float64(total-downloaded) / float64(rate)
	eta := time.Duration(remaining) * time.Second
	return eta.String()
}

// updateAverage calculates the new average given the current average, the new observation,
// and the number of observations.
func updateAverage(currentAvg float64, newObservation float64, n int) float64 {
	return ((currentAvg * float64(n)) + newObservation) / float64(n+1)
}
