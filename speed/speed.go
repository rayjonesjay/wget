// speed package contains functionalities for adjusting the download speed.
package speed

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// RateLimiter controls the rate of data processing
type RateLimiter struct {
	rate      int       // bytes per second
	lastCheck time.Time // last time the limiter was checked
	burstSize int       // maximum burst size in bytes
}

// rateLimitedReader applies rate limiting to the reader
type rateLimitedReader struct {
	reader      io.Reader
	limiter     *RateLimiter
	burstBuffer int // track bytes read for burst management
}

// Read method applies rate limiting to each read operation
func (r *rateLimitedReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if err != nil {
		return n, err
	}

	// apply rate limiting
	r.burstBuffer += n
	if r.burstBuffer >= r.limiter.burstSize {
		r.limiter.Wait(r.burstBuffer)
		r.burstBuffer = 0
	}

	return n, nil
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(rate, burstSize int) RateLimiter {
	return RateLimiter{
		rate:      rate,
		lastCheck: time.Now(),
		burstSize: burstSize,
	}
}

// Wait applies the rate limit by introducing a delay based on the amount of data to be processed
func (rl *RateLimiter) Wait(bytesProcessed int) {
	elapsedTime := time.Since(rl.lastCheck)

	expectedElapsed := time.Duration(bytesProcessed) * time.Second / time.Duration(rl.rate)

	if expectedElapsed > elapsedTime {
		time.Sleep(expectedElapsed - elapsedTime)
	}

	rl.lastCheck = time.Now()
}

// DownloadFileWithRateLimit downloads a file while limiting the download speed
func DownloadFileWithRateLimit(url, filename string, limiter *RateLimiter) error {
	// create the file
	out, err := os.Create(filename)
	if err != nil {
		return err // if an error occurs while creating the file
	}

	defer out.Close()

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create a reader that applies the rate limiter
	limitedReader := &rateLimitedReader{
		reader:  resp.Body,
		limiter: limiter,
	}

	// Copy the data from the response to the file
	_, err = io.Copy(out, limitedReader)
	if err != nil {
		return err
	}

	fmt.Println("download complete")
	return nil
}
