package help

import (
	"fmt"
	"testing"
	"wget/temp"
)

func TestPrintManPage(t *testing.T) {
	// Simply print the help manual to a temporary file
	file, err := temp.File()
	if err != nil {
		t.Fatal(err)
	}

	_, err = file.WriteString(PrintManPage())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Man page written to file %q\n", file.Name())
}
