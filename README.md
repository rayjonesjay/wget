# WGET
WGET is a utility that is used for non-interactive downloads over the net.

## Project Structure
```
.
├── background
│   └── background.go
├── directory
│   └── savetodir.go
├── download
│   ├── download.go
│   └── download_test.go
├── errors
│   └── errors.go
├── LICENSE
├── mirror
│   └── mirror.go
├── multiple
│   └── multiple.go
├── README.md
└── speed
    └── speed.go

8 directories, 10 files
```

## Functionalities
1. wget downloads a file given an URL-(Uniform Resource Locator).

For example:
```
$ go run . https://learn.zone01kisumu.ke/git/root/public/raw/branch/master/subjects/ascii-art/standard.txt
```

2. Downloading a single file and saving it under a different name

For example:
```
$ go run . URL <fileName>
```

3. Downloading and saving the file in a different directory.

For example:
```
$ go run . URL <DirectoryName>
```

4. Set the download speed, limiting the rate speed of a download

For example:
```
$ go run . --speed=100 URL
```

5. Downloading file in background.

For example:
```
$ go run . --background URL
```

6. Downloading multiple files at same time, reading a file containing multiple download links asynchronously.

For example:
```
$ go run . <file_with_url.txt>
```

7. Main Feature will be to download an entire website also known as mirror.

For example:
```
$ go run . <websiteLink>
```
