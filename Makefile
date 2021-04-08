
pack:
	pkger

.PHONY: build-mac
build-mac: pack
	go build -v -o bin/cnvrgctl-darwin-x86_64 main.go pkged.go

.PHONY: install
install: build
	mv cnvrgctl /usr/local/bin
	cnvrgctl completion bash > /usr/local/etc/bash_completion.d/cnvrgctl

.PHONY: build-linux
build-linux:
	docker run --rm -v ${PWD}:/usr/src/cnvrgctl -w /usr/src/cnvrgctl golang:1.14 /bin/bash -c "go get github.com/markbates/pkger/cmd/pkger && pkger && GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/cnvrgctl-linux-x86_64 main.go pkged.go"

get-tools:
	curl -Lso pkg/assets/rke https://github.com/rancher/rke/releases/download/v1.2.7/rke_linux-amd64
	mkdir -p tmp && pushd tmp
	curl -Lso k9s_Linux_x86_64.tar.gz https://github.com/derailed/k9s/releases/download/v0.24.7/k9s_Linux_x86_64.tar.gz
	tar zxvf k9s_Linux_x86_64.tar.gz && cp ./k9s ../pkg/assets/k9s
	popd && rm -fr tmp
	curl -Lso pkg/assets/kubectl https://dl.k8s.io/release/v1.20.5/bin/linux/amd64/kubectl

build-all: build-linux build-mac

upload:
	aws s3 cp ./bin/ s3://whitening-pn38xqkin816s3fk --recursive
