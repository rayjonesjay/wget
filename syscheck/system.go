package syscheck

import (
	"errors"
	"runtime"
)

// CheckOperatingSystem checks if the underlying operating system is neither linux nor macOS
func CheckOperatingSystem() error {
	operatingSystem := runtime.GOOS
	if operatingSystem != "linux" && operatingSystem != "darwin" {
		return errors.New("program cannot run on non-unix operating system")
	}
	return nil
}
