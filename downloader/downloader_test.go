package downloader

import (
	"os"
	"testing"
)

func TestCheckIfFileExists(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"test1", "file.txt", "file1.txt"},
		{"test2", "file", "file1"},
		{"test3", "", ""},
		{"test4", "hello_world.png", "hello_world1.png"},
		{"test5", "20MB.zip", "20MB1.zip"},
		{"test6", ".zip", "1.zip"},
	}

	create := func(fpath string) {
		os.Create(fpath)
	}

	destroy := func(fpath string) {
		os.Remove(fpath)
	}
	for _, tt := range tests {

		if tt.input != "" {
			create(tt.input)
			defer destroy(tt.input)
		}

		t.Run(tt.name, func(t *testing.T) {

			got := CheckIfFileExists(tt.input)

			if got != tt.want {
				t.Errorf("CheckIfFileExists(%q) Failed got %q want %q\n", tt.input, got, tt.want)
			}
		})
	}
}
