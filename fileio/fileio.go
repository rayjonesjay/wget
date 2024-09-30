// Package fileio provides utility functions to help in handling file input/output operations and errors
package fileio

import (
	"fmt"
	"io"
	"os"
)

// Close the given closable reader or writer, ignoring any close errors that may arise.
// This wrapper is essential in defer statements,
// when we actually don't need to care whether the file was closed properly
func Close(closable io.ReadWriteCloser) {
	_ = closable.Close()
}

// ShouldClose ensures the given reader or writer is closed properly, and that,
// should any close errors arise, the error message should be printed to the error stream
func ShouldClose(closable io.ReadWriteCloser) {
	MustClose(closable, nil)
}

// MustClose ensures the given reader or writer is closed properly, and that,
// should any close errors arise the error message should be printed just before the given close error
// handler function is called.
func MustClose(closable io.ReadWriteCloser, handler func()) {
	err := closable.Close()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		if handler != nil {
			handler()
		}
	}
}
