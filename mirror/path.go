package mirror

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"wget/fileio"
	"wget/httpx"
	"wget/mirror/xurl"
	"wget/temp"
)

// GetFile returns the file path where the contents of the provided URL to be mirrored will be written into,
// honoring the given parent folder. The downloaded file will be stored in a directory structure that reflects the
// mirrored URL, with the given parent directory.
//
// The function assumes that the provided downloadUrl is a valid URL and that
// the parent folder (if specified) is a writable directory.
func GetFile(downloadUrl string, header http.Header, parentFolder string) string {
	u, err := url.Parse(downloadUrl)
	if err != nil {
		return ""
	}

	if u.Hostname() == "" && u.Path == "" {
		return ""
	}

	folder := filepath.Dir(u.Path)
	filename := strings.TrimPrefix(u.Path, folder)

	if filename == "" || filename == "/" {
		contentType := httpx.ExtractMimeType(header)
		if ext, ok := contentTypeExtensions[contentType]; ok {
			filename = fmt.Sprintf("index.%s", ext)
		} else {
			filename = "index.html"
		}
	}

	if disposition, err := httpx.FilenameFromContentDisposition(header); err == nil {
		filename = disposition
	}

	host := u.Hostname()
	if host != "" {
		host += "/"
	}
	// host => google.com/
	// include the hostname to the folder name
	folder = host + folder
	return path.Join(parentFolder, folder, filename)
}

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
// be moved to the created parent directory and renamed "index.html"
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
			// create a temporary file to hold `file` as we await to move it
			tempFile, err := temp.File()
			if err != nil {
				return err
			}
			fileio.Close(tempFile)

			// rename the file to the temporary file
			err = os.Rename(file, tempFile.Name())
			if err != nil {
				return err
			}

			// now create a directory instead with the filename
			err = os.MkdirAll(file, 0775)
			if err != nil {
				return err
			}

			// move the file into the newly created directory
			err = os.Rename(tempFile.Name(), path.Join(file, "index.html"))
			if err != nil {
				return err
			}
		}
	}

	return
}
