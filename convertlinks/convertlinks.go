// Package convertlinks defines some utilities to handle the --convert-links options
package convertlinks

import (
	"golang.org/x/net/html"
	"strings"
	"wget/css"
	"wget/globals"
)

var (
	srcElements         = globals.SrcElements
	dataElements        = globals.DataElements
	hrefElements        = globals.HrefElements
	allResourceElements = globals.AllResourceElements
)

// OfCss converts all linked urls in the given CSS string to links defined by the given transformer
func OfCss(inputCss string, transformer func(url string) string) string {
	return css.TransformCssUrl(inputCss, transformer)
}

// OfHtml converts all linked urls in the given html node to links defined by the given transformer.
// If either the html node or the transformer function is nil, then it's a no-op situation
func OfHtml(n *html.Node, transformer func(url string, isA bool) string) {
	if n == nil || transformer == nil {
		return
	}
	e := ExtractConfig{
		transformer: transformer,
		cssTransformer: func(url string) string {
			return transformer(url, false)
		},
	}
	e.extract(n)
}

// ExtractConfig is a wrapper struct to hold state data for all the recursive calls of extract
type ExtractConfig struct {
	// transformer is a function that takes the target url, then converts it to a local filesystem URI,
	//when the --convert-links option is enabled.
	transformer    func(url string, isA bool) string
	cssTransformer func(url string) string
}

// extract recursively traverses through the given html node,
// to extract all links in the HTML tree
func (e *ExtractConfig) extract(n *html.Node) {
	//fmt.Printf("type: %v data: %v\n", n.Type, n.Data)
	if n.Type == html.ElementNode && allResourceElements[n.Data] {
		for i, a := range n.Attr {
			if (a.Key == "src" && srcElements[n.Data]) ||
				(a.Key == "data" && dataElements[n.Data]) ||
				(a.Key == "href" && hrefElements[n.Data]) {
				// n.Data contains a URL such as https://example.com/path/to/image.png
				isA := n.Data == "a"
				n.Attr[i].Val = e.transformer(n.Attr[i].Val, isA)
				break
			}
		}
	}

	if n.Type == html.ElementNode && n.Data == "style" {
		// the current node is a <style> tag, extract its text content (the CSS definitions)
		styleContent, tcNode := textContentNode(n)
		// the style content was extracted from `tcNode.Data`, replace with the converted links
		tcNode.Data = OfCss(styleContent, e.cssTransformer)
	}

	if n.Type == html.ElementNode {
		// for every html element, check the style attribute and make changes to the linked urls in such inline CSS
		for i, a := range n.Attr {
			if a.Key == "style" {
				n.Attr[i].Val = OfCss(a.Val, e.cssTransformer)
			}
		}
	}

	// done converting links from the current html node,
	//now, recursively, extract links from all child nodes to the current node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		e.extract(c)
	}

	return
}

// textContentNode returns the text content of the target html node, and the
// target html.Node from where the text content was extracted
func textContentNode(n *html.Node) (string, *html.Node) {
	var content strings.Builder
	var tcNode *html.Node
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			content.WriteString(child.Data)
			tcNode = child
		}
	}

	return content.String(), tcNode
}
