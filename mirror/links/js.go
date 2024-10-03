package links

import (
	"fmt"
	"regexp"
)

// FromJs given a string of Javascript code,
// extracts links to Javascript modules imported by the given Javascript code; returns a list of the extracted links.
// It supports module imports and commonJs imports (via `require`).
// This uses a poor-man's Javascript parser, therefore,
// it does not filter out `import` statements and `require` import calls that are within comments or embedded within
// strings; neither does it check to affirm that the Javascript code is syntactically correct. As a consequence,
// it may return more linked modules than is the case
func FromJs(js string) (moduleLinks []string) {
	moduleLinks = append(moduleLinks, FromJsImport(js)...)
	moduleLinks = append(moduleLinks, FromJsRequire(js)...)
	return
}

// FromJsImport takes a string containing JavaScript import statements and returns a slice of the linked modules.
func FromJsImport(js string) []string {
	// Regular expression to match import statements
	re := regexp.MustCompile("import\\s+[\\w\\s*{},]*from\\s+[\"'`]([^\"'`]+)[\"`'];?")

	// Find all matches
	matches := re.FindAllStringSubmatch(js, -1)
	re.ReplaceAllString(js, "")

	re = regexp.MustCompile("import\\s+[\"'`]([^\"'`]+)[\"`']")
	matches = append(matches, re.FindAllStringSubmatch(js, -1)...)

	// Extract the module paths
	var modules []string
	for _, match := range matches {
		if len(match) > 1 {
			modules = append(modules, match[1])
		}
	}

	return modules
}

// FromJsRequire takes a string containing JavaScript require statements and returns a slice of the linked modules.
func FromJsRequire(js string) []string {
	// Regular expression to match require statements
	quotes := "\"'`"
	// `require\([`"']([^`"']+)[`"']\)`
	re := regexp.MustCompile(fmt.Sprintf(`require\([%s]([^%s]+)[%s]\)`, quotes, quotes, quotes))

	// Find all matches
	matches := re.FindAllStringSubmatch(js, -1)

	// Extract the module paths
	var modules []string
	for _, match := range matches {
		if len(match) > 1 {
			modules = append(modules, match[1])
		}
	}

	return modules
}
