#!/usr/bin/make -f

test: fmt
	go test -timeout=1s -count=1 ./...

fmt:
	go mod tidy
	goimports -w `find . -type f -name '*.go'`
