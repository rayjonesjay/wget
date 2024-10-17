# WGET

WGET is a utility that recreates some of the core functionalities of the original GNU Wget using Go. It is designed for non-interactive downloads from the web and includes several features such as downloading single or multiple files, limiting download speed, and mirroring entire websites.

## Project Structure

```plaintext
.
├── args
│   ├── eval_args.go
│   └── eval_args_test.go
├── ctx
│   ├── ctx.go
│   └── ctx_test.go
├── downloader
│   └── downloader.go
├── fetch
│   ├── fetch.go
│   └── fetch_test.go
├── fileio
│   ├── fileio.go
│   └── fileio_test.go
├── go.mod
├── help
│   └── help.go
├── httpx
│   ├── extras.go
│   └── extras_test.go
├── LICENSE
├── limitedio
│   ├── limitedio.go
│   └── limitedio_test.go
├── main.go
├── mirror
│   └── mirror.go
├── README.md
├── syscheck
│   ├── system.go
│   └── system_test.go
├── TESTS.md
├── xerr
│   └── xerr.go
└── xurl
    ├── xurl.go
    └── xurl_test.go

13 directories, 25 files
```

## Features

1. **Download a file via URL**  
   Command-line argument accepts a URL to download a file from the web.
   ```bash
   $ go run . URL
   ```

2. **Download and save under a different name**  
   You can specify a custom file name using the `-O` flag.
   ```bash
   $ go run . -O=file_name URL
   ```

3. **Download and save to a different directory**  
   Use `-P` flag to specify a directory for saving the file.
   ```bash
   $ go run . -P=path/to/save URL
   ```

4. **Limit download speed**  
   Control the download speed using the `--rate-limit` flag. The suffix `k` and `M` are used for kilobytes and megabytes, respectively.
   ```bash
   $ go run . --rate-limit=100k URL
   $ go run . --rate-limit=200M URL
   ```

5. **Download in background**  
   The `--background` flag allows the download to proceed in the background.
   ```bash
   $ go run . --background URL
   ```

6. **Download multiple files asynchronously**  
   The `-i` flag reads a file containing multiple URLs and downloads them concurrently.
   ```bash
   $ go run . -i=path/to/file/with/links
   ```

7. **Mirror an entire website**  
   The `--mirror` flag downloads an entire website for offline use.
   ```bash
   $ go run . --mirror URL
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

Visit the [Go official website]("https://golang.org/dl/) and download the latest version suitable for your operating system. Follow the installation instructions for your OS.

2. **Verify Go Installation** 

After installation, verify that Go is properly set up by running the following command:

```bash
$ go version
```

You should see the installed version of Go in the output.


### Setup and Installation 

1. **Clone the Repository**
  
  Clone the wget repository to your local machine:
  ```bash
   $ git clone https://learn.zone01kisumu.ke/git/ramuiruri/wget
```
2. **Navigate to the Project directory

Change into the cloned repository directory and install necessary dependencies
```bash
$ cd wget
$ go mod tidy
```

### Running the utility

Now that the project is set up you can start using the wget utility

####  Download a Single File
```bash
$ go run . https://example.com/file.zip
```

#### Download and Save with a Specific Name
```bash
$ go run . -O=newfile.zip https://example.com/file.zip
```

#### Download to a Specific Directory
```bash
$ go run . -P=~/Downloads/ https://example.com/file.zip
```

#### Limit Download Speed
```bash
$ go run . --rate-limit=500k https://example.com/largefile.zip
```

#### Download Multiple Files Asynchronously
```bash
$ go run . -i=links.txt
```

#### Mirror a Website
```bash
$ go run . --mirror https://example.com
```

## Contribution

We welcome contributions to improve this project! If you wish to contribute:

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

