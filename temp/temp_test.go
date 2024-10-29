package temp

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"testing"
	"wget/info"
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
	err := os.Setenv("TMPDIR", "custom-temp-dir")
	err = os.Setenv("TMP", "custom-temp-dir")
	defer unsetEnvs("TMPDIR", "TMP")

	if err != nil {
		t.Fatalf("unexpected error setting TMP directory: %v", err)
	}

	// Call Dir() to get the temp directory
	dir := Dir()
	expected := filepath.Join("custom-temp-dir", info.Org)
	// Check if the directory matches the custom path
	if filepath.Join(dir) != expected {
		t.Errorf("expected %s, got %s", expected, dir)
	}
}

func unsetEnvs(args ...string) {
	for _, arg := range args {
		_ = os.Unsetenv(arg)
	}
}

// Test concurrent access to Dir()
func TestDirConcurrency(t *testing.T) {
	defer unsetEnvs("TMPDIR", "TMP")
	const numRoutines = 100
	var wg sync.WaitGroup

	// Launch multiple goroutines to call Dir simultaneously
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// reset the `tempDir` after 10 iterations
			if (i+1)%10 == 0 {
				err := os.Setenv("TMPDIR", "")
				err = os.Setenv("TMP", "")
				if err != nil {
					t.Errorf("unexpected error setting TMP directory: %v", err)
					return
				}
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
	// Set TMP Dir to an invalid path to force an error
	err := os.Setenv("TMPDIR", "/invalid-path/")
	err = os.Setenv("TMP", "/invalid-path/")
	defer unsetEnvs("TMPDIR", "TMP")
	if err != nil {
		t.Errorf("unexpected error setting TMP directory: %v", err)
		return
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, but did not occur")
		}
	}()

	// Call File() to trigger the error
	_, err = File()
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
