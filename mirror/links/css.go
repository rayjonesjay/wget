package links

import (
	"strings"
	"wget/css"
)

// FromCssUrl takes a CSS string and returns a slice of URLs extracted from the CSS,
// including those from url() function calls
func FromCssUrl(cssStr string) (urls []string) {
	// the function `css.TransformCssUrl` may match urls with quotes, such as, "img/avatar.png", 'img/avatar.png',
	//or `img/avatar.png`; we use this utility function to remove such quotes, then append the url to the list of urls
	trimAppend := func(url string) string {
		// avoid including malformed urls as:
		// - `"img/avatar.png` -> no closing `"`
		// - `img/avatar.png"` -> no opening `"`
		// - `""img/avatar.png"` -> more than one opening `"`
		// - `"img/avatar.png""` -> more than one closing `"`
		trimmedUrl := strings.TrimFunc(
			url, func(r rune) bool {
				return strings.ContainsRune("\"`'", r)
			},
		)
		urls = append(urls, trimmedUrl)
		return url
	}

	css.TransformCssUrl(cssStr, trimAppend)
	return
}
