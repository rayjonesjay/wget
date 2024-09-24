package args

import (
	"testing"
)

// TestIsHelpFlag is a test file for the IsHelpFlag function
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
