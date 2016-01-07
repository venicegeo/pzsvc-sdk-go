[![GoDoc](https://godoc.org/github.com/venicegeo/pzsvc-sdk-go?status.svg)](https://godoc.org/github.com/venicegeo/pzsvc-sdk-go)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/venicegeo/pzsvc-sdk-go/blob/master/LICENSE)

# pzsvc-sdk-go

This may in the end be an awful name, as I'm not really setting out to create a full-fledged SDK. What it is, for now, is a place for common components that I'm using in my Go-based Piazza services. Things like structs defining JSON messages for job creation, S3 upload/download helper utilities, etc.

## Install

`pzsvc-sdk-go` uses [Glide](https://github.com/Masterminds/glide) to manage its dependencies. Assuming you are on a Mac OS X, Glide can be easily installed via [Homebrew](https://github.com/Homebrew/homebrew) (alternative installation instruction can be found on the Glide webpage).

```console
$ brew install glide
```

We also make use of [Go 1.5's vendor/ experiment](https://medium.com/@freeformz/go-1-5-s-vendor-experiment-fd3e830f52c3#.ueuy8ao53), so you'll need to make sure you are running Go 1.5+.

Installing `pzsvc-sdk-go` is as simple as cloning the repo, installing dependencies, and installing the package.

```console
$ git clone https://github.com/venicegeo/pzsvc-sdk-go
$ cd pzsvc-sdk-go
$ glide install
$ go install
```
