package download

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// createMockServer sets up a mock server with the given responseCode and responseBody
func createMockServer(responseCode int, responseBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(responseCode)
		w.Write([]byte(responseBody))
	}))
}

func TestToFile_Success(t *testing.T) {
	// Create a mock server that returns a valid response
	mockServer := createMockServer(http.StatusOK, "mock file content")
	defer mockServer.Close()

	// Call the ToFile function with the mock server URL
	filename := "testfile.txt"
	err := ToFile(mockServer.URL, filename)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer os.Remove(filename) // Ensure the file is removed after the test

	// Check if the file is created and has the correct content
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}

	expectedContent := "mock file content"
	if string(content) != expectedContent {
		t.Errorf("expected file content %q, got %q", expectedContent, content)
	}
}

func TestToFile_HttpError(t *testing.T) {
	// Create a mock server that returns a 404 error
	mockServer := createMockServer(http.StatusNotFound, "not found")
	defer mockServer.Close()

	// Call the ToFile function with the mock server URL
	fmt.Println(mockServer.URL)
	err := ToFile(mockServer.URL, "testfile.txt")
	fmt.Println(">>><<<",err,err.Error())
	if err == nil || err.Error() != "error downloading file: 404 Not Found" {
		t.Errorf("expected error downloading file: 404 Not Found, got %v", err)
	}
}

func TestToFile_FileCreationError(t *testing.T) {
	// Create a mock server that returns a valid response
	mockServer := createMockServer(http.StatusOK, "mock file content")
	defer mockServer.Close()

	// Call the ToFile function with an invalid filename
	err := ToFile(mockServer.URL, "/invalidpath/testfile.txt")
	if err == nil || err.Error() != "error creating file: open /invalidpath/testfile.txt: no such file or directory" {
		t.Errorf("expected error creating file, got %v", err)
	}
}

func TestToFile_FileWriteError(t *testing.T) {
	// Create a mock server that returns a valid response
	mockServer := createMockServer(http.StatusOK, "mock file content")
	defer mockServer.Close()

	// Use a file that cannot be written to simulate the error
	filename := "/dev/full" // special file that simulates a full disk
	err := ToFile(mockServer.URL, filename)
	if err == nil || err.Error() != "error saving file: write /dev/full: no space left on device" {
		t.Errorf("expected error saving file, got %v", err)
	}
}

// TestGetCurrentTime checks whether GetCurrentTime returns the current time and in YYYY-MM-DD HH:MM:SS format
func TestGetCurrentTime(t *testing.T) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	if formattedTime != GetCurrentTime() {
		t.Error("error getting correct time")
	}
}

// mockserver
func mockServer(statusCode int, body string) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		if body != "" {
			w.Write([]byte(body))
		}
	})
	return httptest.NewServer(handler)
}

func TestDownloadUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Valid URL",
			args:    args{url: mockServer(http.StatusOK, "test content").URL},
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			args:    args{url: "invalid-url"},
			wantErr: true,
		},
		{
			name:    "Non-Existent URL",
			args:    args{url: mockServer(http.StatusNotFound, "").URL},
			wantErr: true,
		},
		{
			name:    "Server Error URL",
			args:    args{url: mockServer(http.StatusInternalServerError, "").URL},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DownloadUrl(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("DownloadUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
