package syscheck

import (
	"os"
	"testing"
	"wget/temp"
)

func TestCheckOperatingSystem(t *testing.T) {
	tests := []struct {
		name            string
		operatingSystem string
		wantErr         bool
	}{
		{"Linux OS", "linux", false},
		{"macOS OS", "darwin", false},
		{"Windows OS", "windows", true},
		{"Other OS", "freebsd", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckOperatingSystem(tt.operatingSystem)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckOperatingSystem(%q) error = %v, wantErr %v", tt.operatingSystem, err, tt.wantErr)
			}
		})
	}
}

func TestGetTerminalWidth(t *testing.T) {
	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout
	}()

	var err error
	os.Stdout, err = temp.File()
	if err != nil {
		t.Fatal(err)
	}

	width := GetTerminalWidth()
	if width != 300 {
		t.Errorf("GetTerminalWidth = %d, want %d", width, 300)
	}
}
