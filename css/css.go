// Package css contains utilities to parse css and optionally transform the stricture of the parsed css
package css

import (
	"bytes"
	"github.com/tdewolff/parse/css"
	"regexp"
	"strings"
)

// TransformCssUrl parses the given CSS, checking for calls to the CSS url() functions,
// then modifying the url link defined in the url() function according to the defined transformer.
// returns the transformed CSS string. Does nothing if the transformer function is nil
func TransformCssUrl(cssStr string, transformer func(url string) string) string {
	if transformer == nil {
		return cssStr
	}

	l := css.NewLexer(strings.NewReader(cssStr))
	b := strings.Builder{}
	for {
		tt, data := l.Next()
		//fmt.Printf("%s -> %s\n", tt, data)
		if tt == css.ErrorToken {
			break
		} else if tt == css.URLToken {
			//fmt.Printf("URL: `%s`\n", data)
			// url("path/to/image.jpg")
			re := regexp.MustCompile("(url)(\\s*)(\\()(\\s*)(['\"`]?)([^'\"`]*)(['\"`]?)(\\s*)(\\))(\\s*)")
			matches := re.FindSubmatch(data)
			if matches != nil {
				urlGroup := 6
				url := string(matches[urlGroup])
				transformedUrl := transformer(url)
				matches[urlGroup] = []byte(transformedUrl)
				data = bytes.Join(matches[1:], []byte(""))
			}
		}

		b.Write(data)
	}
	return b.String()
}
