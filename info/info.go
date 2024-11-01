// Package info contains various information about the program, made available to the program at runtime
package info

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	// Org the name of the organization owning the program
	Org     = "org.zone01.wget"
	Name    = "wget"
	Version = "1.2.0"
	Go      = "1.22.4"
)

// VersionText returns the program version information
func VersionText() string {
	v := `
Zone01 %s %s built on %s-%s
Runtime: Go v%s

+http +https

Copyright (C) 2024 Zone01 Kisumu.
License MIT: The MIT License
https://opensource.org/license/mit.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

authors:
https://github.com/rayjonesjay
https://github.com/Wambita
https://github.com/WycliffeAlphus
https://github.com/nanyona
`
	v = strings.TrimLeft(v, "\n")
	return fmt.Sprintf(v, Name, Version, runtime.GOOS, runtime.GOARCH, Go)
}
