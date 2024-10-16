package syscheck

import (
	"errors"
	"syscall"
	"unsafe"
	"wget/xerr"
)

// CheckOperatingSystem checks if the underlying operating system is neither Linux nor macOS
// allows passing an OS name (used for testing)
func CheckOperatingSystem(operatingSystem string) error {
	if operatingSystem != "linux" && operatingSystem != "darwin" {
		return errors.New("program cannot run on non-unix operating system")
	}
	return nil
}

// terminal represents the structure to store terminal dimensions
type terminal struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// GetTerminalWidth retrieves the width of the terminal at runtime, this this is because in order to print the progress bar to the terminal
// the terminal width must be at least twice the progress bar.
// syscall.Syscall uses SYS_IOCTL syscall to retrieve terminal dimensions by interacting with stdin (file descriptor 0)
// TIOCGWINSZ is a Unix ioctl (input output control) command that gets the terminal's rows and columns.
// terminal holds the dimensions of the terminal. but we are only interested in the Col field, which stores the column
// if the syscall fails we return an error
func GetTerminalWidth() int {
	// syscall to interact with the terminal using file descriptor 0 (stdin)
	ws := &terminal{}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(0), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(ws)))
	if err != 0 {
		xerr.WriteError("cannot get terminal size, make sure you are using darwin or linux os", 1, true)
	}
	width := int(ws.Col)
	if width < 65 {
		xerr.WriteError("terminal size to small adjust width to at least 65", 1, true)
	}
	return int(width)
}
