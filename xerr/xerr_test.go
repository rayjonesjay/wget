// Write error message to standard error successfully
package xerr

import (
	"bytes"
	"os"
	"testing"
)

func TestWriteErrorMessageToStderr(t *testing.T) {
	originalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	errorMessage := "Test error message"
	WriteError(errorMessage, 0, false)

	err := w.Close()
	if err != nil {
		t.Error(err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		t.Error(err)
	}
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
	WriteError(errorMessage, 0, false)

	err := w.Close()
	if err != nil {
		t.Error(err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		t.Error(err)
	}
	os.Stderr = originalStderr

	expectedOutput := "12345\n"
	if buf.String() != expectedOutput {
		t.Errorf("Expected %v, got %v", expectedOutput, buf.String())
	}
}

func Test_writeError(t *testing.T) {
	lastCode := -1
	mockExit := func(code int) {
		lastCode = code
	}

	type args struct {
		errorMessage interface{}
		statusCode   int
		shouldExit   bool
		exit         func(int)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Base",
			args: args{
				errorMessage: "Test message",
				statusCode:   0,
				shouldExit:   false,
				exit:         mockExit,
			},
		},

		{
			name: "Error: exit code 1",
			args: args{
				errorMessage: "Test message",
				statusCode:   1,
				shouldExit:   true,
				exit:         mockExit,
			},
		},

		{
			name: "Error: exit code 0",
			args: args{
				errorMessage: "Test message",
				statusCode:   0,
				shouldExit:   true,
				exit:         mockExit,
			},
		},

		{
			name: "Redundant: exit code 1",
			args: args{
				errorMessage: "Test message",
				statusCode:   1,
				shouldExit:   false,
				exit:         mockExit,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeError(tt.args.errorMessage, tt.args.statusCode, tt.args.shouldExit, tt.args.exit)
			if tt.args.shouldExit {
				if lastCode != tt.args.statusCode {
					t.Errorf("Expected %v, got %v", tt.args.statusCode, lastCode)
				}
			} else {
				if lastCode != -1 {
					t.Errorf("Expected %v, got %v", -1, lastCode)
				}
			}

			lastCode = -1
		})
	}
}
