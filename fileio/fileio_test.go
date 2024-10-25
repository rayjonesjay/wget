package fileio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
)

// MockReadWriteCloser implements io.ReadWriteCloser and allows us to control the Close() behavior for testing
type MockReadWriteCloser struct {
	closed     bool
	closeError error
}

func (m *MockReadWriteCloser) Read([]byte) (n int, err error) {
	return 0, nil
}

func (m *MockReadWriteCloser) Write([]byte) (n int, err error) {
	return 0, nil
}

func (m *MockReadWriteCloser) Close() error {
	m.closed = true
	return m.closeError
}

func TestClose(t *testing.T) {
	m := &MockReadWriteCloser{}
	Close(m)

	if !m.closed {
		t.Error("Expected Close() to be called on the ReadWriteCloser")
	}
}

func TestShouldClose1(t *testing.T) {
	// Test with no error, no handler
	t.Run(
		"Test with no error", func(t *testing.T) {
			m := &MockReadWriteCloser{}
			ShouldClose(m)

			if !m.closed {
				t.Error("Expected MustClose() to call Close() on the ReadWriteCloser")
			}
		},
	)

	// Test with no error, with Exit handler
	t.Run(
		"Test with no error", func(t *testing.T) {
			m := &MockReadWriteCloser{}
			// Since there is no error, the Exit handler should not be called, and thus the Test should continue
			ShouldClose(m)

			if !m.closed {
				t.Error("Expected MustClose() to call Close() on the ReadWriteCloser")
			}
		},
	)
}

func TestMustClose(t *testing.T) {
	Exit := func() {
		os.Exit(1)
	}

	count := 0
	IncrementCount := func() {
		count++
	}

	// Test with no error, no handler
	t.Run(
		"Test with no error", func(t *testing.T) {
			m := &MockReadWriteCloser{}
			MustClose(m, nil)

			if !m.closed {
				t.Error("Expected MustClose() to call Close() on the ReadWriteCloser")
			}
		},
	)

	// Test with no error, with Exit handler
	t.Run(
		"Test with no error", func(t *testing.T) {
			m := &MockReadWriteCloser{}
			// Since there is no error, the Exit handler should not be called, and thus the Test should continue
			MustClose(m, Exit)

			if !m.closed {
				t.Error("Expected MustClose() to call Close() on the ReadWriteCloser")
			}
		},
	)

	stderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w

	// Test with error
	t.Run(
		"Test with error", func(t *testing.T) {
			closeError := "simulated close error"
			m := &MockReadWriteCloser{closeError: errors.New(closeError)}

			i := count
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				// This should panic and cause the go routine to exit
				MustClose(m, IncrementCount)
			}()
			wg.Wait()

			err := w.Close()
			if err != nil {
				t.Fatal(err)
			}
			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			if err != nil {
				t.Fatal(err)
			}
			output := buf.String()

			if !m.closed {
				t.Error("Expected MustClose() to call Close() on the ReadWriteCloser even with error")
			}

			expectedOutput := fmt.Sprintf("%s\n", closeError)
			if output != expectedOutput {
				t.Errorf("Expected stderr output:\n%s\nGot:\n%s", expectedOutput, output)
			}

			if i == count {
				t.Error("Expected MustClose() to call IncrementCount on error")
			}
		},
	)

	// Restore stderr
	os.Stderr = stderr
}
