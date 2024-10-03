package mirror

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

// DownloadLocation records the names of the folder where a given resource should
// be downloaded to, and the name of the file the resource should be written into
type DownloadLocation struct {
	FolderName, FileName string
}

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

// DownloadFolder returns the absolute filepath, from the current working
// directory, where the resource retrieved from the given url should be
// downloaded to.
// The folder name and/or file name may be the `.` character, in
// which case, the file should be downloaded to the current directory, and most
// probably named the same as the current folder
// This function assumes that the given url is a valid
// url; if not, then the returned folder and file names will be empty strings
func DownloadFolder(targetUrl string) (loc DownloadLocation) {
	u, err := url.Parse(targetUrl)
	if err != nil {
		return
	}

	host := ""
	if u.Host != "" {
		host = u.Host + "/"
	}

	loc.FolderName = filepath.Dir(host + u.Path)
	loc.FileName = filepath.Base(u.Path)

	// For consistency, we translate any `/` filename as may be returned by filepath.Base to `.`
	if loc.FileName == "/" {
		loc.FileName = "."
	}

	return
}

// RelativeFolder calculates the relative path of the target folder from the parent folder
func RelativeFolder(parent, target string) string {
	if parent == "" {
		return target
	} else if target == "" {
		return "."
	} else if target == "/" {
		return parent
	}

	// Clean the paths to handle different separators and potential issues.
	parent = filepath.Clean(parent)
	target = filepath.Clean(target)

	// Calculate the relative target
	rel, err := filepath.Rel(parent, target)

	// Handle potential errors, particularly when paths are on different drives in Windows.
	if err != nil {
		return target
	}

	return rel
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
		return "", fmt.Errorf("we don't yet handle opaque URLS: %s", u.Opaque)
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
