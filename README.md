# WGET

WGET is a utility that recreates some of the core functionalities of the original GNU Wget using Go. It is designed for non-interactive downloads from the web and includes several features such as downloading single or multiple files, limiting download speed, and mirroring entire websites.

## Table of Contents

- [Project Structure](#project-structure)
- [Features](#features)
- [Flags](#flags)
- [Usage](#usage)
  - [Prerequisites](#prerequisites)
  - [Setup and Installation](#setup-and-installation)
  - [Running the Utility](#running-the-utility)
- [Contribution](#contribution)
-[Authors](#authors)
- [License](#license)

## Project Structure

```plaintext
.
├── args
│   ├── eval_args.go
│   └── eval_args_test.go
├── convertlinks
│   ├── convertlinks.go
│   └── convertlinks_test.go
├── css
│   ├── css.go
│   └── css_test.go
├── ctx
│   ├── ctx.go
│   └── ctx_test.go
├── downloader
│   ├── downloader.go
│   └── downloader_test.go
├── fetch
│   ├── fetch.go
│   └── fetch_test.go
├── fileio
│   ├── fileio.go
│   └── fileio_test.go
├── globals
│   ├── globals.go
│   └── globals_test.go
├── go.mod
├── go.sum
├── help
│   ├── help.go
│   └── help_test.go
├── httpx
│   ├── extras.go
│   └── extras_test.go
├── info
│   ├── info.go
│   └── info_test.go
├── LICENSE
├── limitedio
│   ├── limitedio.go
│   └── limitedio_test.go
├── main.go
├── mirror
│   ├── dirlimits.go
│   ├── links
│   │   ├── css.go
│   │   ├── css_test.go
│   │   ├── html.go
│   │   └── html_test.go
│   ├── mirror.go
│   ├── path.go
│   ├── path_test.go
│   ├── README.md
│   └── xurl
│       ├── url.go
│       └── url_test.go
├── README.md
├── syscheck
│   ├── system.go
│   └── system_test.go
├── temp
│   ├── temp.go
│   └── temp_test.go
├── xerr
│   ├── xerr.go
│   └── xerr_test.go
└── xurl
    ├── xurl.go
    └── xurl_test.go

20 directories, 48 files
```

## Features

To use the program first build it using, ```go build -o wget```

1. **Download a file via URL**  
   Command-line argument accepts a URL to download a file from the web.
   ```bash
   $ ./wget URL
   ```

2. **Download and save under a different name**  
   You can specify a custom file name using the `-O` flag.
   ```bash
   $ ./wget -O=file_name URL
   ```

3. **Download and save to a different directory**  
   Use `-P` flag to specify a directory for saving the file.
   ```bash
   $ ./wget -P=path/to/save URL
   ```

4. **Limit download speed**  
   Control the download speed using the `--rate-limit` flag. The suffix `k` and `M` are used for kilobytes and megabytes, respectively.
   ```bash
   $ ./wget --rate-limit=100k URL
   $ ./wget --rate-limit=200M URL
   ```

5. **Download in background**  
   The `--background` flag allows the download to proceed in the background.
   ```bash
   $ ./wget --background URL
   ```

6. **Download multiple files asynchronously**  
   The `-i` flag reads a file containing multiple URLs and downloads them concurrently.
   ```bash
   $ ./wget -i=path/to/file/with/links
   ```

7. **Mirror an entire website**  
   The `--mirror` flag downloads an entire website for offline use.
   ```bash
   $ ./wget --mirror URL
   ```

## Flags

Here are the available flags for the WGET utility:

- `-O`: Specify the output file name for the downloaded file.
- `-P`: Specify the directory where the file should be saved.
- `--rate-limit`: Limit the download speed. Use `k` for kilobytes and `M` for megabytes.
- `-i`: Download multiple files by reading URLs from a file.
- `--mirror`: Mirror an entire website.
- `-B`: Download in the background and save logs to `wget-log`.
- `--background`: Download in the background (similar to `-B`).

## Usage

### Prerequisites

Before using this WGET utility, ensure that Go is installed on your system. If not, follow these steps:

1. **Install Go**  
   Visit the [Go official website](https://golang.org/dl/) and download the latest version suitable for your operating system. Follow the installation instructions for your OS.

2. **Verify Go Installation**  
   After installation, verify that Go is properly set up by running the following command:
   ```bash
   $ go version
   ```
   You should see the installed version of Go in the output.

### Setup and Installation

1. **Clone the Repository**  
   Clone this WGET repository to your local machine:
   ```bash
   $ git clone https://learn.zone01kisumu.ke/git/ramuiruri/wget.git
   ```

2. **Navigate into the Project Directory**  
   Change into the cloned repository directory:
   ```bash
   $ cd wget
   ```

3. **Install Dependencies**  
   Run the following command to install any necessary dependencies:
   ```bash
   $ go mod tidy
   ```

### Running the Utility

Now that the project is set up, you can start using the WGET utility:


#### See the help page
```bash
$ ./wget -h
```

#### See the current version
```bash
$ ./wget -v
```

#### Download using url
```bash
$ ./wget https://example.com/file.zip
```

#### Download and Save with a Specifhttps://learn.zone01kisumu.ke/git/ramuiruri/wgetic Name
```bash
$ ./wget -O=newfile.zip https://example.com/file.zip
```

#### Download to a Specific Directory
```bash
$ ./wget -P=~/Downloads/ https://example.com/file.zip
```

#### Limit Download Speed
```bash
$ ./wget --rate-limit=500k https://example.com/largefile.zip
```

#### Download Multiple Files Asynchronously
```bash
$ ./wget -i=links.txt
```

#### Mirror a Website
```bash
$ ./wget --mirror https://example.com
```

## Contribution

We welcome contributions to improve this project! If you wish to contribute:
Contact [Ray Jones](https://github.com/rayjonesjay).

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -am 'Add new feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Create a pull request.

Make sure your code adheres to the coding standards and passes all tests before submitting the pull request.


### Authors
* [**ramuiruri**](https://learn.zone01kisumu.ke/git/ramuiruri)

* [najwang](https://learn.zone01kisumu.ke/git/najwang)

* [wyonyango](https://learn.zone01kisumu.ke/git/wyonyango)

* [shfana](https://learn.zone01kisumu.ke/git/shfana)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
