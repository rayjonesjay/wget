package syscheck

import (
	"fmt"
	"testing"
	"time"
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
func TestGetCurrentTime(t *testing.T) {

	tests := []struct {
		input bool
		want  string
	}{
		{true, fmt.Sprintf("start at %s", time.Now().Format("2006-01-02 15:04:05"))},
		{false, fmt.Sprintf("finished at %s", time.Now().Format("2006-01-02 15:04:05"))},
	}

	for _, tt := range tests {
		got := GetCurrentTime(tt.input)
		if got != tt.want {
			t.Errorf("GetCurrentTime() Failed got %s want %s", got, tt.want)
		}
	}
}
