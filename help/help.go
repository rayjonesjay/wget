package help

import "strings"

var (
	// UsageMessage is used to Display usage message to user on how to use our program
	// it will be triggered with the --help flag
	UsageMessage = `wget: missing URL
Usage: $ ./wget [OPTION] [URL]

Try './wget --help' for more options.`

	Version = "1.1.0"
)

func PrintManPage() string {
	man := `
    Usage: ./wget <link>

    ┌─────────────────────┬──────────────────────────────────────────────┬─────────────────────────────────────────────────────────────┐
    │ FLAG                │ USAGE                                        │ EXPLANATION                                                 │
    ├─────────────────────┼──────────────────────────────────────────────┼─────────────────────────────────────────────────────────────┤
    │ --help              │ --help                                       │ print this manual                                           │
    │ -B                  │ -B URL                                       │ allow download in the background after startup              │
    │ -v                  │ -v | --version                               │ display the current version of wget and exit                │
    │ -O                  │ -O=FILE URL                                  │ log messages to FILE                                        │
    │ --mirror            │ --mirror URL                                 │ mirror a website                                            │
    │ -P                  │ -P=PATH URL                                  │ specify path to save downloaded resource                    │
    │ --rate-limit=NS     │ --rate-limit=NS URL                          │ set speed limit, N is a number, S is either (k or M)        │
    │ -i                  │ -i=FILE                                      │ allow download of resource by reading links stored in a file│
    └─────────────────────┴──────────────────────────────────────────────┴─────────────────────────────────────────────────────────────┘

    Bug reports, questions, issues to:
    - https://github.com/rayjonesjay
    - https://github.com/Wambita
    - https://github.com/WycliffeAlphus
    - https://github.com/nanyona
	
    Ouote By Maya Angelou - "Nothing will work unless you do"
    `
	man = strings.TrimLeft(man, "\n")
	return man
}
