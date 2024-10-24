package xurl

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// AbsoluteUrl returns an absolute URL from the given URL (relative or absolute),
// with a basis on the given parent URL.
func AbsoluteUrl(parentURL, relativeURL string) (string, error) {
	// parse cleans the url before parsing it
	parse := func(_url string) (*url.URL, error) {
		u, err := CleanUrl(_url)
		if err != nil {
			return nil, err
		}

		u2, _ := url.Parse(u)
		return u2, nil
	}

	relative, err := parse(relativeURL)
	if err != nil {
		return "", err
	}

	if relative.IsAbs() {
		return relative.String(), nil
	}

	parent, err := parse(parentURL)
	if err != nil {
		return "", err
	}

	// Resolve the relative URL against the parent URL.
	// Note that the `url.ResolveReference` function assumes the target url is a resource and not a
	// resource folder when it's not suffixed by a `/` character, we work around this
	// by appending `/` character
	if !strings.HasSuffix(parent.Path, "/") {
		parent.Path += "/"
	}

	absolute := parent.ResolveReference(relative)
	return absolute.String(), nil
}

// TrimSlash returns a new url, with trailing slash character in the given url removed
func TrimSlash(targetUrl string) string {
	return strings.TrimSuffix(targetUrl, "/")
}

// CleanUrl returns a new string with all extra `/` characters, that may have
// been included in the URL, removed
func CleanUrl(targetUrl string) (string, error) {
	u, err := url.Parse(targetUrl)
	if err != nil {
		return "", err
	} else if u.Opaque != "" {
		return "", fmt.Errorf("we don't yet handle opaque URLs: %s", u.Opaque)
	}

	u.Path = CleanSlash(u.Path)

	if u.Host == "" && u.Path != "" && u.Scheme != "" && u.Scheme != "file" {
		// If the URI defined a scheme, ensure only file URI supports leading / in the URI path
		u.Path = strings.TrimLeft(u.Path, "/")
	}

	return u.String(), nil
}

// CleanSlash returns a new path string, with all series of more than one `/`
// character, replaced with a single `/` character. Note that the path doesn't
// have to be a valid URL, just any path, potentially with `/` characters to be
// checked
func CleanSlash(path string) string {
	re := regexp.MustCompile(`/{2,}`)
	return re.ReplaceAllString(path, "/")
}

// SameHost checks if two URLs have the same host.
//
// It parses both input URLs and compares their host components.
// If either URL is malformed, as defined by [url.Parse], an error is returned.
//
// Returns:
// - true if the hosts are the same, false otherwise.
// - An error if either URL is malformed.
func SameHost(url1, url2 string) bool {
	parsedURL1, err1 := url.Parse(url1)
	parsedURL2, err2 := url.Parse(url2)
	if err1 == nil && err2 == nil {
		return parsedURL1.Host == parsedURL2.Host
	}
	return false
}
