package syscheck

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
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
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(ws)))
	if err != 0 {
		xerr.WriteError("cannot get terminal size, make sure you are using darwin or linux os", 1, true)
	}
	width := int(ws.Col)
	if width < 65 {
		xerr.WriteError("terminal size to small adjust width to at least 65", 1, true)
	}
	return int(width)
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

// GetCurrentTime prints a textual representation of the current time, it takes isStart
// if isStart is true it indicates the download process has started, if false means end of the current download
func GetCurrentTime(isStart bool) string {
	// Get the current time
	currentTime := time.Now()

	// Format time to print up to seconds
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	if isStart {
		return fmt.Sprintf("start at %s", formattedTime)
	}
	return fmt.Sprintf("finished at %s", formattedTime)
}
