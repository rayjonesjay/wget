// Package globals contains various one-shot utilities,
// that could have otherwise been put in different packages; we instead mix them all here
package globals

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
)

var (
	// SrcElements defines html elements that typically define their linked resource in the "src" attribute
	SrcElements = map[string]bool{
		"video": true, "audio": true, "img": true, "script": true, "iframe": true,
	}
	// DataElements defines html elements that typically define their linked resource in the "data" attribute
	DataElements = map[string]bool{"object": true}
	// HrefElements defines html elements that typically define their linked resource in the "href" attribute
	HrefElements = map[string]bool{"a": true, "link": true}
	// SrcDataElements defines html elements that typically define their linked resource in either the "src" or "data"
	//attribute
	SrcDataElements = MergeMaps(SrcElements, DataElements)
	// AllResourceElements defines html elements that typically define their linked resource in
	//either the "src", "href", or "data" attribute
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
