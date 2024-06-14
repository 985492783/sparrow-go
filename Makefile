
.PHONY: go-mod-vendor
go-mod-vendor:
	go mod vendor

.PHONY: install
install:
	go build -o ./bin/sparrow cmd/sparrow/sparrow.go

