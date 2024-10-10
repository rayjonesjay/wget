package xurl

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

func TestAbsoluteUrl(t *testing.T) {
	type args struct {
		parentURL   string
		relativeURL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Relative root URL",
			args: args{
				parentURL:   "https://example.com/path/to/file",
				relativeURL: "/",
			},
			want:    "https://example.com",
			wantErr: false,
		},

		{
			name: "Relative root URL path",
			args: args{
				parentURL:   "https://example.com/path/to/file",
				relativeURL: "/another/path/to/file",
			},
			want:    "https://example.com/another/path/to/file",
			wantErr: false,
		},

		{
			name: "Relative path URL (parent)",
			args: args{
				parentURL:   "https://example.com/user/profile/settings",
				relativeURL: "../",
			},
			want:    "https://example.com/user/profile",
			wantErr: false,
		},

		{
			name: "Relative path URL (two level parent)",
			args: args{
				parentURL:   "https://example.com/user/profile/settings",
				relativeURL: "../../",
			},
			want:    "https://example.com/user",
			wantErr: false,
		},

		{
			name: "Relative path URL (two level parent)",
			args: args{
				parentURL:   "https://example.com/user/profile/settings",
				relativeURL: "../../name/firstname",
			},
			want:    "https://example.com/user/name/firstname",
			wantErr: false,
		},

		{
			name: "Relative path URL (parent no trailing slash)",
			args: args{
				parentURL:   "https://example.com/user/profile/settings",
				relativeURL: "..",
			},
			want:    "https://example.com/user/profile",
			wantErr: false,
		},

		{
			name: "Relative path URL (current)",
			args: args{
				parentURL:   "https://example.com/user/profile/settings",
				relativeURL: "./",
			},
			want:    "https://example.com/user/profile/settings",
			wantErr: false,
		},

		{
			name: "Relative path URL (current) with multiple slash separators in parent URL",
			args: args{
				parentURL:   "https://example.com///user//profile////settings//",
				relativeURL: "./",
			},
			want:    "https://example.com/user/profile/settings/",
			wantErr: false,
		},

		{
			name: "Relative path URL (two level parent) with multiple slash separators in parent and relative URLs",
			args: args{
				parentURL:   "https://example.com///user//profile////settings//",
				relativeURL: "..//..///",
			},
			want:    "https://example.com/user",
			wantErr: false,
		},

		{
			name: "Relative path URL (current no trailing slash)",
			args: args{
				parentURL:   "https://example.com/user/profile/settings",
				relativeURL: ".",
			},
			want:    "https://example.com/user/profile/settings",
			wantErr: false,
		},

		{
			name: "Same Parent and relative URL",
			args: args{
				parentURL:   "https://example.com/user/profile/settings",
				relativeURL: "https://example.com/user/profile/settings",
			},
			want:    "https://example.com/user/profile/settings",
			wantErr: false,
		},

		{
			name: "Different domain for Parent and relative URL",
			args: args{
				parentURL:   "https://example.com/user/profile/settings",
				relativeURL: "https://example101.com/user/profile/settings",
			},
			want:    "https://example101.com/user/profile/settings",
			wantErr: false,
		},

		{
			name: "Different sub-domain for Parent and relative URL",
			args: args{
				parentURL:   "https://example.com/user/profile/settings",
				relativeURL: "https://dev.example.com/user/profile/settings",
			},
			want:    "https://dev.example.com/user/profile/settings",
			wantErr: false,
		},

		{
			name: "Error: opaque parent URL (:)",
			args: args{
				parentURL:   "https:://example.com/user/profile/settings",
				relativeURL: ".",
			},
			want:    "",
			wantErr: true,
		},

		{
			name: "Malformed parent URL (/)",
			args: args{
				parentURL:   "https:///example.com/user/profile/settings",
				relativeURL: ".",
			},
			want:    "https://example.com/user/profile/settings",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := AbsoluteUrl(tt.args.parentURL, tt.args.relativeURL)
				if (err != nil) != tt.wantErr {
					t.Errorf("AbsoluteUrl() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				// Allow there to be trailing slashes in the output
				if got != tt.want && TrimSlash(got) != tt.want {
					t.Errorf("AbsoluteUrl() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestCleanSlash(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty string",
			args: args{path: ""},
			want: "",
		},

		{
			name: "No slash in string",
			args: args{path: "home"},
			want: "home",
		},

		{
			name: "Single leading slash in string",
			args: args{path: "/home"},
			want: "/home",
		},

		{
			name: "Single trailing slash in string",
			args: args{path: "home/"},
			want: "home/",
		},

		{
			name: "Mixed single slash in string",
			args: args{path: "/home/user/.local/share"},
			want: "/home/user/.local/share",
		},

		{
			name: "Mixed single slash in string (trailing slash)",
			args: args{path: "/home/user/.local/share/"},
			want: "/home/user/.local/share/",
		},

		{
			name: "Double slash leading",
			args: args{path: "//home/user/Downloads"},
			want: "/home/user/Downloads",
		},

		{
			name: "Double slash leading and trailing",
			args: args{path: "//home/user/Downloads//"},
			want: "/home/user/Downloads/",
		},

		{
			name: "Double slash leading, middle and trailing",
			args: args{path: "//home//user//Downloads//"},
			want: "/home/user/Downloads/",
		},

		{
			name: "Triple slash leading, middle and trailing",
			args: args{path: "///home///user///Downloads///"},
			want: "/home/user/Downloads/",
		},

		{
			name: "Mixed slash count; leading, middle and trailing",
			args: args{path: "////home///user//Downloads//////"},
			want: "/home/user/Downloads/",
		},

		{
			name: "Mixed slash count; leading, middle and trailing, variant 2",
			args: args{path: "////home///user//Downloads"},
			want: "/home/user/Downloads",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := CleanSlash(tt.args.path); got != tt.want {
					t.Errorf("CleanSlash() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestCleanUrl(t *testing.T) {
	type args struct {
		_url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Empty url",
			args:    args{_url: ""},
			want:    "",
			wantErr: false,
		},

		{
			name:    "Bad url",
			args:    args{_url: "http://example.com:xxxx"},
			want:    "",
			wantErr: true,
		},

		{
			name:    "Clean URL",
			args:    args{_url: "http://example.com/user/profile/settings"},
			want:    "http://example.com/user/profile/settings",
			wantErr: false,
		},

		{
			name:    "Clean URL path (trailing slash)",
			args:    args{_url: "http://example.com/user/profile/settings/"},
			want:    "http://example.com/user/profile/settings/",
			wantErr: false,
		},

		{
			name:    "Triple slashes in http URL scheme",
			args:    args{_url: "http:///example.com/user/profile/settings/"},
			want:    "http://example.com/user/profile/settings/",
			wantErr: false,
		},

		{
			name:    "Triple slashes in https URL scheme",
			args:    args{_url: "https:///example.com/user/profile/settings/"},
			want:    "https://example.com/user/profile/settings/",
			wantErr: false,
		},

		{
			name:    "Triple slashes in file URL scheme",
			args:    args{_url: "file:///example.com/user/profile/settings/"},
			want:    "file:///example.com/user/profile/settings/",
			wantErr: false,
		},

		{
			name:    "Four leading slashes in file URL scheme",
			args:    args{_url: "file:////example.com/user/profile/settings/"},
			want:    "file:///example.com/user/profile/settings/",
			wantErr: false,
		},

		{
			name:    "Only two leading slashes in file URL scheme",
			args:    args{_url: "file://example.com/user/profile/settings/"},
			want:    "file://example.com/user/profile/settings/",
			wantErr: false,
		},

		{
			name:    "Triple slashes in ftp URL scheme",
			args:    args{_url: "ftp:///www.example.com/user/profile/settings/"},
			want:    "ftp://www.example.com/user/profile/settings/",
			wantErr: false,
		},

		{
			name:    "More slashes in URL path",
			args:    args{_url: "http://www.example.com///user//profile/settings/////"},
			want:    "http://www.example.com/user/profile/settings/",
			wantErr: false,
		},

		{
			name:    "More slashes in URL scheme and path",
			args:    args{_url: "http:////www.example.com///user//profile/settings/////"},
			want:    "http://www.example.com/user/profile/settings/",
			wantErr: false,
		},

		{
			name:    "More slashes in file URL scheme and path",
			args:    args{_url: "file:////www.example.com///user//profile/settings/////"},
			want:    "file:///www.example.com/user/profile/settings/",
			wantErr: false,
		},

		{
			name:    "More slashes in file URL path (only two slashes in scheme)",
			args:    args{_url: "file://www.example.com///user//profile/settings/////"},
			want:    "file://www.example.com/user/profile/settings/",
			wantErr: false,
		},

		{
			name:    "Case study",
			args:    args{_url: "https://example.com///user//profile////settings//"},
			want:    "https://example.com/user/profile/settings/",
			wantErr: false,
		},

		{
			name:    "Case study 2",
			args:    args{_url: "/"},
			want:    "/",
			wantErr: false,
		},

		{
			name:    "Case study 3",
			args:    args{_url: "../../"},
			want:    "../../",
			wantErr: false,
		},

		{
			name:    "Case study 4",
			args:    args{_url: "..///..//"},
			want:    "../../",
			wantErr: false,
		},

		{
			name:    "Malformed URL",
			args:    args{_url: "https://example.com:xxxx/user/profile/settings"},
			want:    "",
			wantErr: true,
		},

		{
			name:    "Opaque URL",
			args:    args{_url: "https:://example.com/user/profile/settings"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := CleanUrl(tt.args._url)
				if (err != nil) != tt.wantErr {
					t.Errorf("CleanUrl() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("CleanUrl() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestDownloadFolder(t *testing.T) {
	type args struct {
		_url string
	}
	tests := []struct {
		name string
		args args
		want DownloadLocation
	}{
		{
			name: "Empty URL",
			args: args{_url: ""},
			want: DownloadLocation{
				FolderName: ".",
				FileName:   ".",
			},
		},

		{
			name: "Malformed URL",
			args: args{_url: "https://example.com:xxxxx/user/profile/settings/"},
			want: DownloadLocation{
				FolderName: "",
				FileName:   "",
			},
		},

		{
			name: "Domain root",
			args: args{_url: "https://example.com"},
			want: DownloadLocation{
				FolderName: "example.com",
				FileName:   ".",
			},
		},

		{
			name: "Domain root trailing slash",
			args: args{_url: "https://example.com/"},
			want: DownloadLocation{
				FolderName: "example.com",
				FileName:   ".",
			},
		},

		{
			name: "Domain root index",
			args: args{_url: "https://example.com/index.html"},
			want: DownloadLocation{
				FolderName: "example.com",
				FileName:   "index.html",
			},
		},

		{
			name: "Path 1",
			args: args{_url: "https://example.com/user/profile/settings"},
			want: DownloadLocation{
				FolderName: "example.com/user/profile",
				FileName:   "settings",
			},
		},

		{
			name: "Path 1 trailing slash",
			args: args{_url: "https://example.com/user/profile/settings/"},
			want: DownloadLocation{
				FolderName: "example.com/user/profile/settings",
				FileName:   ".",
			},
		},

		{
			name: "Path 2 trailing slash",
			args: args{_url: "https://theuselessweb.com/privacy"},
			want: DownloadLocation{
				FolderName: "theuselessweb.com",
				FileName:   "privacy",
			},
		},

		{
			name: "Root url (no scheme)",
			args: args{_url: "/example.com/user/profile/settings"},
			want: DownloadLocation{
				FolderName: "/example.com/user/profile",
				FileName:   "settings",
			},
		},

		{
			name: "Relative url (no scheme)",
			args: args{_url: "example.com/user/profile/settings"},
			want: DownloadLocation{
				FolderName: "example.com/user/profile",
				FileName:   "settings",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := DownloadFolder(tt.args._url); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("DownloadFolder() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestRelativeFolder(t *testing.T) {
	type args struct {
		parentFolder string
		rpath        string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// Example Case:
		{
			name: "BasicRelativePath",
			args: args{
				parentFolder: "/home/user/Documents",
				rpath:        "/home/user/Documents/Projects/GoProject",
			},
			want: "Projects/GoProject",
		},
		// Edge Cases:
		{
			name: "SameFolder",
			args: args{
				parentFolder: "/home/user",
				rpath:        "/home/user",
			},
			want: ".", // Relative path to the same folder is "."
		},
		{
			name: "ParentFolderIsRoot",
			args: args{
				parentFolder: "/",
				rpath:        "/home/user/Documents",
			},
			want: "home/user/Documents",
		},
		{
			name: "RpathIsRoot",
			args: args{
				parentFolder: "/home/user/Documents",
				rpath:        "/",
			},
			want: "/home/user/Documents",
		},

		{
			name: "Different Folders",
			args: args{
				parentFolder: "/home/username/Download",
				rpath:        "/home/username2/Download",
			},
			want: "../../username2/Download",
		},

		{
			name: "Different relative folders",
			args: args{
				parentFolder: "username/Download",
				rpath:        "username2/Download",
			},
			want: "../../username2/Download",
		},
		{
			name: "EmptyParent",
			args: args{
				parentFolder: "", // Empty parent, should default to current directory
				rpath:        "/home/user/Documents",
			},
			want: "/home/user/Documents", // Should return the absolute path of rpath
		},
		{
			name: "EmptyRpath",
			args: args{
				parentFolder: "/home/user/Documents",
				rpath:        "", // Empty rpath, should default to current directory
			},
			want: ".", // Should return "." (current directory relative to itself)
		},
		{
			name: "TrailingSlashParent",
			args: args{
				parentFolder: "/home/user/Documents/",
				rpath:        "/home/user/Documents/Projects",
			},
			want: "Projects", // Trailing slash shouldn't affect the result
		},
		{
			name: "TrailingSlashRpath",
			args: args{
				parentFolder: "/home/user/Documents",
				rpath:        "/home/user/Documents/Projects/",
			},
			want: "Projects", // Trailing slash shouldn't affect the result
		},
		{
			name: "DifferentDriveWindows",
			args: args{
				parentFolder: "C:/Users/username/Docs",
				rpath:        "D:/Projects/Go",
			},
			want: "../../../../D:/Projects/Go", // Returns absolute path for different drives
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := RelativeFolder(tt.args.parentFolder, tt.args.rpath); got != tt.want {
					t.Errorf("RelativeFolder() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func ExampleRelativeFolder() {
	parent := "/Users/username/Documents"
	child := "/Users/username/Documents/Projects/GoProject"

	// Get the relative path
	relativePath := RelativeFolder(parent, child)

	// Print the cleaned relative path:
	fmt.Println(filepath.Clean(relativePath))

	// Expected Output:
	// Projects/GoProject
}

func TestSameHost(t *testing.T) {
	type args struct {
		url1 string
		url2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Same host with same protocol",
			args: args{
				url1: "https://example.com/path",
				url2: "https://example.com/anotherpath",
			},
			want: true,
		},
		{
			name: "Same host with different protocols",
			args: args{
				url1: "http://example.com/path",
				url2: "https://example.com/anotherpath",
			},
			want: true,
		},
		{
			name: "Different hosts",
			args: args{
				url1: "https://example.com/path",
				url2: "https://another.com/anotherpath",
			},
			want: false,
		},
		{
			name: "Subdomain vs main domain",
			args: args{
				url1: "https://sub.example.com/path",
				url2: "https://example.com/anotherpath",
			},
			want: false,
		},
		{
			name: "Same host with different ports",
			args: args{
				url1: "https://example.com:8080/path",
				url2: "https://example.com:9090/anotherpath",
			},
			want: false,
		},
		{
			name: "Non-standard scheme",
			args: args{
				url1: "htp://example.com/path", // intentional typo in scheme
				url2: "https://example.com/anotherpath",
			},
			want: true,
		},
		{
			name: "Empty URLs",
			args: args{
				url1: "",
				url2: "https://example.com/anotherpath",
			},
			want: false,
		},
		{
			name: "Both URLs empty",
			args: args{
				url1: "",
				url2: "",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := SameHost(tt.args.url1, tt.args.url2)
				if got != tt.want {
					t.Errorf("SameHost() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func ExampleSameHost() {
	isSameHost := SameHost("https://example.com/path", "https://example.com/anotherpath")
	if isSameHost {
		fmt.Println("The URLs have the same host.")
	} else {
		fmt.Println("The URLs have different hosts.")
	}
	// Output: The URLs have the same host.
}
