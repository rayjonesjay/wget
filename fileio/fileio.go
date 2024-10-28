// Package fileio provides utility functions to help in handling file input/output operations and errors
package fileio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Close the given closable reader or writer, ignoring any close errors that may arise.
// This wrapper is essential in defer statements,
// when we actually don't need to care whether the file was closed properly
func Close(closable io.Closer) {
	_ = closable.Close()
}

// ShouldClose ensures the given reader or writer is closed properly, and that,
// should any close errors arise, the error message should be printed to the error stream
func ShouldClose(closable io.Closer) {
	MustClose(closable, nil)
}

// MustClose ensures the given reader or writer is closed properly, and that,
// should any close errors arise the error message should be printed just before the given close error
// handler function is called.
func MustClose(closable io.Closer, handler func()) {
	err := closable.Close()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		if handler != nil {
			handler()
		}
	}
}

// PrependDotIfInCurrentDir checks if the file is in the current directory
// and prepends "./" if it is.
func PrependDotIfInCurrentDir(path string) string {
	if slices.Contains([]string{"", "/", ".", "..", "./", "../", "~/", "~"}, path) {
		return path
	}

	path = filepath.Clean(path)
	// Find the relative path from the current directory
	rel, err := filepath.Rel("./", path)
	if err != nil {
		return path
	}

	hasAnyPrefix := func(rel string) bool {
		p := []string{"/", "./", "../", "~/"}
		for _, prefix := range p {
			if strings.HasPrefix(rel, prefix) {
				return true
			}
		}

		return false
	}

	if hasAnyPrefix(rel) {
		return rel
	}

	return fmt.Sprintf(".%s%s", string(os.PathSeparator), rel)
}

// DoAliasUserDir given a path, and the user's HOME path, returns the given path
// with the user's HOME path replaced with the tilde symbol (~).
//
// Note: This API strips off the `/` suffix, as defined by [filepath.Clean]
func DoAliasUserDir(path, homePath string) string {
	homePath = filepath.Clean(homePath)
	newPath := strings.TrimPrefix(path, homePath)

	if newPath == path {
		return path
	}

	if !strings.HasPrefix(newPath, "/") {
		return filepath.Clean(path)
	}

	return filepath.Join("~", newPath)
}

// AliasUserDir same as [DoAliasUserDir], but gets the user's HOME directory from
// the environment variables as returned by [os.UserHomeDir]
//
// # Example
//
// Code:
//
//	// assuming the target user is `ubuntu` and his HOME path is set to `/home/ubuntu`
//	result := AliasUserDir("/home/ubuntu/Downloads")
//	fmt.Println(result)
//
// Output:
//
//	~/Downloads
func AliasUserDir(path string) string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return DoAliasUserDir(path, homePath)
}
