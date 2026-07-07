SHELL = /bin/bash

GITHUB_REPO:=eduncan911/podcast

build:
	go build ./...

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

README:
	godoc2ghmd -play -ex -verify_import_links=0 github.com/$(GITHUB_REPO) > README.md.tmp
	echo "[![GoDoc](https://godoc.org/github.com/$(GITHUB_REPO)?status.svg)](https://godoc.org/github.com/$(GITHUB_REPO))" > README.md	
	echo "[![Build Status](https://github.com/$(GITHUB_REPO)/workflows/CI/badge.svg)](https://github.com/$(GITHUB_REPO)/actions?workflow=CI)" >> README.md
	echo "[![Coverage Status](https://codecov.io/gh/$(GITHUB_REPO)/branch/master/graph/badge.svg)](https://codecov.io/gh/$(GITHUB_REPO))" >> README.md
	echo "[![Go Report Card](https://goreportcard.com/badge/github.com/$(GITHUB_REPO))](https://goreportcard.com/report/github.com/$(GITHUB_REPO))" >> README.md
	echo "[![MIT License](https://img.shields.io/npm/l/mediaelement.svg)](https://eduncan911.mit-license.org/)" >> README.md
	echo  >>README.md
	cat README.md.tmp >> README.md
	rm README.md.tmp

clean:
	rm -rf corpus crashers suppressions workdir podcast-fuzz.zip
