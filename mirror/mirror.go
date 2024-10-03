// Package mirror contains functionalities that allow one to download an entire website
package mirror

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
	"sync"
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

	return m.Site(parse.String())
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
func (a *arg) Site(mirrorUrl string) error {
	{ // check if the given URL has already been downloaded by this instance
		a.mutex.Lock()
		if _, ok := a.downloaded[mirrorUrl]; ok {
			// skip, already downloaded or in the download queue
			return nil
		}
		a.downloaded[mirrorUrl] = true
		a.mutex.Unlock()
	}

	// TODO add progress and rate listeners
	info, err := fetch.URL(
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
		return err
	}

	contentType := httpx.ExtractMimeType(info.Headers)
	// TODO: parse other web formats as CSS
	if contentType != "text/html" {
		// Not a html file, done downloading
		return nil
	}

	// the downloaded file is HTML; attempt to extract linked contents and pages
	htmlFile, err := os.Open(info.Name)
	if err != nil {
		return err
	}
	defer fileio.Close(htmlFile)

	doc, err := html.Parse(htmlFile)
	if err != nil {
		return err
	}

	links := Extract(doc)
	fmt.Printf("Found links: %v\n", links)
	if len(links) == 0 {
		// No more links to download
		return nil
	}

	// Add a link to the site icon, /favicon.ico, most clients love it
	{
		faviconUrl, err := AbsoluteUrl(mirrorUrl, "/favicon.ico")
		if err != nil {
			return err
		}
		links = append(links, UrlExtract{Url: faviconUrl})
	}

	// Download each link
	wg := sync.WaitGroup{}
	for _, link := range links {
		mUrl, err := AbsoluteUrl(mirrorUrl, link.Url)
		if err != nil {
			fmt.Println(err)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := a.Site(mUrl)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	wg.Wait()
	if a.ConvertLinks {
		err := html.Render(htmlFile, doc)
		if err != nil {
			return err
		}
	}

	return nil
}
