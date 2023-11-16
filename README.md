# url-manager

## Test exam

A console application for fetch url from file

## Install

### Go

```sh
go install github.com/robotomize/url-manager/cmd/url-manager@latest
```

### Docker

```sh
cd tests
docker pull robotomize/urlmanager:latest
docker run --rm -v $(pwd):/app  -w /app robotomize/urlmanager urlmanager -s ./fixtures.txt
```

## Usage

```shell
Check websites console application

Usage:
  url-manager [flags]

Flags:
  -d, --debug           debug logging
  -h, --help            help for url-manager
  -s, --source string   source file with urls
  -c, --sync            sync mode with one thread
```