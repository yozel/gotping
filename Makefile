GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=gotping

define cross_build_func
	CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) $(GOBUILD) -o bin/$(BINARY_NAME)-$(1)-$(2) -v
endef

all: test build

test: 
	$(GOTEST) -v .

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v

build-cross: test build-linux build-darwin build-windows

build-linux:
	$(call cross_build_func,linux,amd64)
	$(call cross_build_func,linux,386)

build-darwin:
	$(call cross_build_func,darwin,amd64)
	$(call cross_build_func,darwin,386)

build-windows:
	$(call cross_build_func,windows,amd64)
	$(call cross_build_func,windows,386)
