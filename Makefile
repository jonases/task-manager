
all: test

clean:
	rm -f task-manager

install: prepare
	godep go install

prepare:
	go get github.com/tools/godep

build: prepare
	env GOOS="linux" GOARCH="amd64" godep go build

test: prepare build
	echo "no tests"

.PHONY: install prepare build test
