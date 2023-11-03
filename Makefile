IMG ?= quay.io/skeeey/mochimqtt:latest

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)


## build
.PHONY: build
build:
	go build -o bin/server-source cmd/server/source/main.go
	go build -o bin/server cmd/server/main.go
	go build -o bin/client cmd/client/main.go

.PHONY: image
image:
	docker build -t ${IMG} .
