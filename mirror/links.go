package mirror

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"wget/mirror/links"

	"golang.org/x/net/html"
)

var (
	srcElements = map[string]bool{
		"video": true, "audio": true, "link": true, "img": true, "script": true, "iframe": true,
	}
	dataElements    = map[string]bool{"object": true}
	hrefElements    = map[string]bool{"a": true}
	srcDataElements = mergeMaps(srcElements, dataElements)
)

// UrlExtract structures some Url with the target html attribute from where the target url was extracted
type UrlExtract struct {
	Url  string
	Attr *html.Attribute
}

// Extract links to external content from a html node. Returns a list of all extracted links, with a pointer
// to the target node attribute's value; this allows you to change the target
// source of the external content in the target html node
func Extract(node *html.Node) (extractedLinks []UrlExtract) {
	//extractedLinks = make([]UrlExtract, 0, 10)
	extract(node, &extractedLinks)
	return
}

// extract recursively traverses through the node,
// to extract all links in the HTML tree; the links are saved in the input UrlExtract array
// pointed to by `extractedLinks`
func extract(n *html.Node, extractedLinks *[]UrlExtract) {
	//fmt.Printf("type: %v data: %v\n", n.Type, n.Data)
	if n.Type == html.ElementNode && srcDataElements[n.Data] {
		for i, a := range n.Attr {
			if (a.Key == "src" && srcElements[n.Data]) ||
				(a.Key == "data" && dataElements[n.Data]) ||
				(a.Key == "href" && hrefElements[n.Data]) {
				*extractedLinks = append(*extractedLinks, UrlExtract{a.Val, &n.Attr[i]})
				break
			}
		}
	}

	if n.Type == html.ElementNode && n.Data == "style" {
		// the current node is a <style> tag, extract its content
		styleContent := textContent(n)
		cssAssets := links.FromCssUrl(styleContent)
		*extractedLinks = append(*extractedLinks, fromLinks(cssAssets)...)
	}

	// done extracting links from the current html node,
	//now, recursively, extract links from all child nodes to the current node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extract(c, extractedLinks)
	}

	return
}

// textContent returns the text content of the target html node.
func textContent(n *html.Node) string {
	var content strings.Builder
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			content.WriteString(child.Data)
		}
	}

	return content.String()
}

// fromLinks creates an array of UrlExtract from the given list of URL links
func fromLinks(links []string) (out []UrlExtract) {
	for _, l := range links {
		out = append(out, UrlExtract{l, nil})
	}
	return
}

// ExtractFirst returns the first link extracted from the node by Extract, or an
// error if the first link cannot be extracted
func ExtractFirst(node *html.Node) (UrlExtract, error) {
	extracts := Extract(node)
	if len(extracts) == 0 {
		return UrlExtract{}, errors.New("no elements found")
	}

	return extracts[0], nil
}

// RenderToString renders the given html node to a string; returns the string the
// html node was rendered to, or an error if the render was unsuccessful
func RenderToString(node *html.Node) (string, error) {
	var buffer bytes.Buffer

	// Create an io.Writer that writes to the buffer
	writer := io.Writer(&buffer)

	err := html.Render(writer, node)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// mergeMaps creates a new map, whose keys and values includes merged keys and
// values from the given maps. If a key exists in both maps, then the value in
// the second map takes precedence
func mergeMaps(a, b map[string]bool) map[string]bool {
	m := make(map[string]bool, len(a)+len(b))
	for k, v := range a {
		m[k] = v
	}
	for k, v := range b {
		m[k] = v
	}

	return m
}
