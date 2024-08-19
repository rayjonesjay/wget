// types contains user defined types and constants
package types

// Global Variable
var (
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

	DefaultDownloadDirectory = "$HOME/Downloads"

	// default port used by http if not specified
	DefaultHTTPPort = "80"

	// default port used by https if not specified
	DefaultHTTPSPort = "443"

)

const (
	// Size of Data in KilobBytes
	KB = 1000 * 1

	// Size of Data in KibiBytes, same as 2^10
	KiB = 1 << (10 * 1)

	// Size of Data in MegaBytes
	MB = 1000 * KB

	// Size of Data in MebiBytes, same as 2^20
	MiB = 1 << 20

	// Size of Data in GigaBytes
	GB = 1000 * MB

	// Size of Data in GibiBytes, same as 2^30
	GiB = 1 << 30

	// Size of Data in TeraBytes
	TB = 1000 * GB

	// Size of Data in TebiBytes, same as 2^40
	TiB = 1 << 40
)
