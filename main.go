package main

import (
	"wget/download"
)

func main() {
	url := "https://learn.zone01kisumu.ke/git/root/public/src/branch/master/subjects/ascii-art/"
	file := "standard.txt"
	download.DownloadUrl(url + file)
}
