package reject

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"wget/ctx"
)

//
//type ABString struct {
//	isAB bool
//	string
//}
//
//const (
//	WildcardStar = iota
//	WildcardRepetitionZeroOrOne
//	WildcardCharacterClass
//)
//
//var regexCache = make(map[string]*regexp.Regexp)
//
//type Wildcard struct {
//	Type    int
//	Pattern string
//}

func Reject(ctx ctx.Context, mirrorPath string) bool {
	wildcards := []string{"*", "?", "[", "]"}
	if !slices.Contains(wildcards, mirrorPath) {
		// suffix
		for _, suffix := range ctx.Rejects {
			if strings.HasSuffix(mirrorPath, suffix) {
				return true
			}
		}
		return false
	}

	// pattern such as -> *.png
	for _, reject := range ctx.Rejects {
		re, err := buildRegex(reject)
		if err != nil {
			continue
		}

		if re.MatchString(mirrorPath) {
			return true
		}
	}

	return false
}

func buildRegex(pattern string) (*regexp.Regexp, error) {
	re := regexp.MustCompile(`(\*)|(\?)|(\[.*])`)
	//pattern := `hello-[1-3]-?.*.png`
	//pattern := `hello-[1-3][1-3]-?.*.png`
	//pattern := `hello-[1-3][1-3]-??.*.*.png`
	//pattern := `opo`
	fmt.Printf("%q\n", re.FindAllStringSubmatch(pattern, -1))
	fmt.Printf("%q\n", re.Split(pattern, -1))

	matches := re.FindAllStringSubmatch(pattern, -1)
	splits := re.Split(pattern, -1)

	out := ""
	for i, sp := range splits {
		out += regexp.QuoteMeta(sp)
		if i != len(splits)-1 && i < len(matches) {
			m := matches[i][0]
			switch m {
			case "*":
				out += ".*"
			case "?":
				out += ".?"
			default:
				// []
				out += m
			}
		}
	}

	return regexp.Compile(out)
}
