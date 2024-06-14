FROM golang:1.22.3-alpine
ARG ARCH

RUN mkdir -p /sparrow
COPY ./ /sparrow
WORKDIR /sparrow
EXPOSE 9854
EXPOSE 9800

RUN GOOS=$(OS) GOARCH=$(ARCH) go build -o ./bin/sparrow cmd/sparrow/sparrow.go
CMD ["./bin/sparrow"]