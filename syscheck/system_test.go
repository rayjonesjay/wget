package syscheck

import (
	"testing"
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
