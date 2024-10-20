package mirror

import (
	"path"
	"regexp"
)

// ShouldReject returns true if the given mirror URL path, refers to a file that
// should not be downloaded; false otherwise
func (a *arg) ShouldReject(mirrorPath string) bool {
	mirrorPath = path.Base(mirrorPath)
	for _, pattern := range a.rejectPatterns {
		if pattern.MatchString(mirrorPath) {
			return true
		}
	}

	return false
}

// ShouldExclude returns true if the given mirror URL path, refers to a directory that
// should not be downloaded; false otherwise
func (a *arg) ShouldExclude(mirrorPath string) bool {
	mirrorPath = path.Dir(mirrorPath)
	for _, pattern := range a.excludePatterns {
		if pattern.MatchString(mirrorPath) {
			return true
		}
	}

	return false
}

// initReject creates a list of compiled regular expressions, that would be used
// to match whether a given file should not be downloaded
func (a *arg) initReject() {
	if len(a.rejectPatterns) != 0 {
		panic("rejectPatterns must be empty; init should ideally be called once after struct creation")
	}

	for _, reject := range a.Rejects {
		regex, err := buildRegex(reject, true)
		if err != nil {
			continue
		}
		a.rejectPatterns = append(a.rejectPatterns, regex)
	}
}

// initExclude creates a list of compiled regular expressions, that would be used
// to match whether a given directory should not be downloaded
func (a *arg) initExclude() {
	if len(a.excludePatterns) != 0 {
		panic("excludePatterns must be empty; init should ideally be called once after struct creation")
	}
	for _, exclude := range a.Exclude {
		regex, err := buildRegex(exclude, false)
		if err != nil {
			continue
		}
		a.excludePatterns = append(a.excludePatterns, regex)
	}
}

// buildRegex attempts to compile a regex, from the given pattern as passed to
// the directory-based-limits command-line arguments as --reject and --exclude,
// (it takes a boolean to differentiate which of the two the pattern was
// extracted)
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
