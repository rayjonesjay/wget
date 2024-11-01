// Package mirror contains functionalities that allow one to download an entire website
package mirror

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
	"wget/convertlinks"
	"wget/ctx"
	"wget/fetch"
	"wget/fileio"
	"wget/globals"
	"wget/httpx"
	"wget/mirror/links"
	"wget/mirror/xurl"
	"wget/syscheck"
	"wget/temp"
)

// map of content types to preferred file extensions
var contentTypeExtensions = map[string]string{
	"text/html":              "html",
	"text/plain":             "txt",
	"text/css":               "css",
	"text/javascript":        "js",
	"application/javascript": "js",
	"application/json":       "json",
	"application/xml":        "xml",
	"application/pdf":        "pdf",
	"image/jpeg":             "jpg",
	"image/png":              "png",
	"image/gif":              "gif",
	"image/svg+xml":          "svg",
	"audio/mpeg":             "mp3",
	"audio/wav":              "wav",
	"video/mp4":              "mp4",
	"video/webm":             "webm",
	"application/zip":        "zip",
	"application/x-gzip":     "gz",
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
	// urlDownloadInfo maps the results of downloading a given URL,
	//such as the filename the contents of the URL was written to,
	//and whether there was any error when downloading the file
	urlDownloadInfo map[string]UrlDownloadInfo
	// rejectPatterns keeps a list of regex patterns to match files to be rejected for download
	rejectPatterns []*regexp.Regexp
	// excludePatterns keeps a list of regex patterns to match directories to be rejected for download
	excludePatterns []*regexp.Regexp
	// d keeps the download id of the current mirror URL. Since, we are downloading
	// recursively, this integer keeps the number of mirrored URLS, thus, can be used
	// to assign unique integer IDs to each downloads
	d int
	// df records the total number of bytes downloaded in the process of this mirror
	df int64
}

// UrlDownloadInfo keeps the results of downloading a given URL,
// such as the filename the contents of the URL was written to,
// and whether there was any error when downloading the file
type UrlDownloadInfo struct {
	url string
	fetch.FileInfo
	error
}

var ErrFileAlreadyDownloaded = errors.New("file already downloaded")

// Site downloads the entire website being possible to use "part" of the website offline.
// If no scheme is detected in the mirror URL, then, the HTTP scheme is assumed
func Site(cxt ctx.Context, mirrorUrl string) error {
	m := &arg{
		Context:         &cxt,
		downloaded:      make(map[string]bool),
		mutex:           &sync.Mutex{},
		urlDownloadInfo: make(map[string]UrlDownloadInfo),
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

	syscheck.MoveCursor(1)
	syscheck.ClearScreen()
	syscheck.HideCursor()
	defer syscheck.ShowCursor() // Ensure cursor is shown again when done

	m.init()
	startTime := time.Now()
	_, err = m.Site(parse.String())
	endTime := time.Now()
	duration := endTime.Sub(startTime).Truncate(time.Second)
	fmt.Printf("\n\nFINISHED --%s--\n"+
		"Total wall clock time: %s\n"+
		"Downloaded: %d files, %s in %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		duration.String(),
		m.d,
		globals.FormatSize(m.df),
		duration.String(),
	)
	return err
}

// init is called once, when the struct instance is created, to initialize various fields to their usable values
func (a *arg) init() {
	a.initReject()
	a.initExclude()
}

// GetFile returns a writable file, where the downloaded file will be written into,
// or an error if it fails. GetFile honours the current download context as specified by this instance
func (a *arg) GetFile(downloadUrl string, header http.Header) (*os.File, error) {
	downloadPath := GetFile(downloadUrl, header, a.SavePath)
	log.Printf("downloading url %q -> %q\n", downloadUrl, downloadPath)
	err := ForceMkdirAll(downloadPath)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(downloadPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
}

// ShouldDownload will be called to validate whether the file from the given url
// should be downloaded, based on the given headers as retrieved from the server.
// This will always download HTML files, so that we can extract linked URLs from
// them, and later delete them if the directory-based-limits infer
func (a *arg) ShouldDownload(mirrorUrl string, header http.Header) bool {
	if httpx.ExtractMimeType(header) == "text/html" {
		return true
	}

	parse, err := url.Parse(mirrorUrl)
	if err != nil {
		return false
	}

	if a.ShouldReject(parse.Path) || a.ShouldExclude(parse.Path) {
		return false
	}

	return true
}

// Site downloads the entire website being possible to use "part" of the website offline,
// respecting the download context defined by the given instance.
// If no scheme is detected in the mirror URL, then, the HTTP scheme is assumed
func (a *arg) Site(mirrorUrl string) (info fetch.FileInfo, err error) {
	log.Printf("[1] Fetching >> %q\n", mirrorUrl)
	defer log.Printf("[1] Done\n")
	// check if the given URL has already been downloaded by this instance
	err = func() error {
		a.mutex.Lock()
		defer a.mutex.Unlock()
		if _, ok := a.downloaded[mirrorUrl]; ok {
			// skip, already downloaded or in the download queue
			return ErrFileAlreadyDownloaded
		}
		a.downloaded[mirrorUrl] = true
		return nil
	}()

	if err != nil {
		return
	}

	// define a download status listener for the current mirror URL
	status := fetch.DownloadStatus{}
	status.OnUpdate = func(status *fetch.DownloadStatus, hint int) {
		// whenever the status of this download has changed, we print the progress to its
		// respective position in the terminal
		globals.PrintLines(
			(a.d-1)*8, []string{
				status.Start,
				status.Status,
				status.ContentLength,
				status.SavePath,
				status.Progress,
				status.Finished,
			},
		)
	}

	// configure an Advanced Progress Listener for the GET request
	advancedProgressListener := *status.ProgressListener()
	{
		// need to increment the ID when the download actually starts, inject an onstart listener
		originalOnStart := advancedProgressListener.OnStart
		advancedProgressListener.OnStart = func(time time.Time) {
			a.d++
			originalOnStart(time)
		}
	}
	{
		// need to keep track of the total bytes downloaded
		originalOnDownloadFinished := advancedProgressListener.OnDownloadFinished
		advancedProgressListener.OnDownloadFinished = func(url string, time time.Time) {
			originalOnDownloadFinished(url, time)
			a.mutex.Lock()
			a.df += status.Downloaded
			a.mutex.Unlock()
		}
	}

	info, err = fetch.URL(
		mirrorUrl,
		fetch.Config{
			GetFile:                  a.GetFile,
			ShouldDownload:           a.ShouldDownload,
			Limit:                    int32(a.RateLimitValue),
			ProgressListener:         nil,
			RateListener:             nil,
			Body:                     nil,
			Method:                   "GET",
			AllowedStatusCodes:       []int{http.StatusOK},
			AdvancedProgressListener: advancedProgressListener,
		},
	)
	if err != nil {
		return
	}

	// Save the results of the downloaded resource
	a.urlDownloadInfo[mirrorUrl] = UrlDownloadInfo{
		url:      mirrorUrl,
		FileInfo: info,
		error:    err,
	}

	if !a.ShouldDownload(mirrorUrl, http.Header{}) {
		defer func(name string) {
			_ = os.Remove(name)
		}(info.Name)
	}

	contentType := httpx.ExtractMimeType(info.Headers)
	if contentType == "text/css" {
		a.FetchCss(mirrorUrl, info.Name)
		return
	} else if contentType != "text/html" {
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

	linkedUrls := links.FromHtml(doc)
	log.Printf("Found linkedUrls: %v\n", linkedUrls)
	if len(linkedUrls) == 0 {
		// No more links to download
		return
	}

	convertUrls := make(map[string]string)
	// Download each link, synchronously, continuing to the next regardless of errors
	for _, link := range linkedUrls {
		linkUrl, err := xurl.AbsoluteUrl(mirrorUrl, link)
		if err != nil {
			log.Println(err)
			continue
		} else if !xurl.SameHost(mirrorUrl, linkUrl) {
			continue
		}

		linkInfo, err := a.Site(linkUrl)
		if _, ok := a.urlDownloadInfo[link]; errors.Is(err, ErrFileAlreadyDownloaded) && ok {
			linkInfo = a.urlDownloadInfo[link].FileInfo
			err = nil
		}

		if err != nil {
			log.Println(err)
		} else {
			log.Printf("saved to -> %s\n", linkInfo.Name)
			log.Printf("parent: %s -> relative: %s\n", info.Name, linkInfo.Name)
			convertUrls[link], _ = relativePath(info.Name, linkInfo.Name)
		}
	}

	log.Printf("Converter Map >> %v\n", convertUrls)
	convertLinks := func() {
		linkConverter := func(url string, isA bool) string {
			if toUrl, ok := convertUrls[url]; ok {
				return toUrl
			}
			return url
		}
		convertlinks.OfHtml(doc, linkConverter)

		// Write the new doc `html.Node` to a new temporary file
		convertHtmlFile, err := temp.File()
		if err != nil {
			log.Println(err)
			return
		}
		defer fileio.Close(convertHtmlFile)

		err = html.Render(convertHtmlFile, doc)
		if err != nil {
			log.Println(err)
			return
		}

		// successfully converted the links, and wrote the HTML node to the new temporary file,
		// move the temporary file to the actual downloaded HTML file
		fileio.Close(htmlFile)
		fileio.Close(convertHtmlFile)
		err = os.Rename(convertHtmlFile.Name(), htmlFile.Name())
		if err != nil {
			log.Println(err)
			return
		}
	}

	if a.ConvertLinks {
		convertLinks()
	}

	return
}

// FetchCss assumes the file of the given filename, is a CSS file and thus,
// extracts all linked URLs, downloads the linked resources, then optionally
// converting the links defined in the CSS file to local filesystem based paths
func (a *arg) FetchCss(mirrorUrl, fileName string) {
	// the downloaded file is CSS; attempt to extract linked url resources
	cssFile, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer fileio.Close(cssFile)
	cssStr, _ := io.ReadAll(cssFile)

	// get all urls that are linked in the given CSS file
	var linkedUrls = links.FromCssUrl(string(cssStr))

	// download each linked url synchronously, continuing to the next regardless of errors
	convertUrls := make(map[string]string)
	for _, link := range linkedUrls {
		linkedUrl, err := xurl.AbsoluteUrl(mirrorUrl, link)
		if err != nil {
			log.Println(err)
			continue
		} else if !xurl.SameHost(mirrorUrl, linkedUrl) {
			continue
		}

		linkInfo, err := a.Site(linkedUrl)
		if _, ok := a.urlDownloadInfo[linkedUrl]; errors.Is(err, ErrFileAlreadyDownloaded) && ok {
			linkInfo = a.urlDownloadInfo[linkedUrl].FileInfo
			err = nil
		}

		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("saved to -> %s\n", linkInfo.Name)
		log.Printf("parent: %s -> relative: %s\n", fileName, linkInfo.Name)

		if a.ConvertLinks {
			convertUrls[linkedUrl], _ = relativePath(fileName, linkInfo.Name)
		}
	}

	if a.ConvertLinks {
		linkConverter := func(url string) string {
			if toUrl, ok := convertUrls[url]; ok {
				return toUrl
			}
			return url
		}
		// convert all linked urls in the css, to the local files the linked urls were downloaded to
		newCss := convertlinks.OfCss(string(cssStr), linkConverter)

		// write the new CSS to a new temporary file
		convertCssFile, err := temp.File()
		if err != nil {
			log.Println(err)
			return
		}
		defer fileio.Close(convertCssFile)

		_, err = convertCssFile.WriteString(newCss)
		if err != nil {
			log.Println(err)
			return
		}

		// successfully converted the links, and wrote the new CSS to the new temporary file,
		// move the temporary file to the actual downloaded CSS file
		fileio.Close(cssFile)
		fileio.Close(convertCssFile)
		err = os.Rename(convertCssFile.Name(), cssFile.Name())
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// relativePath computes the relative path from file1 to file2
func relativePath(file1, file2 string) (string, error) {
	// Get the directory of file1
	dir1 := filepath.Dir(file1)

	// Get the absolute path of file1's directory
	absDir1, err := filepath.Abs(dir1)
	if err != nil {
		return "", err
	}

	// Get the absolute path of file2
	absFile2, err := filepath.Abs(file2)
	if err != nil {
		return "", err
	}

	// Calculate the relative path from file1's directory to file2
	relPath, err := filepath.Rel(absDir1, absFile2)
	if err != nil {
		return "", err
	}

	return relPath, nil
}
