package mirror

import (
	"path"
	"regexp"
)

func (a *arg) ShouldReject(mirrorPath string) bool {
	mirrorPath = path.Base(mirrorPath)
	for _, pattern := range a.rejectPatterns {
		if pattern.MatchString(mirrorPath) {
			return true
		}
	}

	return false
}

func (a *arg) ShouldExclude(mirrorPath string) bool {
	mirrorPath = path.Dir(mirrorPath)
	for _, pattern := range a.excludePatterns {
		if pattern.MatchString(mirrorPath) {
			return true
		}
	}

	return false
}

func (a *arg) initReject() {
	for _, reject := range a.Rejects {
		regex, err := buildRegex(reject, true)
		if err != nil {
			continue
		}
		a.rejectPatterns = append(a.rejectPatterns, regex)
	}
}

func (a *arg) initExclude() {
	for _, exclude := range a.Exclude {
		regex, err := buildRegex(exclude, false)
		if err != nil {
			continue
		}
		a.excludePatterns = append(a.excludePatterns, regex)
	}
}

func buildRegex(pattern string, isReject bool) (*regexp.Regexp, error) {
	re := regexp.MustCompile(`(\*)|(\?)|(\[.*])`)
	//pattern := `hello-[1-3]-?.*.png`
	//pattern := `hello-[1-3][1-3]-?.*.png`
	//pattern := `hello-[1-3][1-3]-??.*.*.png`
	//pattern := `hello`

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

	if out == pattern {
		if isReject {
			return regexp.Compile(`.*` + pattern + `$`)
		} else {
			return regexp.Compile(`^/?` + pattern + `.*`)
		}
	}

	return regexp.Compile(out)
}
