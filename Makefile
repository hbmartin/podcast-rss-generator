SHELL = /bin/bash

GITHUB_REPO := hbmartin/podcast-rss-generator
MODULE := github.com/$(GITHUB_REPO)/v2

.PHONY: build test cover lint fmt vet tidy clean README

build:
	go build ./...

test:
	go test -race ./...

cover:
	go test -coverprofile=profile.out ./...
	go tool cover -func=profile.out

lint:
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy
	go mod vendor

README:
	godoc2ghmd -play -ex -verify_import_links=0 $(MODULE) > README.md.tmp
	echo "[![Go Reference](https://pkg.go.dev/badge/$(MODULE).svg)](https://pkg.go.dev/$(MODULE))" > README.md
	echo "[![CI](https://github.com/$(GITHUB_REPO)/actions/workflows/ci.yml/badge.svg)](https://github.com/$(GITHUB_REPO)/actions/workflows/ci.yml)" >> README.md
	echo "[![codecov](https://codecov.io/gh/$(GITHUB_REPO)/branch/master/graph/badge.svg)](https://codecov.io/gh/$(GITHUB_REPO))" >> README.md
	echo "[![Go Report Card](https://goreportcard.com/badge/$(MODULE))](https://goreportcard.com/report/$(MODULE))" >> README.md
	echo "[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)" >> README.md
	echo  >>README.md
	cat README.md.tmp >> README.md
	rm README.md.tmp

clean:
	rm -rf corpus crashers suppressions workdir podcast-fuzz.zip profile.out coverage.txt
