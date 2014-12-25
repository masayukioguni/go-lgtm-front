all: clean test

setup:
	go get github.com/golang/lint/golint
	go get github.com/jteeuwen/go-bindata/...

test:
	go test $(TESTFLAGS) ./...

clean:
	go clean

templates:
	go-bindata 


lint:
	golint ./...

vet:
	go vet ./...

.PHONY: setup test clean lint vet