// Returns existing tempDir if already set
package temp

import (
	"os"
	"testing"
)

func TestReturnsExistingTempDir(t *testing.T) {
	expectedDir := "/tmp/org.zone01.wget"
	result := Dir()
	if result != expectedDir {
		t.Errorf("Expected %s, but got %s", expectedDir, result)
	}
}

func TestHandlesMkdirAllError(t *testing.T) {
	originalTempDir := os.TempDir()
	defer os.Setenv("TMPDIR", originalTempDir)

	os.Setenv("TMPDIR", "/invalid/path")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to MkdirAll failure, but did not panic")
		}
	}()

	Dir()
}

func TestFileSuccess(t *testing.T) {
	// Mocked Dir function for testing
	Dir := func() string {
		return os.TempDir() // Default to system temp dir
	}
	// temporary directory
	tempDir := t.TempDir()

	originalDir := Dir
	Dir = func() string { return tempDir }
	defer func() { Dir = originalDir }()

	file, err := File()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if file == nil {
		t.Fatal("expected a file, got nil")
	}

	file.Close()

	if _, err := os.Stat(file.Name()); os.IsNotExist(err) {
		t.Fatalf("expected file %s to exist, but it does not", file.Name())
	}
}

func TestFileError(t *testing.T) {
	// Mocked Dir function for testing
	Dir := func() string {
		return os.TempDir() // Default to system temp dir
	}

	originalDir := Dir
	Dir = func() string { return "/non/existent/directory" }
	defer func() { Dir = originalDir }()

	file, err := File()

	if err == nil {
		t.Fatal("expected an error, got none")
	}
	if file != nil {
		t.Fatal("expected file to be nil, got a file")
	}
}
