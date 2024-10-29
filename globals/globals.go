// Package globals contains various one-shot utilities,
// that could have otherwise been put in different packages; we instead mix them all here
package globals

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"wget/httpx"
	"wget/syscheck"

	"golang.org/x/net/html"
)

var (
	// SrcElements defines html elements that typically define their linked resource in the "src" attribute
	SrcElements = map[string]bool{"img": true}
	// DataElements defines html elements that typically define their linked resource in the "data" attribute
	DataElements = map[string]bool{}
	// HrefElements defines html elements that typically define their linked resource in the "href" attribute
	HrefElements = map[string]bool{"a": true, "link": true}
	// SrcDataElements defines html elements that typically define their linked resource in either the "src" or "data"
	// attribute
	SrcDataElements = MergeMaps(SrcElements, DataElements)
	// AllResourceElements defines html elements that typically define their linked resource in
	// either the "src", "href", or "data" attribute
	AllResourceElements = MergeMaps(SrcDataElements, HrefElements)
)

// MergeMaps creates a new map, whose keys and values includes merged keys and
// values from the given maps. If a key exists in both maps, then the value in
// the second map takes precedence
func MergeMaps(a, b map[string]bool) map[string]bool {
	m := make(map[string]bool, len(a)+len(b))
	for k, v := range a {
		m[k] = v
	}
	for k, v := range b {
		m[k] = v
	}

	return m
}

// RenderToString renders the given html node to a string; returns the string the
// html node was rendered to. If a nil html.Node is supplied, then an empty string would be returned, otherwise,
// the returned string will always be a valid HTML based on the given html.Node.
func RenderToString(node *html.Node) string {
	if node == nil {
		return ""
	}

	var buffer bytes.Buffer

	// Create an io.Writer that writes to the buffer
	writer := io.Writer(&buffer)

	_ = html.Render(writer, node)

	return buffer.String()
}

// FormatSize helper function to format byte size
func FormatSize(size int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GiB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MiB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KiB", float64(size)/KB)
	default:
		if size < 0 {
			return "--.- B"
		}
		return fmt.Sprintf("%d B", size)
	}
}

// RoundBytes rounds the given bytes to the nearest MB or GB.
func RoundBytes(bytes int64) string {
	if bytes < httpx.GB {
		return fmt.Sprintf("%.2f%s", float64(bytes)/float64(httpx.MB), "MB")
	} else {
		return fmt.Sprintf("%.2f%s", float64(bytes)/float64(httpx.GB), "GB")
	}
}

// printLinesMutex is a global mutex locker for PrintLines
var printLinesMutex = &sync.Mutex{}

// PrintLines prints multiple lines of text starting from a specific row.
func PrintLines(baseRow int, lines []string) {
	printLinesMutex.Lock()
	defer printLinesMutex.Unlock()
	for i, line := range lines {
		i++
		syscheck.MoveCursor(baseRow + i) // move to the correct line
		fmt.Print("\033[K")              // clear the line
		fmt.Print(line)
	}
}

// StringTimes returns an array containing n instances of the given string
func StringTimes(s string, n int) (out []string) {
	for i := 0; i < n; i++ {
		out = append(out, s)
	}
	return
}
