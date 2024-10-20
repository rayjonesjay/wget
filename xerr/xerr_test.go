// Write error message to standard error successfully
package xerr_test

import (
	"bytes"
	"os"
	"testing"

	"wget/xerr"
)

func TestWriteErrorMessageToStderr(t *testing.T) {
	originalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	errorMessage := "Test error message"
	xerr.WriteError(errorMessage, 0, false)

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stderr = originalStderr

	if buf.String() != errorMessage+"\n" {
		t.Errorf("Expected %v, got %v", errorMessage+"\n", buf.String())
	}
}

// Handle non-string errorMessage inputs gracefully

func TestHandleNonStringErrorMessage(t *testing.T) {
	originalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	errorMessage := 12345
	xerr.WriteError(errorMessage, 0, false)

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stderr = originalStderr

	expectedOutput := "12345\n"
	if buf.String() != expectedOutput {
		t.Errorf("Expected %v, got %v", expectedOutput, buf.String())
	}
}
