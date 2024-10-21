package temp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"testing"
)

// Test creating files in the temp directory
func TestFileCreation(t *testing.T) {
	// Create a temporary file using File()
	file, err := File()
	if err != nil {
		t.Fatalf("unexpected error creating temp file: %v", err)
	}
	defer func(name string) {
		_ = os.Remove(name)
	}(file.Name()) // Cleanup

	// Check if the file was created
	if _, err := os.Stat(file.Name()); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist, but it does not", file.Name())
	}

	// Check if the file name has the correct pattern
	re := regexp.MustCompile(`.*\.wget\.tmp`)
	if !re.MatchString(file.Name()) {
		t.Errorf("expected file name to match pattern *.wget.tmp, got %s", file.Name())
	}
}

// Test Dir() when tempDir is manually set
func TestDirWithCustomTempDir(t *testing.T) {
	// Set a custom temp dir
	tempDir = filepath.Join(os.TempDir(), "custom-temp-dir")
	defer func() {
		tempDir = ""
	}()

	// Call Dir() to get the temp directory
	dir := Dir()

	// Check if the directory matches the custom path
	if dir != tempDir {
		t.Errorf("expected %s, got %s", tempDir, dir)
	}
}

// Test concurrent access to Dir()
func TestDirConcurrency(t *testing.T) {
	const numRoutines = 100
	var wg sync.WaitGroup

	// Launch multiple goroutines to call Dir simultaneously
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// reset the `tempDir` after 10 iterations
			if (i+1)%10 == 0 {
				tempDir = ""
			}
			dir := Dir()
			if dir == "" {
				t.Error("Dir() returned an empty string")
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

// Test concurrent file creation
func TestFileConcurrency(t *testing.T) {
	const numFiles = 100
	var wg sync.WaitGroup
	fileNames := make(map[string]struct{})
	var mu sync.Mutex

	// Launch multiple goroutines to create temporary files
	for i := 0; i < numFiles; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			file, err := File()
			if err != nil {
				t.Errorf("unexpected error creating temp file: %v", err)
				return
			}
			defer func() {
				_ = os.Remove(file.Name())
			}()

			// Store the file name in the map (protected by a mutex)
			mu.Lock()
			fileNames[file.Name()] = struct{}{}
			mu.Unlock()
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Ensure all file names are unique
	if len(fileNames) != numFiles {
		t.Errorf("expected %d unique files, but got %d", numFiles, len(fileNames))
	}
}

// Test error handling when directory creation fails
func TestDirFailure(t *testing.T) {
	// Set tempDir to an invalid path to force an error
	tempDir = "/invalid-path/"
	defer func() {
		tempDir = ""
	}()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, but did not occur")
		}
	}()

	// Call File() to trigger the error
	_, err := File()
	if err != nil {
		panic(err)
	}
}

func ExampleFile() {
	file, err := File()
	if err != nil {
		panic(err)
	}

	// remove the temporary file when we are done with it
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err)
		}
	}(file.Name())

	// Now read and write to the file
	n, err := file.WriteString("Hello World")
	if err != nil {
		panic(err)
	}

	fmt.Println(n)
	// Output:
	// 11
}

func TestReturnsExistingTempDir(t *testing.T) {
	expectedDir := "/tmp/org.zone01.wget"
	result := Dir()
	if result != expectedDir {
		t.Errorf("Expected %s, but got %s", expectedDir, result)
	}
}

func TestHandlesMkdirAllError(t *testing.T) {
	originalTempDir := os.TempDir()
	defer func(key, value string) {
		err := os.Setenv(key, value)
		if err != nil {
			log.Fatalf("failed to set environment variable: %v\n", err)
		}
	}("TMPDIR", originalTempDir)

	// get the tmp dir before changing the TMPDIR environment variable
	dir := Dir()

	err := os.Setenv("TMPDIR", "/invalid/path")
	if err != nil {
		log.Fatalf("failed to set environment variable: %v\n", err)
	}

	// When the TMPDIR environment variable is changed, but the Dir() function was
	// already called, then, Dir will stick to the first TMPDIR it succeeded in
	// creating
	if Dir() != dir {
		t.Errorf("expected %s, but got %s", originalTempDir, Dir())
	}
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

	err = file.Close()
	if err != nil {
		log.Fatalf("failed to close temp file [%q]: %v\n", file.Name(), err)
	}

	if _, err := os.Stat(file.Name()); os.IsNotExist(err) {
		t.Fatalf("expected file %s to exist, but it does not", file.Name())
	}
}
