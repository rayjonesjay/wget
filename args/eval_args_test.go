package args

import (
	"errors"
	"reflect"
	"testing"
	"wget/ctx"
)

// TestToBytes is a test function for ToBytes function
func TestToBytes(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  int64
	}{
		{"test1", "20k", 20_000},
		{"test2", "20M", 20_000_000},
		{"test3", "2", 2},
		{"test4", "1KB", 0},
		{"test5", "k", 0},
		{"test6", "M", 0},
		{"test7", "", 0},
		{"test8", "200000000000M", 200_000_000_000_000_000},
		{"test9", "1000k", 1_000_000},
		{"test10", "1.0", 1},
		{"test11", "14.5", 14},
		{"test12", "15.5k", 15_000},
		{"test13", "123.56M", 123_000_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToBytes(tt.input)
			if got != tt.want {
				t.Errorf("ToBytes(%q) got %d want %d", tt.input, got, tt.want)
			}
		})
	}

}

// TestIsPathFlag is a test function for IsPathFlag function
func TestIsPathFlag(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want1 bool
		want2 string
	}{
		{"test1", "-P=/home/Downloads", true, "/home/Downloads"},
		{"test2", "-P=/home", true, "/home"},
		{"test3", "-P=", false, ""},
		{"test4", "-P=...", true, "..."},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got1, got2 := IsPathFlag(tt.input)
				if got1 != tt.want1 || got2 != tt.want2 {
					t.Errorf("IsPathFlag() = [%v %v], want [%v %v]", got1, got2, tt.want1, tt.want2)
				}
			},
		)
	}
}

// TestInputFile is a test function for InputFile function
func TestInputFile(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want1   bool
		want2   string
		wantErr error
	}{
		{"test1", "-i=file.txt", true, "file.txt", nil},
		{"test2", "-i=file", true, "file", nil},
		{"test3", "-i", false, "", errors.New("path might be empty")},
		{"test4", "-i=.", false, "", errors.New("path is a directory")},
		{"test5", "-i=..", false, "", errors.New("path is a directory")},
		{"test6", "-i=...", true, "...", nil},
		{"test7", "-i=/", false, "", errors.New("path is a directory")},
		{"test8", "-i=/////", false, "", errors.New("path is a directory")},
		{"test9", "-i=//", false, "", errors.New("path is a directory")},
		{"test10", `-i=\`, true, `\`, nil},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got1, got2, gotError := InputFile(tt.input)

				if gotError != nil {
					if gotError.Error() != tt.wantErr.Error() {
						t.Errorf("InputFile() => %v | %v ", gotError, tt.wantErr)
					}
				}

				if got1 != tt.want1 || got2 != tt.want2 {
					t.Errorf("InputFile() = [%v %v] , want [%v %v]", got1, got2, tt.want1, tt.want2)
				}
			},
		)
	}
}

// TestIsOutputFlag is a test function for IsOutputFlag function
func TestIsOutputFlag(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want1 bool
		want2 string
	}{
		{"test1", "-O=file.txt", true, "file.txt"},
		{"test2", "-O=file", true, "file"},
		{"test3", "-O=", false, ""},
		{"test4", "-o=file.txt", false, ""},
		{"test5", "-O=/", false, ""},
		{"test6", "-O=.", false, ""},
		{"test7", "-O=..", false, ""},
		{"test8", "-O=...", true, "..."},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got1, got2 := IsOutputFlag(tt.input)
				if got1 != tt.want1 || got2 != tt.want2 {
					t.Errorf("IsOutputFlag() = got [%v %v], want [%v %v]", got1, got2, tt.want1, tt.want2)
				}
			},
		)
	}
}

func TestEvalArgs(t *testing.T) {
	type args struct {
		arguments []string
	}

	// Initialize arguments map for test cases
	mappy := map[string][]string{
		"Omega": {},
		"Beta":  {"-O=file.txt", "-i=urls.txt", "-P=/home/Downloads"},
		"Alpha": {"https://learn.zone01kisumu.ke/git/root/public/raw/branch/master/subjects/ascii-art/shadow.txt"},
	}

	tests := []struct {
		name          string
		args          args
		wantArguments ctx.Context
	}{
		// when no arguments have been parsed
		{"Omega", args{arguments: mappy["Omega"]}, ctx.Context{}},

		{
			"Beta", args{arguments: mappy["Beta"]},
			ctx.Context{OutputFile: "file.txt", InputFile: "urls.txt", SavePath: "./"}, // if SavePath cannot be accessed(permissions) default to current directory
		},

		{
			name: "Alpha", args: args{arguments: mappy["Alpha"]},
			wantArguments: ctx.Context{
				Links: []string{
					"https://learn.zone01kisumu.ke" +
						"/git/root/public/raw/branch/master/subjects/ascii-art/shadow.txt",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if gotArguments := DownloadContext(tt.args.arguments); !reflect.DeepEqual(
					gotArguments, tt.wantArguments,
				) {
					s := "-------------------------------------------------------------"
					t.Errorf("DownloadContext() %s \ngot %v\nwant %v \n%s", s, gotArguments, tt.wantArguments, s)
				}
			},
		)
	}
}
