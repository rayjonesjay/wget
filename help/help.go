package help

const (
	// UsageMessage is used to Display usage message to user on how to use our program
	// it will be triggered with the --help flag
	UsageMessage = `wget: missing URL
Usage: $ ./wget [OPTION] [URL]

Try './wget --help' for more options.`

	Manual = `wget 1.0.0, a non-interactive retriever.
Usage: ./wget [OPTION] [URL]

Supported options as of version 1.0.0:

FLAG:				 USAGE:						EXPLANATION:

--help				 ******						print this manual
-B 	  		        -B URL					        allow download in the background after startup
-v 			        ******					        display the current version of wget and exit
-O   				-O=FILE	 URL					log messages to FILE
--mirror  			--mirror URL					mirror a website
-P			        -P=PATH  URL				        specify path to save downloaded resource
--rate-limit=NS									set speed limit, N is a number, S is either (k or M)
-i  				-i=FILE						allow download of resource by reading links stored in a file


bug reports, questions, issues to https://github.com/rayjonesjay
`

	Version = "1.0.0"
)
