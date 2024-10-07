package mirror

import (
	"net/http"
	"wget/httpx"
	"wget/mirror/xurl"
)

// GetFile returns the file path where the contents of the provided URL to be mirrored will be written into,
// honoring the given parent folder. The downloaded file will be stored in a directory structure that reflects the
// mirrored URL, with the given parent directory.
//
// The function assumes that the provided downloadUrl is a valid URL and that
// the parent folder (if specified) is a writable directory.
func GetFile(downloadUrl string, header http.Header, parentFolder string) string {
	// Download the base html file, say `index.html`
	loc := xurl.DownloadFolder(downloadUrl)
	path := xurl.TrimSlash(loc.FolderName)
	if parentFolder != "" {
		path = xurl.TrimSlash(parentFolder) + "/" + path
	}

	if loc.FileName == "." {
		loc.FileName = "index.html"
		contentType := httpx.ExtractMimeType(header)
		if ext, ok := contentTypeExtensions[contentType]; ok {
			loc.FileName += "." + ext
		}
	}

	path = xurl.TrimSlash(path) + "/" + loc.FileName
	return path
}
