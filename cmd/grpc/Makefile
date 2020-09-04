.PHONY=build up

GO111MODULE?=on
CGO_ENABLED?=0
GOOS?=linux
GO_BIN?=app
GO?=go

.EXPORT_ALL_VARIABLES:


build:
	$(GO) build -a -installsuffix nocgo -o $(GO_BIN) main.go
