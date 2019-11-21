#!/usr/bin/make -f

test: fmt
	go test -v ./...

fmt:
	goimports -d -w $(find . -type f -name '*.go')
