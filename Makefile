#!/usr/bin/make -f

test: fmt
	go test -cover -timeout=1s -count=1 ./...

fmt:
	go mod tidy
	/Users/mike/bin/goimports -w `find . -type f -name '*.go'`

install: test
	go install github.com/mdwhatcott/huguinho/cmd/...
