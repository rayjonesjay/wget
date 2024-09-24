package args

import (
	"testing"
)

// TestIsHelpFlag is a test function for the IsHelpFlag function
func TestIsHelpFlag(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"test1", "--help", true},
		{"test2", "-help", false},
		{"test3", "-h", false},
		{"test4", "help", false},
		{"test5", "h", false},
		{"test6", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHelpFlag(tt.input); got != tt.want {
				t.Errorf("IsHelpFlag() = %v, want %v", got, tt.want)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := IsPathFlag(tt.input)
			if got1 != tt.want1 || got2 != tt.want2 {
				t.Errorf("IsPathFlag() = [%v %v], want [%v %v]", got1, got2, tt.want1, tt.want2)
			}
		})
	}
}
