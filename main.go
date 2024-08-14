package main

import (
	"fmt"
	"os"

	"wget/download"
	"wget/errors"
	"wget/syscheck"
)

func main() {
	if err := syscheck.CheckOperatingSystem(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	download.DownloadUrl(errors.Link1)
}
