package downloader

import (
	"os"
	"path/filepath"
	"testing"

	"wget/ctx"
	"wget/httpx"
)

// func TestGetResource_Success(t *testing.T) {
// 	// Create a temporary directory to save the file.
// 	tempDir := t.TempDir()

// 	// Set up the context with necessary fields.
// 	c := ctx.Context{
// 		OutputFile: "",
// 		SavePath:   tempDir,
// 	}

// 	a := arg{Context: &c}

// 	// Mock the HTTP server for testing.
// 	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("Test file content"))
// 	})
// 	go http.ListenAndServe(":8080", nil)

// 	// Test the resource download.
// 	err := a.GetResource("http://localhost:8080/file")
// 	if err != nil {
// 		t.Fatalf("GetResource failed: %v", err)
// 	}

// 	// Verify file is saved correctly.
// 	expectedFile := filepath.Join(tempDir, "file")
// 	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
// 		t.Fatalf("Expected file %s not found", expectedFile)
// 	}
// }

func TestCheckIfFileExists(t *testing.T) {
	tempDir := t.TempDir()

	// Create a temporary file.
	fname := filepath.Join(tempDir, "file.txt")
	os.WriteFile(fname, []byte("test content"), 0o644)

	// Test that CheckIfFileExists adds a number to duplicate filenames.
	newFileName := CheckIfFileExists(fname)
	expected := filepath.Join(tempDir, "file(1).txt")
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
		{int64(httpx.MB), "1.00MB"},
		{int64(httpx.GB), "1.00GB"},
	}

	for _, tc := range testCases {
		result := httpx.RoundOfSizeOfData(tc.bytes)
		if result != tc.expected {
			t.Fatalf("Expected %s but got %s", tc.expected, result)
		}
	}
}

func TestCheckIfFileExists_Multi(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"test1", "file.txt", "file(1).txt"},
		{"test2", "file", "file(1)"},
		{"test3", "", ""},
		{"test4", "hello_world.png", "hello_world(1).png"},
		{"test5", "20MB.zip", "20MB(1).zip"},
		{"test6", ".zip", "(1).zip"},
	}

	create := func(fpath string) {
		os.Create(fpath)
	}

	destroy := func(fpath string) {
		os.Remove(fpath)
	}
	for _, tt := range tests {

		if tt.input != "" {
			create(tt.input)
			defer destroy(tt.input)
		}

		t.Run(tt.name, func(t *testing.T) {

			got := CheckIfFileExists(tt.input)

			if got != tt.want {
				t.Errorf("CheckIfFileExists(%q) Failed got %q want %q\n", tt.input, got, tt.want)
			}
		})
	}
}
