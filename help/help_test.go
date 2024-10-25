package help

import (
	"strings"
	"testing"
)

func TestPrintManPage(t *testing.T) {
	expected := strings.Clone(PrintManPage())
	if got := PrintManPage(); got != expected {
		t.Errorf("PrintManPage() got: \n%s\nwant: %s\n", got, expected)
	}
}
