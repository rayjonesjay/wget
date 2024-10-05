// Package mirror contains functionalities that allow one to download an entire website
package mirror

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
	"sync"
	"wget/convertlinks"
	"wget/ctx"
	"wget/fetch"
	"wget/fileio"
	"wget/httpx"
)

// map of content types to preferred file extensions
var contentTypeExtensions = map[string]string{
	"text/html":                "html",
	"text/css":                 "css",
	"text/javascript":          "js",
	"application/javascript":   "js",
	"application/json":         "json",
	"application/xml":          "xml",
	"application/pdf":          "pdf",
	"image/jpeg":               "jpg",
	"image/png":                "png",
	"image/gif":                "gif",
	"image/svg+xml":            "svg",
	"audio/mpeg":               "mp3",
	"audio/wav":                "wav",
	"video/mp4":                "mp4",
	"video/webm":               "webm",
	"application/zip":          "zip",
	"application/x-gzip":       "gz",
	"application/octet-stream": "bin",
}

// arg embeds the download context to add custom receiver functions
type arg struct {
	*ctx.Context
	// downloaded keeps a map based list of urls that have been downloaded,
	//or scheduled for download by this instance. A map is preferred for O(1), constant time, existential checks
	//This helps avoid re-downloading a file from the same URL more than once
	downloaded map[string]bool
	// mutex locks and unlocks this instance when accessed by multiple goroutines
	mutex *sync.Mutex
}

var ErrFileAlreadyDownloaded = errors.New("file already downloaded")

// Site downloads the entire website being possible to use "part" of the website offline.
// If no scheme is detected in the mirror URL, then, the HTTP scheme is assumed
func Site(cxt ctx.Context, mirrorUrl string) error {
	m := &arg{
		Context:    &cxt,
		downloaded: make(map[string]bool),
		mutex:      &sync.Mutex{},
	}
	parse, err := url.Parse(mirrorUrl)
	if err != nil {
		return err
	}

	if parse.Scheme == "" {
		// to be compatible with other clients,
		//assume the scheme is http
		parse.Scheme = "http"
	}

	_, err = m.Site(parse.String())
	return err
}

// GetFile returns a writable file, where the downloaded file will be written into,
// or an error if it fails. GetFile honours the current download context as specified by this instance
func (a *arg) GetFile(downloadUrl string, header http.Header) (*os.File, error) {
	// Download the base html file, say `index.html`
	loc := DownloadFolder(downloadUrl)
	path := TrimSlash(loc.FolderName)
	if a.SavePath != "" {
		path = TrimSlash(a.SavePath) + "/" + path
	}

	err := os.MkdirAll(path, 0775)
	if err != nil {
		return nil, err
	}

	if loc.FileName == "." {
		loc.FileName = "index"
		contentType := httpx.ExtractMimeType(header)
		if ext, ok := contentTypeExtensions[contentType]; ok {
			loc.FileName += "." + ext
		}
	}

	path = TrimSlash(path) + "/" + loc.FileName
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
}

// Site downloads the entire website being possible to use "part" of the website offline,
// respecting the download context defined by the given instance.
// If no scheme is detected in the mirror URL, then, the HTTP scheme is assumed
func (a *arg) Site(mirrorUrl string) (info fetch.FileInfo, err error) {
	{ // check if the given URL has already been downloaded by this instance
		a.mutex.Lock()
		if _, ok := a.downloaded[mirrorUrl]; ok {
			// skip, already downloaded or in the download queue
			return info, ErrFileAlreadyDownloaded
		}
		a.downloaded[mirrorUrl] = true
		a.mutex.Unlock()
	}

	// TODO add progress and rate listeners
	info, err = fetch.URL(
		mirrorUrl,
		fetch.Config{
			GetFile:          a.GetFile,
			Limit:            int32(a.RateLimitValue),
			ProgressListener: nil,
			RateListener:     nil,
			Body:             nil,
			Method:           "GET",
		},
	)
	if err != nil {
		return
	}

	contentType := httpx.ExtractMimeType(info.Headers)
	// TODO: parse other web formats as CSS
	if contentType != "text/html" {
		// Not a html file, done downloading
		return
	}

	// the downloaded file is HTML; attempt to extract linked contents and pages
	htmlFile, err := os.Open(info.Name)
	if err != nil {
		return
	}
	defer fileio.Close(htmlFile)

	doc, err := html.Parse(htmlFile)
	if err != nil {
		return
	}

	links := Extract(doc)
	fmt.Printf("Found links: %v\n", links)
	if len(links) == 0 {
		// No more links to download
		return
	}

	// Add a link to the site icon, /favicon.ico (most clients love it),
	//and to the /robots.txt file (wget downloads this file in mirror mode)
	for _, relUrl := range []string{"/robots.txt", "/favicon.ico"} {
		faviconUrl, err := AbsoluteUrl(mirrorUrl, relUrl)
		if err != nil {
			return info, err
		}
		links = append(links, UrlExtract{Url: faviconUrl})
	}

	convertUrls := make(map[string]string)
	// Download each link, synchronously, continuing to the next regardless of errors
	for _, link := range links {
		mUrl, err := AbsoluteUrl(mirrorUrl, link.Url)
		if err != nil {
			fmt.Println(err)
			continue
		}

		linkInfo, err := a.Site(mUrl)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("saved to -> %v\n", linkInfo)
			convertUrls[link.Url] = RelativeFolder(info.Name, linkInfo.Name)
		}
	}

	// TODO: remove tester assignment
	a.ConvertLinks = true
	for a.ConvertLinks {
		linkConverter := func(url string, isA bool) string {
			if toUrl, ok := convertUrls[url]; ok {
				return toUrl
			}
			return url
		}
		convertlinks.OfHtml(doc, linkConverter)

		// Write the new doc `html.Node` to a new temporary file
		convertHtmlFile, err := createTempFile()
		if err != nil {
			break
		}
		defer fileio.Close(convertHtmlFile)

		err = html.Render(convertHtmlFile, doc)
		if err != nil {
			fmt.Println(err)
			break
		}

		// successfully converted the links, and wrote the HTML node to the new temporary file,
		// move the temporary file to the actual downloaded HTML file
		fileio.Close(htmlFile)
		err = os.Rename(convertHtmlFile.Name(), htmlFile.Name())
		if err != nil {
			break
		}

		break
	}

	return
}

func createTempFile() (*os.File, error) {
	dir := os.TempDir() + "com.zone01.wget"
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return nil, err
	}

	// Create a temporary file inside the directory
	tempFile, err := os.CreateTemp(dir, "*.tmp")
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}
