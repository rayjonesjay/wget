package links

import (
	"golang.org/x/net/html"
	"wget/convertlinks"
)

// FromHtml extracts and returns a list of all URLs linked in the target [html.Node].
//
// It retrieves linked URLs from the following HTML elements:
// - <a>: Hyperlinks
// - <img>: Images
// - <link>: Stylesheets
// - <style>: Style sheets, including @import statements
// - <object>: Embedded objects
//
// If no linked URLs are found, it returns a nil slice. If the input document is malformed,
// it may result in an empty slice or unexpected behavior.
func FromHtml(doc *html.Node) (urls []string) {
	urlsExtractor := func(url string, isA bool) string {
		urls = append(urls, url)
		return url
	}
	convertlinks.OfHtml(doc, urlsExtractor)
	return
}
