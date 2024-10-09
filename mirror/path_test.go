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

	ContentTypeDisposition := func(contentType, disposition string) http.Header {
		withContentDisposition := http.Header{}

		withContentDisposition.Set("Content-Type", contentType+"; charset=utf-8")
		withContentDisposition.Set("Content-Disposition", "attachment; filename="+disposition)

		return withContentDisposition
	}

	headerText := http.Header{}
	{
		headerText.Set("Content-Type", "text/plain; charset=utf-8")
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
			name: "Root URL",
			args: args{
				downloadUrl:  "http://example.com",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "/downloads/example.com/index.html",
		},

		{
			name: "Root URL (trailing slash)",
			args: args{
				downloadUrl:  "http://example.com/",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "/downloads/example.com/index.html",
		},

		{
			name: "Root URL (host:port)",
			args: args{
				downloadUrl:  "http://example.com:9000",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "/downloads/example.com/index.html",
		},

		{
			name: "Root URL (host:port) (trailing slash)",
			args: args{
				downloadUrl:  "http://example.com:8000/",
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
			name: "Valid URL with path (trailing slash)",
			args: args{
				downloadUrl:  "http://example.com/resource/file.txt/",
				header:       headerText,
				parentFolder: "/downloads",
			},
			want: "/downloads/example.com/resource/file.txt/index.txt",
		},

		{
			name: "Valid wget URL with path (trailing slash)",
			args: args{
				downloadUrl:  "http://example.com/wget/",
				header:       headerHTML,
				parentFolder: "downloads",
			},
			want: "downloads/example.com/wget/index.html",
		},

		{
			name: "Valid wget URL with path (no trailing slash)",
			args: args{
				downloadUrl:  "http://example.com/wget",
				header:       headerHTML,
				parentFolder: "downloads",
			},
			want: "downloads/example.com/wget",
		},

		{
			name: "Valid wget CSS URL with path (no trailing slash)",
			args: args{
				downloadUrl:  "http://example.com/wget",
				header:       headerCSS,
				parentFolder: "downloads",
			},
			want: "downloads/example.com/wget",
		},

		{
			name: "Valid wget CSS URL with path (trailing slash)",
			args: args{
				downloadUrl:  "http://example.com/wget/",
				header:       headerCSS,
				parentFolder: "downloads",
			},
			want: "downloads/example.com/wget/index.css",
		},

		{
			name: "URL with query parameters",
			args: args{
				downloadUrl:  "http://example.com/resource/file.txt?version=1",
				header:       headerText,
				parentFolder: "/downloads",
			},
			want: "/downloads/example.com/resource/file.txt",
		},
		{
			name: "URL with special characters",
			args: args{
				downloadUrl:  "http://example.com/resource/file@name.txt",
				header:       headerText,
				parentFolder: "/downloads",
			},
			want: "/downloads/example.com/resource/file@name.txt",
		},
		{
			name: "Empty URL",
			args: args{
				downloadUrl:  "",
				header:       headerHTML,
				parentFolder: "/downloads",
			},
			want: "",
		},

		{
			name: "Non-HTTP scheme",
			args: args{
				downloadUrl:  "ftp://invalid-scheme",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "/downloads/invalid-scheme/index.html",
		},

		{
			name: "Non-HTTP scheme trailing slash",
			args: args{
				downloadUrl:  "ftp://invalid-scheme/",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "/downloads/invalid-scheme/index.html",
		},

		{
			name: "Invalid URL trailing slash",
			args: args{
				downloadUrl:  "https://example.com:xxxx/",
				header:       headerHTML,
				parentFolder: "/downloads",
			},
			want: "",
		},

		{
			name: "Invalid URL",
			args: args{
				downloadUrl:  "https://example.com:xxxx",
				header:       headerHTML,
				parentFolder: "/downloads",
			},
			want: "",
		},

		{
			name: "Incomplete URL",
			args: args{
				downloadUrl:  "https://",
				header:       headerHTML,
				parentFolder: "/downloads",
			},
			want: "",
		},

		{
			name: "Test case 1 trailing slash",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/css-beer/",
				header:       headerHTML,
				parentFolder: "/downloads",
			},
			want: "/downloads/wizard254.github.io/wget/css-beer/index.html",
		},

		{
			name: "Test case 1",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/css-beer",
				header:       headerHTML,
				parentFolder: "/downloads",
			},
			want: "/downloads/wizard254.github.io/wget/css-beer",
		},

		{
			name: "Test case 1 trailing slash (CSS)",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/css-beer/",
				header:       headerCSS,
				parentFolder: "/downloads",
			},
			want: "/downloads/wizard254.github.io/wget/css-beer/index.css",
		},

		{
			name: "Test case 1 trailing slash (stream)",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/beer/",
				header:       headerStream,
				parentFolder: "/downloads",
			},
			want: "/downloads/wizard254.github.io/wget/beer/index.html",
		},

		{
			name: "Test case 1 no trailing slash (stream)",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/beer",
				header:       headerStream,
				parentFolder: "/downloads",
			},
			want: "/downloads/wizard254.github.io/wget/beer",
		},

		{
			name: "Test case 1 no trailing slash (CSS)",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/beer",
				header:       headerCSS,
				parentFolder: "/downloads",
			},
			want: "/downloads/wizard254.github.io/wget/beer",
		},

		{
			name: "css-beer stylesheet",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/css-beer/style.css",
				header:       http.Header{},
				parentFolder: "/downloads",
			},
			want: "/downloads/wizard254.github.io/wget/css-beer/style.css",
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

		{
			name: "Save path does not exist (Wrong Content-type)",
			args: args{
				downloadUrl:  "http://example.com/resource/file.txt",
				header:       headerHTML,
				parentFolder: "/nonexistent/path",
			},
			want: "/nonexistent/path/example.com/resource/file.txt",
		},

		{
			name: "Test case 2 no trailing slash (text) disposition",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/beer",
				header:       ContentTypeDisposition("text/plain", "beer.txt"),
				parentFolder: "/downloads",
			},
			want: "/downloads/wizard254.github.io/wget/beer.txt",
		},

		{
			name: "Test case 2 no trailing slash (text) disposition markdown",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/beer",
				header:       ContentTypeDisposition("text/plain", "beer.md"),
				parentFolder: "/downloads",
			},
			want: "/downloads/wizard254.github.io/wget/beer.md",
		},

		{
			name: "Test case 2 no trailing slash (stream) disposition",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/beer",
				header:       ContentTypeDisposition("octet-stream", "beer.txt"),
				parentFolder: "/downloads",
			},
			want: "/downloads/wizard254.github.io/wget/beer.txt",
		},

		{
			name: "Test case 2 no trailing slash (stream) disposition markdown",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/beer",
				header:       ContentTypeDisposition("octet-stream", "beer.md"),
				parentFolder: "/downloads",
			},
			want: "/downloads/wizard254.github.io/wget/beer.md",
		},

		{
			name: "Test case 2 no trailing slash (stream) malformed disposition",
			args: args{
				downloadUrl:  "https://wizard254.github.io/wget/beer",
				header:       ContentTypeDisposition("octet-stream", ""),
				parentFolder: "/downloads",
			},
			want: "/downloads/wizard254.github.io/wget/beer",
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
