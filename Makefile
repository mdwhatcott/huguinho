#!/usr/bin/make -f

test: fmt
	go test ./...

fmt:
	goimports -d -w `find . -type f -name '*.go'`
