package syscheck

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
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
	Row     uint16
	Col     uint16
	XPixels uint16
	YPixels uint16
}

// GetTerminalWidth retrieves the width of the terminal at runtime, to be
// utilized for example in calculating the width of the progress bar to the
// terminal
//
//	width must be at least twice the progress bar.
//
// syscall.Syscall uses SYS_IOCTL syscall to retrieve terminal dimensions by interacting with stdin (file descriptor 0)
// TIOCGWINSZ is a Unix ioctl (input output control) command that gets the terminal's rows and columns.
// terminal holds the dimensions of the terminal. but we are only interested in the Col field, which stores the column
// if the syscall fails we return an error
func GetTerminalWidth() int {
	fd := os.Stdout.Fd()
	ws := &terminal{}
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(ws)))
	width := int(ws.Col)
	if width <= 0 {
		return 300
	}
	return width
}

// MoveCursor moves the terminal cursor to the specified row, we don't need columns here
func MoveCursor(row int) {
	// 0H is the column part, all printing is done from left
	fmt.Printf("\033[%d;0H", row)
}

// HideCursor hides the terminal cursor to avoid it from blinking
func HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor makes the cursor visible again
func ShowCursor() {
	fmt.Print("\033[?25h")
}

// ClearScreen clears the terminal screen.
func ClearScreen() {
	fmt.Print("\033[2J")
}
