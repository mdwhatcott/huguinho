#!/usr/bin/make -f

test: fmt
	go test -v -timeout=1s -count=1 ./...

fmt:
	goimports -d -w `find . -type f -name '*.go'`
