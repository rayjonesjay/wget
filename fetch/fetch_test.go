package fetch

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	rand2 "math/rand/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var randomFileHash = "Random File Hash"
var randomData []byte

func init() {
	randomData = make([]byte, 10*1024*1024)
	for i := range randomData {
		randomData[i] = byte(rune(rand2.IntN(126-32) + 32))
	}
	// Calculate the SHA256 hash
	randomFileHash = fmt.Sprintf("%x", sha256.Sum256(randomData))
}

func TestFetch(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Set headers
				h := w.Header()
				h.Set("Content-Type", "application/octet-stream")
				h.Set("Content-Disposition", "attachment; filename=\"random.bin\"")
				h.Set("Content-Length", fmt.Sprintf("%d", len(randomData)))
				_, _ = w.Write(randomData)
			},
		),
	)

	ShortRequestServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if _, err := w.Write([]byte("H")); err != nil {
					http.Error(w, "Failed to write response", http.StatusInternalServerError)
					return
				}
			},
		),
	)

	EmptyRequestServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if _, err := w.Write([]byte("")); err != nil {
					http.Error(w, "Failed to write response", http.StatusInternalServerError)
					return
				}
			},
		),
	)

	BadContentLengthServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Set a header that indicates an incompatible content length.
				w.Header().Set("Content-Length", "100")

				// Write a smaller amount of data than what the Content-Length header claims
				_, err := fmt.Fprint(w, "Hello")
				if err != nil {
					http.Error(w, "Failed to write response", http.StatusInternalServerError)
					return
				}
			},
		),
	)

	GetFile := func(url string, headers http.Header) (*os.File, error) {
		file, err := createTempReadWriteFile()
		if err != nil {
			return nil, err
		}
		return file, nil
	}

	ErroneousGetFile := func(url string, headers http.Header) (*os.File, error) {
		return nil, errors.New("ErroneousGetFile: won't create a file")
	}

	GetFile2 := func(url string, header http.Header) (*os.File, error) {
		file, err := createTempReadWriteFile()
		if err != nil {
			return nil, err
		}
		return file, nil

	}

	ClosingGetFile := func(url string, header http.Header) (*os.File, error) {
		file, err := createTempReadWriteFile()
		if err != nil {
			return nil, err
		}
		_, _ = file.WriteString("No Op")
		err = file.Close()
		if err != nil {
			return nil, err
		}
		return file, nil

	}

	var ClosedBody io.ReadCloser
	{
		file, err := createTempReadWriteFile()
		if err != nil {
			t.Fatalf("Error creating temp file: %v", err)
		}
		_, _ = file.WriteString("No Op")
		err = file.Close()
		if err != nil {
			t.Fatalf("Error closing temp file: %v", err)
		}

		ClosedBody = file
	}

	type args struct {
		url    string
		config Config
	}
	tests := []struct {
		name               string
		args               args
		wantInfo           FileInfo
		wantErr            bool
		wantPanic          bool
		compareInfo        func(a, b FileInfo) bool
		dontCheckHashMatch bool
	}{
		{
			name: "No GetFile function",
			args: args{
				url: server.URL,
				config: Config{
					GetFile:          nil,
					Limit:            0,
					ProgressListener: nil,
					RateListener:     nil,
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr: true,
		},

		{
			name: "Error-prone GetFile function",
			args: args{
				url: server.URL,
				config: Config{
					GetFile:          ErroneousGetFile,
					Limit:            0,
					ProgressListener: nil,
					RateListener:     nil,
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr: true,
		},

		{
			name: "GetFile function returning closed file",
			args: args{
				url: server.URL,
				config: Config{
					GetFile:          ClosingGetFile,
					Limit:            0,
					ProgressListener: nil,
					RateListener:     nil,
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr: true,
		},

		{
			name: "No Listener functions",
			args: args{
				url: server.URL,
				config: Config{
					GetFile:          GetFile,
					Limit:            0,
					ProgressListener: nil,
					RateListener:     nil,
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr: false,
		},

		{
			name: "No Listener functions negative rate limit",
			args: args{
				url: server.URL,
				config: Config{
					GetFile:          GetFile,
					Limit:            -10,
					ProgressListener: nil,
					RateListener:     nil,
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr: false,
		},

		{
			name: "Custom Listener functions",
			args: args{
				url: server.URL,
				config: Config{
					GetFile:          GetFile2,
					Limit:            0,
					ProgressListener: func(downloaded, total int64) {},
					RateListener:     func(rate int32) {},
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr: false,
		},

		{
			name: "Unsupported scheme:Custom Listener functions",
			args: args{
				url: "http",
				config: Config{
					GetFile:          GetFile2,
					Limit:            0,
					ProgressListener: func(downloaded, total int64) {},
					RateListener:     func(rate int32) {},
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr: true,
		},

		{
			name: "Bad-URL:Custom Listener functions",
			args: args{
				url: "https://",
				config: Config{
					GetFile:          GetFile2,
					Limit:            0,
					ProgressListener: func(downloaded, total int64) {},
					RateListener:     func(rate int32) {},
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr: true,
		},

		{
			name: "POST Request with Body: Custom Listener functions",
			args: args{
				url: server.URL,
				config: Config{
					GetFile:          GetFile2,
					Limit:            0,
					ProgressListener: func(downloaded, total int64) {},
					RateListener:     func(rate int32) {},
					Body:             strings.NewReader("Hello"),
					Method:           "POST",
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr: false,
		},

		{
			name: "Invalid HTTP request method, Request with Body: Custom Listener functions",
			args: args{
				url: server.URL,
				config: Config{
					GetFile:          GetFile2,
					Limit:            0,
					ProgressListener: func(downloaded, total int64) {},
					RateListener:     func(rate int32) {},
					Body:             strings.NewReader("Hello"),
					Method:           "SUPERFLUOUS!!",
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr: false,
		},

		{
			name: "Empty URL: Invalid HTTP request method Request with Body: Custom Listener functions",
			args: args{
				url: "",
				config: Config{
					GetFile:          GetFile,
					Limit:            0,
					ProgressListener: func(downloaded, total int64) {},
					RateListener:     func(rate int32) {},
					Body:             strings.NewReader("Hello"),
					Method:           "SUPERFLUOUS!!",
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr: true,
		},

		{
			name: "POST Request with closed Body: Custom Listener functions",
			args: args{
				url: server.URL,
				config: Config{
					GetFile:          GetFile2,
					Limit:            0,
					ProgressListener: func(downloaded, total int64) {},
					RateListener:     func(rate int32) {},
					Body:             ClosedBody,
					Method:           "POST",
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr: true,
		},

		{
			name: "GET Request with Body: Custom Listener functions",
			args: args{
				url: server.URL,
				config: Config{
					GetFile:          GetFile2,
					Limit:            0,
					ProgressListener: func(downloaded, total int64) {},
					RateListener:     func(rate int32) {},
					Body:             strings.NewReader("Hello"),
					Method:           "GET",
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr: false,
		},

		{
			name: "GET Request to BadContentLengthServer with Custom Listener functions",
			args: args{
				url: BadContentLengthServer.URL,
				config: Config{
					GetFile:          GetFile,
					Limit:            0,
					ProgressListener: func(downloaded, total int64) {},
					RateListener:     func(rate int32) {},
					Method:           "GET",
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr:            true,
			dontCheckHashMatch: true,
		},

		{
			name: "GET Request with limiter on Short Request Server Custom Listener functions",
			args: args{
				url: ShortRequestServer.URL,
				config: Config{
					GetFile:          GetFile,
					Limit:            10,
					ProgressListener: func(downloaded, total int64) {},
					RateListener:     func(rate int32) {},
					Body:             nil,
					Method:           "GET",
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr:            false,
			dontCheckHashMatch: true,
		},

		{
			name: "GET Request with limiter on Empty Response Server Custom Listener functions",
			args: args{
				url: EmptyRequestServer.URL,
				config: Config{
					GetFile:          GetFile,
					Limit:            0,
					ProgressListener: func(downloaded, total int64) {},
					RateListener:     func(rate int32) {},
					Body:             nil,
					Method:           "GET",
				},
			},
			wantInfo: FileInfo{
				Name:    "",
				Headers: nil,
			},
			wantErr:            false,
			dontCheckHashMatch: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotInfo, err := URL(tt.args.url, tt.args.config)
				if (err != nil) != tt.wantErr {
					t.Errorf("URL() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if err != nil {
					return
				}

				// check that the downloaded file has the same hash as the original file from the test server
				fileHash, err := calculateFileHash(gotInfo.Name)
				if err != nil {
					t.Fatalf("calculateFileHash() error = %v", err)
				}

				if fileHash != randomFileHash && !tt.dontCheckHashMatch {
					t.Errorf("calculateFileHash() = \n%v, want \n%v", fileHash, randomFileHash)
				}

				if tt.compareInfo != nil {
					tt.compareInfo(gotInfo, tt.wantInfo)
				}
			},
		)
	}
}

func createTempReadWriteFile() (*os.File, error) {
	// Create a temporary file with a ".tmp" extension in the default temp folder path
	tempFile, err := os.CreateTemp("", "*.tmp")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	// Open the file for reading and writing
	file, err := os.OpenFile(tempFile.Name(), os.O_RDWR, 0600)
	if err != nil {
		// Clean up the temporary file on error
		_ = tempFile.Close()
		_ = os.Remove(tempFile.Name())
		return nil, fmt.Errorf("failed to open temp file for read/write: %w", err)
	}
	return file, nil
}

func calculateFileHash(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	fileData, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return fmt.Sprintf("%x", sha256.Sum256(fileData)), nil
}
