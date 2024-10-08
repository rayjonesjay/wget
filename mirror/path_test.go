package mirror

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestGetFile(t *testing.T) {
	headerHTML := http.Header{}
	{
		headerHTML.Set("Content-Type", "text/html; charset=utf-8")
	}

	headerCSS := http.Header{}
	{
		headerCSS.Set("Content-Type", "text/css; charset=utf-8")
	}

	headerStream := http.Header{}
	{
		headerStream.Set("Content-Type", "octet-stream; charset=utf-8")
	}

	type args struct {
		downloadUrl  string
		header       http.Header
		parentFolder string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid URL with no path",
			args: args{
				downloadUrl:  "http://example.com",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "/downloads/example.com/index.html",
		},
		{
			name: "Valid URL with path",
			args: args{
				downloadUrl:  "http://example.com/resource/file.txt",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "/downloads/example.com/resource/file.txt",
		},

		{
			// TODO
			name: "Valid URL with path (trailing slash)",
			args: args{
				downloadUrl:  "http://example.com/wget/",
				header:       headerHTML,
				parentFolder: "downloads",
			},
			want: "downloads/example.com/wget/index.html",
		},

		{
			name: "Valid URL with path (no trailing slash)",
			args: args{
				downloadUrl:  "http://example.com/wget",
				header:       headerHTML,
				parentFolder: "downloads",
			},
			want: "downloads/example.com/wget/index.html",
		},

		{
			name: "Valid URL with path (no trailing slash)",
			args: args{
				downloadUrl:  "http://example.com/wget",
				header:       headerCSS,
				parentFolder: "downloads",
			},
			want: "downloads/example.com/wget",
		},

		{
			name: "URL with query parameters",
			args: args{
				downloadUrl:  "http://example.com/resource/file.txt?version=1",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "/downloads/example.com/resource/file.txt",
		},
		{
			name: "URL with special characters",
			args: args{
				downloadUrl:  "http://example.com/resource/file@name.txt",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "/downloads/example.com/resource/file@name.txt",
		},
		{
			name: "Empty URL",
			args: args{
				downloadUrl:  "",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "/downloads/index.html",
		},
		{
			name: "Invalid URL",
			args: args{
				downloadUrl:  "htp://invalid-url",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "/downloads/invalid-url/index.html",
		},

		{
			// TODO: TWO
			name: "Invalid URL",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/css-beer/",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "wizard254.github.io/wget/css-beer/index.html",
		},

		{
			// TODO: TWO
			name: "Invalid URL",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/css-beer/style.css",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "wizard254.github.io/wget/css-beer/style.css",
		},

		{
			// TODO: TWO
			name: "Invalid URL",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "wizard254.github.io/wget/index.html",
		},

		{
			name: "Save path does not exist",
			args: args{
				downloadUrl:  "http://example.com/resource/file.txt",
				header:       http.Header{},
				parentFolder: "/nonexistent/path",
			},
			want: "/nonexistent/path/example.com/resource/file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := GetFile(tt.args.downloadUrl, tt.args.header, tt.args.parentFolder); got != tt.want {
					t.Errorf("GetFile() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

// Example test case
func ExampleGetFile() {
	header := http.Header{}
	header.Set("User-Agent", "my-agent")
	path := GetFile("http://example.com/resource/file.txt", header, "/downloads")
	fmt.Println(path)
	// Output: /downloads/example.com/resource/file.txt
}

func TestFolderStructure(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name          string
		args          args
		wantStructure []string
	}{
		{
			name:          "Simple file path",
			args:          args{filePath: "/a/b/c/file.txt"},
			wantStructure: []string{"/a/b/c", "/a/b", "/a"},
		},
		{
			name:          "Root file path",
			args:          args{filePath: "/file.txt"},
			wantStructure: nil,
		},
		{
			name:          "File in a nested folder (absolute)",
			args:          args{filePath: "/home/user/docs/report.pdf"},
			wantStructure: []string{"/home/user/docs", "/home/user", "/home"},
		},
		{
			name:          "File in a nested folder (relative)",
			args:          args{filePath: "user/docs/report.pdf"},
			wantStructure: []string{"user/docs", "user"},
		},
		{
			name:          "Non-existent file path",
			args:          args{filePath: "/nonexistent/path/to/file.txt"},
			wantStructure: []string{"/nonexistent/path/to", "/nonexistent/path", "/nonexistent"},
		},
		{
			name:          "Empty path",
			args:          args{filePath: ""},
			wantStructure: nil,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if gotStructure := FolderStructure(tt.args.filePath); !reflect.DeepEqual(
					gotStructure, tt.wantStructure,
				) {
					t.Errorf("FolderStructure() = %v, want %v", gotStructure, tt.wantStructure)
				}
			},
		)
	}
}

// ExampleFolderStructure provides an example usage of the FolderStructure function.
func ExampleFolderStructure() {
	structure := FolderStructure("/a/b/c/file.txt")
	for _, folder := range structure {
		fmt.Println(folder)
	}
	// Output:
	// /a/b/c
	// /a/b
	// /a
}
