
.PHONY: build
build:
#	COMMIT=$(shell git rev-parse --short HEAD)
#	VERSION=$(shell echo "v1.0")
#	go build -v -ldflags="-X 'main.version=${VERSION}'" -o cnvrgctl cmd/cnvrgctl/*.go
	go build -v -o cnvrgctl cmd/cnvrgctl/*.go

.PHONY: install
install: build
	mv cnvrgctl /usr/local/bin
	cnvrgctl completion bash > /usr/local/etc/bash_completion.d/cnvrgctl

.PHONY: build-linux
build-linux:
	docker run --rm -v ${PWD}:/usr/src/cnvrgctl -w /usr/src/cnvrgctl golang:1.14 /bin/bash -c "GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cnvrgctl cmd/cnvrgctl/*.go"

