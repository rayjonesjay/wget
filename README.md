# WGET

WGET is a utility that is used for non-interactive downloads over the net.

## Project Structure
```
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

## Functionalities
1. Downloading a file given an URL-(Uniform Resource Locator) parsed through the command line.

For example:
```
$ go run . URL
```

2. Downloading a single file and saving it under a different name.

For example:
```
$ go run . -O=file_name URL
```

3. Downloading and saving the file in a different path.

For example:
```
$ go run . -P=path/to/save URL
```

4. Set the download speed, limiting the rate speed of a download. Only **k** and **M** for kilo bytes and Mega bytes respectively are allowed if k or M not used as suffix then the value is assumed as bytes per second.

For example:
```
$ go run . --rate-limit=100k URL
$ go run . --rate-limit=200M URL
```

5. Downloading file in background.

For example:
```
$ go run . --background URL
```

6. Downloading multiple files at same time, reading a file containing multiple download links asynchronously.

For example:
```
$ go run . -i=path/to/file/with/links
```

7. Main Feature will be to download an entire website also known as mirror.

For example:
```
$ go run . --mirror URL
```


### Authors
* [**ramuiruri**](https://learn.zone01kisumu.ke/git/ramuiruri)

* [najwang](https://learn.zone01kisumu.ke/git/najwang)

* [wyonyango](https://learn.zone01kisumu.ke/git/wyonyango)

* [shfana](https://learn.zone01kisumu.ke/git/shfana)