package syscheck

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	// HideCursor hides the terminal cursor to avoid it from blinking
	HideCursor = from("\033[?25l")
	// ShowCursor makes the cursor visible again
	ShowCursor = from("\033[?25h")
	// ClearScreen clears the terminal screen.
	ClearScreen = from("\033[2J")
	// MoveCursor moves the terminal cursor to the specified row, we don't need columns here
	MoveCursor func(row int)
)

func init() {
	MoveCursor = fromArg("\033[%d;0H")
}

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
		return 120
	}
	return width
}

func from(format string) func() {
	return func() {
		fmt.Printf(format)
	}
}

func fromArg(format string) func(int) {
	return func(arg int) {
		fmt.Printf(format, arg)
	}
}
