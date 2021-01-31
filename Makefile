
.PHONY: build
build:
	go build -o cnvrgctl cmd/cnvrgctl/*.go

.PHONY: install
install: build
	mv cnvrgctl /usr/local/bin
	cnvrgctl completion bash > /usr/local/etc/bash_completion.d/cnvrgctl
