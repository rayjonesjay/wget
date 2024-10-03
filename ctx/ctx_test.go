package ctx

import (
	"fmt"
	"testing"
)

type Arg struct {
	*Context
}

// Receiver function to set the OutputFile of the embedded Context
func (a *Arg) SetOutputFile(filename string) {
	a.OutputFile = filename // Accessing and modifying Context's field directly
}

// Receiver function to check if mirroring is enabled
func (a *Arg) IsMirroringEnabled() bool {
	return a.Mirror // Accessing Context's field directly
}

func TestContext(t *testing.T) {
	arg := Arg{
		Context: &Context{}, // Initialize the embedded Context
	}

	// Using receiver functions
	arg.SetOutputFile("output.txt")
	fmt.Println("Output file:", arg.OutputFile) // Accessing Context's field directly

	if arg.OutputFile != "output.txt" {
		t.Errorf("Output file not equal to output.txt")
	}

	arg.Mirror = true
	fmt.Println("Mirroring enabled:", arg.IsMirroringEnabled()) // Using a receiver function

	if !arg.IsMirroringEnabled() {
		t.Errorf("IsMirroring not enabled yet it was just set to true")
	}
}
