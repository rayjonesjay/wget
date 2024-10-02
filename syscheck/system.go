package syscheck

import (
	"errors"
)

// CheckOperatingSystem checks if the underlying operating system is neither Linux nor macOS
// allows passing an OS name (used for testing)
func CheckOperatingSystem(operatingSystem string) error {
	if operatingSystem != "linux" && operatingSystem != "darwin" {
		return errors.New("program cannot run on non-unix operating system")
	}
	return nil
}