package info

import (
	"strings"
	"testing"
)

func TestVersionText(t *testing.T) {
	v := VersionText()
	if !strings.Contains(v, "Zone01") || !strings.Contains(v, "wget") {
		t.Errorf("Expected this program to be wget, and the property of Zone 01 Kisumu")
	}
}
