package terminal

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

// Progress displays text in the terminal in real-time, for reporting the progress of a background activity
type Progress struct {
	mu                sync.RWMutex // Mutex for thread-safe access to lines
	terminal          *Terminal    // where the progress will be written, must not be nil
	termWidth         int
	termHeight        int
	shouldUpdate      chan struct{}
	currentLine       int
	lines             []string              // The lines displayed on the screen
	activeAnimations  map[int]chan struct{} // Tracks active animations per line
	close             chan struct{}
	run               bool
	exit              chan struct{}
	runExit           chan struct{}
	alwaysScrollToEnd bool
}

// New returns a new Progress instance
func New(term *Terminal, exit chan struct{}) *Progress {
	if term == nil {
		panic("terminal is nil")
	}

	// register and update the terminal size
	p := Progress{
		alwaysScrollToEnd: true,
		exit:              exit,
		runExit:           make(chan struct{}),
	}
	p.mu = sync.RWMutex{}
	p.terminal = term
	p.termWidth, p.termHeight = Size()
	p.shouldUpdate = make(chan struct{}, 1)
	p.currentLine = 0
	p.lines = make([]string, 0, p.termHeight)
	p.activeAnimations = make(map[int]chan struct{})
	p.close = make(chan struct{}, 1)

	// Listen for changes in the terminal size, and update accordingly
	go func() {
		// Create a channel to receive signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGWINCH)

		for {
			select {
			case <-sigChan:
				log.Println("[WATCHER] Terminal size changed: updating...")
				p.termWidth, p.termHeight = Size()
				p.update()
			case <-p.close:
				log.Println("[WATCHER] -> closing...")
				return
			}
		}
	}()

	return &p
}

// Run is a no return function, thus should be called in another go routine to process progress.
// Note: This function should also only be called once
func (p *Progress) Run() {
	p.mu.Lock()
	log.Println("[RUN] starting...")
	if p.run {
		p.mu.Unlock()
		p.Exit()
		fmt.Println("error: calling `Run()` more than once on the same `terminal.Progress` instance")
		return
	}
	p.run = true
	p.mu.Unlock()

	// Checks succeeded, register a termination handler
	// Channel to listen for OS signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		// Block until a signal is received
		sig := <-sigs
		log.Printf("Received signal: %s\n", sig)
		p.Exit()
	}()

	log.Println("[RUN] starting input handler...")
	go func() {
		err := p.handleInput()
		if err != nil {
			log.Println("Error handling input:", err)
		}
	}()

	// Run indefinitely, processing terminal display on-demand
	log.Println("[RUN] Main render loop...")
	for {
		select {
		case <-p.runExit:
			return
		case <-p.shouldUpdate:
			p.render()
		}
	}
}

// ensureSize checks and expands the screen's size to include the specified line index.
// It grows the screen by adding blank lines if the line index exceeds the current number of lines.
func (p *Progress) ensureSize(line int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if line >= len(p.lines) {
		newLines := make([]string, line-len(p.lines)+1)
		p.lines = append(p.lines, newLines...)
		if p.alwaysScrollToEnd {
			currentLine := line - p.termHeight + 1
			if currentLine >= 0 && currentLine < len(p.lines) {
				p.currentLine = currentLine
			}
		}
	}
}

// Clear erases all lines on the screen by setting them to empty strings.
// It invokes the line update callback for each line cleared.
func (p *Progress) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lines = make([]string, 0, p.termHeight)
	p.update()
}

// SetLine sets the text of a specific line on the screen, expanding the screen size if necessary.
// The callback function is called to indicate the line update.
func (p *Progress) SetLine(line int, text string) {
	p.ensureSize(line)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lines[line] = text
	p.update()
}

// SetTextAt sets text at a specific line and column position, expanding the screen size if necessary.
// The screen line is adjusted to include the text at the specified position.
// It invokes the line update callback after setting the text.
func (p *Progress) SetTextAt(line, col int, text string) {
	p.ensureSize(line)
	p.mu.Lock()
	defer p.mu.Unlock()
	if col < 0 {
		return
	}
	lineContent := p.lines[line]
	if col < len(lineContent) {
		p.lines[line] = lineContent[:col] + text + lineContent[col+len(text):]
	} else {
		// Extend line if necessary
		p.lines[line] = lineContent + strings.Repeat(" ", col-len(lineContent)) + text
	}
	p.update()
}

// Exit the indefinite Run function
func (p *Progress) Exit() {
	_ = p.terminal.EndWin()

	select {
	case p.close <- struct{}{}:
	default:
	}

	select {
	case p.exit <- struct{}{}:
	default:
	}

	select {
	case p.runExit <- struct{}{}:
	default:
	}
}

func (p *Progress) handleInput() error {
	buffer := make([]byte, 3)
	for {
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			return err
		}
		isEnd := false
		if n == 1 {
			switch buffer[0] {
			case 'q':
				p.Exit()
				return nil
			case 'j': // Down arrow (first byte of escape sequence)
				p.scrollDown(1)
			case 'k': // Up
				p.scrollUp(1)
			}
		} else if n == 3 && buffer[0] == '\x1b' && buffer[1] == '[' {
			switch buffer[2] {
			case 'A': // Up arrow
				p.scrollUp(1)
			case 'B': // Down arrow
				p.scrollDown(1)
			case '5': // Page Up
				p.scrollUp(p.termHeight)
			case '6': // Page Down
				p.scrollDown(p.termHeight)
			case 0x48: // Home
				p.scrollHome()
			case 0x46: // End
				isEnd = true
				p.scrollEnd()
			}
		}

		if !isEnd {
			p.alwaysScrollToEnd = false
		}
	}
}

func (p *Progress) scrollUp(lines int) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	p.currentLine -= lines
	if p.currentLine < 0 {
		p.currentLine = 0
	}
	p.update()
}

func (p *Progress) scrollHome() {
	p.mu.RLock()
	defer p.mu.RUnlock()
	p.currentLine = 0
	p.update()
}

func (p *Progress) scrollEnd() {
	p.mu.Lock()
	p.alwaysScrollToEnd = true
	p.mu.Unlock()

	p.mu.RLock()
	defer p.mu.RUnlock()
	currentLine := len(p.lines) - p.termHeight
	if currentLine >= 0 && currentLine < len(p.lines) {
		p.currentLine = currentLine
	}
	p.update()
}

func (p *Progress) scrollDown(lines int) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	maxLine := len(p.lines) - p.termHeight
	if maxLine < 0 {
		maxLine = 0
	}
	p.currentLine += lines
	if p.currentLine > maxLine {
		p.currentLine = maxLine
	}
	p.update()
}

func (p *Progress) render() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Clear screen and move cursor to top
	p.terminal.Clear()

	// Display content
	endLine := p.currentLine + p.termHeight
	if endLine > len(p.lines) {
		endLine = len(p.lines)
	}

	for i, line := range p.lines[p.currentLine:endLine] {
		p.terminal.Print(i, line)
	}

	// Display status line
	fmt.Printf("\x1b[%d;0H\x1b[7m", p.termHeight+1)
	status := fmt.Sprintf(
		"Output Line %d-%d/%d (%.0f%%) - (press q to quit)",
		p.currentLine+1,
		endLine,
		len(p.lines),
		float64(endLine)/float64(len(p.lines))*100,
	)
	fmt.Printf("%-*s\x1b[0m", p.termWidth, status)
}

func (p *Progress) update() {
	select {
	case p.shouldUpdate <- struct{}{}:
	default:
	}
}
