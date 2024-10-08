package mirror

import (
	"net/http"
	"os"
	"path/filepath"
	fpath "path/filepath"
	"slices"
	"strings"
	"wget/mirror/xurl"
)

// GetFile returns the file path where the contents of the provided URL to be mirrored will be written into,
// honoring the given parent folder. The downloaded file will be stored in a directory structure that reflects the
// mirrored URL, with the given parent directory.
//
// The function assumes that the provided downloadUrl is a valid URL and that
// the parent folder (if specified) is a writable directory.
//func GetFile(downloadUrl string, header http.Header, parentFolder string) string {
//	// Download the base html file, say `index.html`
//	loc := xurl.DownloadFolder(downloadUrl)
//	path := xurl.TrimSlash(loc.FolderName)
//	if parentFolder != "" {
//		path = xurl.TrimSlash(parentFolder) + "/" + path
//	}
//
//	if loc.FileName == "." {
//		contentType := httpx.ExtractMimeType(header)
//		if ext, ok := contentTypeExtensions[contentType]; ok {
//			loc.FileName = fmt.Sprintf("index.%s", ext)
//		} else {
//			loc.FileName = "index.html"
//		}
//	}
//
//	path = xurl.TrimSlash(path) + "/" + loc.FileName
//	return path
//}

// FolderStructure returns all parent folders necessary for the given filepath to exist.
// The returned slice contains the directories from the provided filepath up to the root.
// For example, given "/a/b/c/file.txt", it will return ["/a/b/c", "/a/b", "/a"].
func FolderStructure(filePath string) (structure []string) {
	specialCases := []string{"", ".", "/"}
	isRootFolder := func(path string) bool {
		return slices.Contains(specialCases, path)
	}

	for {
		filePath = xurl.TrimSlash(filePath)
		if isRootFolder(filePath) {
			break
		}
		dir := filepath.Dir(filePath)
		if isRootFolder(dir) {
			break
		}
		structure = append(structure, dir)
		filePath = dir
	}
	return
}

// ForceMkdirAll is a destructive variant of the MkdirAll function, that instead
// ensures that all parent folders anticipated by the given filepath will be
// created even if there exists a file with the same name as any of its parent
// folders. Any file with the same name as any parent directory of filepath will
// be deleted
func ForceMkdirAll(filePath string) (err error) {
	structure := FolderStructure(filePath)
	var info os.FileInfo
	for i := len(structure) - 1; i >= 0; i-- {
		file := structure[i]
		err = os.MkdirAll(file, 0775)
		if err == nil {
			continue
		}
		info, err = os.Stat(file)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			err = os.Remove(file)
			if err != nil {
				return err
			}
			err = os.MkdirAll(file, 0775)
			if err != nil {
				return err
			}
		}

	}

	return
}

// GetFile determines the destination file path based on the URL and headers
func GetFile(link string, header http.Header, parentFolder string) (destination string) {
	indexFile := "index.html"

	// Remove the scheme (http:// or https://) from the URL
	link = removeScheme(link)

	// Extract Content-Type from the header
	contentType := header.Get("Content-Type")

	// Get the subdirectory and file name from the URL
	subDir, file := base(link, '/')

	// Determine the appropriate file name based on content type and URL structure
	switch {
	case strings.Contains(contentType, "text/html") || strings.Contains(contentType, "text/css"):
		// HTML or CSS files
		if len(file) == 0 {
			return fpath.Join(parentFolder, subDir, indexFile)
		}
		return fpath.Join(parentFolder, subDir, file)
	case contentType == "":
		// When content type is empty, assume HTML and default to index.html
		if len(file) == 0 {
			return fpath.Join(parentFolder, subDir, indexFile)
		}
		return fpath.Join(parentFolder, subDir, file)
	default:
		// For other content types (images, binaries, etc.)
		if len(file) == 0 {
			file = "unknownfile"
		}
		return fpath.Join(parentFolder, subDir, file)
	}
}

// removeScheme removes the "http://" or "https://" scheme from a URL
func removeScheme(link string) string {
	if strings.HasPrefix(link, "https://") {
		return strings.TrimPrefix(link, "https://")
	}
	return strings.TrimPrefix(link, "http://")
}

// base extracts the subdirectory and file name from the URL
func base(link string, delimiter byte) (string, string) {
	// Strip query parameters if any (everything after '?')
	queryIndex := strings.Index(link, "?")
	if queryIndex != -1 {
		link = link[:queryIndex]
	}

	// If the link ends with a slash, it's a directory, not a file
	if strings.HasSuffix(link, "/") {
		return link, ""
	}

	// Find the last occurrence of the delimiter ('/')
	index := strings.LastIndexByte(link, delimiter)

	if index == -1 {
		// No slash found, so it's just the file name without a directory
		return "", link
	}

	// Extract the directory and file name
	directory, file := link[:index+1], link[index+1:]
	return directory, file
}
