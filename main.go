package main

import (
	"fmt"
	"os"

	"wget/syscheck"
)

func main() {
	if err := syscheck.CheckOperatingSystem(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	// download.DownloadUrl()
}
