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

func TestPrependDotIfInCurrentDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "File in current directory",
			args: args{path: "file.txt"},
			want: "./file.txt",
		},

		{
			name: "Dot file in current directory",
			args: args{path: ".file.txt"},
			want: "./.file.txt",
		},

		{
			name: "Dot dot file in current directory",
			args: args{path: "..file.txt"},
			want: "./..file.txt",
		},

		{
			name: "Tilde file in current directory",
			args: args{path: "~file.txt"},
			want: "./~file.txt",
		},

		{
			name: "Double Tilde file in current directory",
			args: args{path: "~~file.txt"},
			want: "./~~file.txt",
		},

		{
			name: "File in current directory (with spaces)",
			args: args{path: "file one to two.txt"},
			want: "./file one to two.txt",
		},

		{
			name: "File in current directory (with special characters)",
			args: args{path: "file one to two#$2@1.txt"},
			want: "./file one to two#$2@1.txt",
		},

		{
			name: "File in current directory (nested)",
			args: args{path: "nested/path/to/file one to two#$2@1.txt"},
			want: "./nested/path/to/file one to two#$2@1.txt",
		},

		{
			name: "File in parent directory",
			args: args{path: "../file.txt"},
			want: "../file.txt",
		},

		{
			name: "File in parent directory (nested)",
			args: args{path: "../nested/path/to/file.txt"},
			want: "../nested/path/to/file.txt",
		},

		{
			name: "File in parent directory (nested)",
			args: args{path: "../../nested/path/to/file.txt"},
			want: "../../nested/path/to/file.txt",
		},

		{
			name: "File in absolute path",
			args: args{path: "/file.txt"},
			want: "/file.txt",
		},

		{
			name: "File in absolute path (nested)",
			args: args{path: "/home/user/file.txt"},
			want: "/home/user/file.txt",
		},

		{
			name: "File in user home",
			args: args{path: "~/file.txt"},
			want: "~/file.txt",
		},

		{
			name: "File in user home (nested)",
			args: args{path: "~/nested/path/to/file.txt"},
			want: "~/nested/path/to/file.txt",
		},

		{
			name: "Root folder",
			args: args{path: "/"},
			want: "/",
		},

		{
			name: "Path leading to current dir",
			args: args{path: "user/home/../../file.txt"},
			want: "./file.txt",
		},

		{
			name: "Path not leading to current dir",
			args: args{path: "user/home/../../../file.txt"},
			want: "../file.txt",
		},

		{
			name: "Empty string",
			args: args{path: ""},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrependDotIfInCurrentDir(tt.args.path); got != tt.want {
				t.Errorf("PrependDotIfInCurrentDir() = %q, want %q", got, tt.want)
			}
		})
	}
}

func ExamplePrependDotIfInCurrentDir() {
	result := PrependDotIfInCurrentDir("exampleFile.txt")
	fmt.Println(result)
	// Output:
	// ./exampleFile.txt
}

func TestDoAliasUserDir(t *testing.T) {
	type args struct {
		path     string
		homePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty path",
			args: args{
				path:     "",
				homePath: "/home/user",
			},
			want: "",
		},

		{
			name: "Empty path, empty homePath",
			args: args{
				path:     "",
				homePath: "",
			},
			want: "",
		},

		{
			name: "Downloads folder",
			args: args{
				path:     "/home/user/Downloads",
				homePath: "/home/user",
			},
			want: "~/Downloads",
		},

		{
			name: "Downloads folder (already with tilde)",
			args: args{
				path:     "~/Downloads",
				homePath: "/home/user",
			},
			want: "~/Downloads",
		},

		{
			name: "Downloads folder (trailing slashes)",
			args: args{
				path:     "/home/user/Downloads/",
				homePath: "/home/user/",
			},
			want: "~/Downloads",
		},

		{
			name: "Downloads folder (invalid user)",
			args: args{
				path:     "/home/user/Downloads/",
				homePath: "/home/user2/",
			},
			want: "/home/user/Downloads/",
		},

		{
			name: "Downloads folder (invalid user in path)",
			args: args{
				path:     "/home/user2/Downloads/",
				homePath: "/home/user/",
			},
			want: "/home/user2/Downloads",
		},

		{
			name: "opt folder (invalid user in path)",
			args: args{
				path:     "/opt/home/user/Downloads/",
				homePath: "/home/user/",
			},
			want: "/opt/home/user/Downloads/",
		},

		{
			name: "Root path as homePath",
			args: args{
				path:     "/opt/home/user/Downloads/",
				homePath: "/",
			},
			want: "/opt/home/user/Downloads",
		},

		{
			name: "Empty homePath",
			args: args{
				path:     "/opt/home/user/Downloads/",
				homePath: "",
			},
			want: "/opt/home/user/Downloads/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DoAliasUserDir(tt.args.path, tt.args.homePath); got != tt.want {
				t.Errorf("DoAliasUserDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAliasUserDir(t *testing.T) {
	// For constituency among user's, let the default user HOME path be `/home/user`
	err := os.Setenv("HOME", "/home/user")
	defer func() {
		err := os.Unsetenv("HOME")
		if err != nil {
			t.Fatal(err)
		}
	}()

	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// DRY: fewer test case, the other test cases have been covered by the DoAliasUserDir tests
		{
			name: "Downloads folder",
			args: args{
				path: "/home/user/Downloads",
			},
			want: "~/Downloads",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AliasUserDir(tt.args.path); got != tt.want {
				t.Errorf("AliasUserDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
