package help

import (
	"fmt"
	"math/rand"
	"strings"
	"wget/info"
)

var (
	// UsageMessage is used to Display usage message to user on how to use our program
	// it will be triggered with the --help flag
	UsageMessage = `wget: missing URL
Usage: wget [OPTION] [URL]

Try 'wget --help' for more options.`
)

// PrintManPage prints the program's help text
func PrintManPage() string {
	intro := fmt.Sprintf("Zone01 Wget %s, a non-interactive network retriever.\n", info.Version)
	man := `
Usage: wget [OPTION]... [URL]...

Mandatory arguments to long options are mandatory for short options too.

    ┌──────────────────────────┬────────────────────────────────────────────────────────────────────┐
    │ OPTION                   │ EXPLANATION                                                        │
    ├──────────────────────────┼────────────────────────────────────────────────────────────────────┤
    │ --help                   │ print this manual and exit                                         │
    │ -v | --version           │ display the current version of wget and exit                       │
    │ -B                       │ download a file immediately to the background,                     │
    │                          │ redirecting the output to the log file (wget-log)                  │
    │ -O=FILENAME              │ download a file and save it under a different name                 │
    │ -P=PATH                  │ specify the path where to save downloaded resource                 │
    │ --rate-limit=AMOUNT      │ Limit the download speed to AMOUNT bytes per second.               │
    │                          │ Amount may be expressed in bytes, kilobytes with the ‘k’ suffix,   │ 
    │                          │ or megabytes with the ‘M’ suffix. For example,                     │
    │                          │ ‘--limit-rate=20k’ will limit the retrieval rate to 20KiB/s        │                         
    │ -i=FILE                  │ Read URLs from the local file                                      │
    │ --mirror                 │ mirror a website                                                   │
    │ --reject | -R=list       │ list of file suffixes to avoid downloading during the retrieval    │
    │ --exclude | -X=list      │ list of directories excluded from the download                     │
    │ --convert-links          │ convert the links in the document,                                 │
    │                          │ to make them suitable for local viewing                            │
    └──────────────────────────┴────────────────────────────────────────────────────────────────────┘

    Bug reports, questions, issues to:
    - https://github.com/rayjonesjay
    - https://github.com/Wambita
    - https://github.com/WycliffeAlphus
    - https://github.com/nanyona

    `
	type quote struct {
		Text   string
		Author string
	}
	var quotes = []quote{
		{"Programs must be written for people to read, and only incidentally for machines to execute.", "Harold Abelson"},
		{"Any fool can write code that a computer can understand. Good programmers write code that humans can understand.", "Martin Fowler"},
		{"First, solve the problem. Then, write the code.", "John Johnson"},
		{"Experience is the name everyone gives to their mistakes.", "Oscar Wilde"},
		{"In order to be irreplaceable, one must always be different.", "Coco Chanel"},
		{"Java is to JavaScript what car is to Carpet.", "Chris Heilmann"},
		{"Knowledge is power.", "Francis Bacon"},
		{"Sometimes it pays to stay in bed on Monday, rather than spending the rest of the week debugging Monday's code.", "Dan Salomon"},
		{"Code is like humor. When you have to explain it, it’s bad.", "Cory House"},
		{"Fix the cause, not the symptom.", "Steve Maguire"},
		{"Simplicity is the soul of efficiency.", "Austin Freeman"},
		{"Before software can be reusable it first has to be usable.", "Ralph Johnson"},
		{"Make it work, make it right, make it fast.", "Kent Beck"},
		{"Programs are meant to be read by humans and only incidentally for computers to execute.", "Donald Knuth"},
		{"The best method for accelerating a computer is the one that boosts it by 9.8 m/s².", "Anonymous"},
		{"Walking on water and developing software from a specification are easy if both are frozen.", "Edward V. Berard"},
		{"If debugging is the process of removing software bugs, then programming must be the process of putting them in.", "Edsger Dijkstra"},
		{"Software is a great combination of artistry and engineering.", "Bill Gates"},
		{"Measuring programming progress by lines of code is like measuring aircraft building progress by weight.", "Bill Gates"},
		{"The most important property of a program is whether it accomplishes the intention of its user.", "C.A.R. Hoare"},
		{"The trouble with programmers is that you can never tell what a programmer is doing until it’s too late.", "Seymour Cray"},
		{"The function of good software is to make the complex appear to be simple.", "Grady Booch"},
		{"Adding manpower to a late software project makes it later.", "Fred Brooks"},
		{"The best performance improvement is the transition from the nonworking state to the working state.", "J. Osterhout"},
		{"One of my most productive days was throwing away 1000 lines of code.", "Ken Thompson"},
		{"The best way to get a project done faster is to start sooner.", "Jim Highsmith"},
		{"When to use iterative development? You should use iterative development only on projects that you want to succeed.", "Martin Fowler"},
		{"Controlling complexity is the essence of computer programming.", "Brian Kernighan"},
		{"Deleted code is debugged code.", "Jeff Sickel"},
		{"Testing can only prove the presence of bugs, not their absence.", "Edsger Dijkstra"},
	}
	var randomQuote string
	{ // Get a random quote
		randomIndex := rand.Intn(len(quotes))
		helpQuote, author := quotes[randomIndex].Text, quotes[randomIndex].Author
		randomQuote = fmt.Sprintf("Quote By %s:\n      %q", author, helpQuote)
	}

	man = strings.TrimLeft(man, "\n")
	return intro + man + randomQuote
}
