.PHONY=build test protoc

GO111MODULE?=on
CGO_ENABLED?=0
GOOS?=linux
GO_BIN?=/go/bin/app
GO?=go
MAIN_DIR?=/var/app
APP_NAME?=grpc
PROTOC?=protoc

.EXPORT_ALL_VARIABLES:


test:
	$(GO) test ./... -v -covermode=count -coverprofile=coverage.out -cover

build:
	$(MAKE) -C cmd/$(APP_NAME) build

protoc:
	$(PROTOC) --proto_path=pkg/ --go_opt=paths=source_relative --go_out=plugins=grpc:pkg/ status/update.proto status/status.proto
