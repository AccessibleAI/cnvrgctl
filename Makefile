
.PHONY: build
build-mac:
	go build -v -o bin/cnvrgctl-darwin-x86_64 cmd/cnvrgctl/*.go

.PHONY: install
install: build
	mv cnvrgctl /usr/local/bin
	cnvrgctl completion bash > /usr/local/etc/bash_completion.d/cnvrgctl

.PHONY: build-linux
build-linux:
	docker run --rm -v ${PWD}:/usr/src/cnvrgctl -w /usr/src/cnvrgctl golang:1.14 /bin/bash -c "GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/cnvrgctl-linux-x86_64 cmd/cnvrgctl/*.go"

build-all: build-linux build-mac

upload:
	aws s3 cp ./bin/ s3://whitening-pn38xqkin816s3fk --recursive
