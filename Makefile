.PHONY: test

export PATH := bin:$(PATH)

staticcheck := /home/runner/go/bin/staticcheck

all: test

cidep:
	go get -u honnef.co/go/tools/cmd/staticcheck
	go get -u github.com/kisielk/errcheck

test: gotest gostaticcheck goerrcheck

gotest:
	go test ./... -v -timeout=45s -failfast

gostaticcheck:
	staticcheck ./...

goerrcheck:
	errcheck -exclude .lint/errcheck-excludes -ignoretests ./...

build:
	go build -o bin/section section.go