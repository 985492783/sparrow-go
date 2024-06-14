ARCH?=arm64
OS?=darwin
VERSION?=1.0.0

.PHONY: go-mod-vendor
go-mod-vendor:
	go mod vendor

.PHONY: install
install:
	GOOS=$(OS) GOARCH=$(ARCH) go build -o ./bin/sparrow cmd/sparrow/sparrow.go

.PHONY: install_docker
install_docker:
	docker build -t sparrow/sparrow:$(VERSION) \
	--build-arg ARCH=$(ARCH) \
	--platform linux/${ARCH} \
	.

.PHONY: start_docker
start_docker:
	docker run -d -it \
	-p 9854:9854 \
	-p 9800:9800 \
	sparrow/sparrow:$(VERSION)