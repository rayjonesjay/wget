package downloader

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"wget/ctx"
)


func TestGetResource_Success(t *testing.T) {
	// Create a temporary directory to save the file.
	tempDir := t.TempDir()

	// Set up the context with necessary fields.
	c := ctx.Context{
		OutputFile: "",
		SavePath:   tempDir,
	}

	a := arg{Context: &c}

	// Mock the HTTP server for testing.
	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test file content"))
	})
	go http.ListenAndServe(":8080", nil)

	// Test the resource download.
	err := a.GetResource("http://localhost:8080/file")
	if err != nil {
		t.Fatalf("GetResource failed: %v", err)
	}

	// Verify file is saved correctly.
	expectedFile := filepath.Join(tempDir, "file")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Fatalf("Expected file %s not found", expectedFile)
	}
}

func TestCheckIfFileExists(t *testing.T) {
	tempDir := t.TempDir()

	// Create a temporary file.
	fname := filepath.Join(tempDir, "file.txt")
	os.WriteFile(fname, []byte("test content"), 0o644)

	// Test that CheckIfFileExists adds a number to duplicate filenames.
	newFileName := CheckIfFileExists(fname)
	expected := filepath.Join(tempDir, "file1.txt")
	if newFileName != expected {
		t.Fatalf("Expected %s but got %s", expected, newFileName)
	}
}

func TestIsEmpty(t *testing.T) {
	c := ctx.Context{}
	a := arg{Context: &c}

	// Test for string type
	if !a.IsEmpty("") {
		t.Fatal("Expected string to be empty")
	}

	// Test for slice type
	if !a.IsEmpty([]int{}) {
		t.Fatal("Expected slice to be empty")
	}

	// Test for array type
	if !a.IsEmpty([0]int{}) {
		t.Fatal("Expected array to be empty")
	}

	// Test for non-empty data
	if a.IsEmpty("data") {
		t.Fatal("Expected string to not be empty")
	}
}

func TestDetermineOutputPath(t *testing.T) {
	c := ctx.Context{
		OutputFile: "",
		SavePath:   "",
	}
	a := arg{Context: &c}

	url := "http://example.com/file.txt"
	expected := "file.txt"
	result := a.determineOutputPath(url)
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}

	c.SavePath = "/path/to/dir"
	expected = "/path/to/dir/file.txt"
	result = a.determineOutputPath(url)
	if result != expected {
		t.Fatalf("Expected %s but got %s", expected, result)
	}
}

func TestRoundOfSizeOfData(t *testing.T) {
	testCases := []struct {
		bytes    int64
		expected string
	}{
		{MB, "1.00MB"},
		{GB, "1.00GB"},
	}

	for _, tc := range testCases {
		result := RoundOfSizeOfData(tc.bytes)
		if result != tc.expected {
			t.Fatalf("Expected %s but got %s", tc.expected, result)
		}
	}
}
