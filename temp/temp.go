// Package temp defines various utilities to handle the creation and deletion of temporary files
package temp

import (
	"os"
	"path"
	"sync"
	"wget/info"
)

var tempDir = ""
var cm = sync.Mutex{}

// Dir returns the default directory to use for temporary files for the program.
// This is a wrapper around [os.TempDir], so as not to end up polluting the user's temp dir.
func Dir() string {
	if tempDir != "" {
		return tempDir
	}

	cm.Lock()
	defer cm.Unlock()
	tempDir = path.Join(os.TempDir(), info.Org)
	err := os.MkdirAll(tempDir, 0775)
	if err != nil {
		panic(err)
	}
	return tempDir
}

// File creates a new temporary file in the program's temp directory, opens the file for reading and writing,
// and returns the resulting file. The filename is randomly generated,
// with the pattern "*.wget.tmp". Multiple programs or goroutines calling [temp.File] simultaneously will not choose
// the same file. The caller can use the file's Name method to find the pathname of the file.
// It is the caller's responsibility to remove the file when it is no longer needed.
func File() (*os.File, error) {
	// Create a temporary file inside the tmp directory
	tempFile, err := os.CreateTemp(Dir(), "*.wget.tmp")
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}
