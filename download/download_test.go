package download

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Mock server to simulate file download
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
	defer os.Remove(filename)

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
	err := ToFile(mockServer.URL, "testfile.txt")
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
