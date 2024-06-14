ARCH?=arm64
OS?=darwin

.PHONY: go-mod-vendor
go-mod-vendor:
	go mod vendor

.PHONY: install
install:
	GOOS=$(OS) GOARCH=$(ARCH) go build -o ./bin/sparrow cmd/sparrow/sparrow.go

