# Copyright 2016 Michal Witkowski. All Rights Reserved.
# See LICENSE for licensing terms.

PROTOC_VERSION=3.5.1
PROTOC_REPO_URL=https://github.com/google/protobuf/releases/download/v$(PROTOC_VERSION)
PROTOC_BIN=_tools/bin/protoc

NAME = protoc-gen-gomeetfaker

OS_ARCH=$(shell go env GOARCH)
OS_NAME=$(shell go env GOOS)

PROTOC_VERSION=3.5.1
PROTOC_REPO_URL=https://github.com/google/protobuf/releases/download/v$(PROTOC_VERSION)

ifeq ($(OS_NAME),windows)
  PROTOC_PKG_NAME := protoc-$(PROTOC_VERSION)-win32.zip
endif
ifeq ($(OS_NAME),darwin)
  ifeq ($(OS_ARCH),amd64)
    PROTOC_PKG_NAME := protoc-$(PROTOC_VERSION)-osx-x86_64.zip
  else
    PROTOC_PKG_NAME := protoc-$(PROTOC_VERSION)-osx-x86_32.zip
  endif
endif
ifeq ($(OS_NAME),linux)
  ifeq ($(OS_ARCH),arm64)
    PROTOC_PKG_NAME := protoc-$(PROTOC_VERSION)-linux-aarch_64.zip
  else
    ifeq ($(OS_ARCH),amd64)
      PROTOC_PKG_NAME := protoc-$(PROTOC_VERSION)-linux-x86_64.zip
    else
      PROTOC_PKG_NAME := protoc-$(PROTOC_VERSION)-linux-x86_32.zip
    endif
  endif
endif

CMD_SHASUM = shasum -a 256
ifeq ($(UNAME_S),OpenBSD)
	CMD_SHASUM = sha256 -r
endif

.PHONY: build
build: proto
	@echo "$(NAME): build task"
	-mkdir -p _build
	CGO_ENABLED=0 go build \
		-ldflags '-extldflags "-lm -lstdc++ -static"' \
		-o _build/$(NAME) \
	protoc-gen-gomeetfaker/main.go

.PHONY: clean
clean: tools-clean proto-clean
	@echo "$(NAME): clean task"
	-rm _build/protoc-gen-gomeetfaker

.PHONY: proto
proto: tools
	@echo "$(NAME): proto task"
	PATH=$(shell pwd)/_tools/bin:${PATH} && \
		protoc -I. \
			-I$(shell pwd)/third_party \
			--gogo_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:. \
		gomeetfaker.proto

.PHONY: proto-examples
proto-examples: build
	@echo "$(NAME): Regenerating example directory"
	PATH=$(shell pwd)/_tools/bin:$(shell pwd)/_build:${PATH} && \
		protoc -I. \
			-I${GOPATH}/src \
			-I$(shell pwd)/third_party \
			--go_out=. \
			--gomeetfaker_out=. \
		examples/full/pb/*.proto
	PATH=$(shell pwd)/_tools/bin:$(shell pwd)/_build:${PATH} && \
		protoc -I. \
			-I${GOPATH}/src \
			-I$(shell pwd)/third_party \
			--go_out=. \
			--gomeetfaker_out=. \
		examples/address_book/pb/*.proto
	PATH=$(shell pwd)/_tools/bin:$(shell pwd)/_build:${PATH} && \
		protoc -I. \
			-I${GOPATH}/src \
			-I$(shell pwd)/third_party \
			--go_out=. \
			--gomeetfaker_out=. \
		examples/book/pb/*.proto

.PHONY: proto-clean
proto-clean:
	@echo "$(NAME): proto-clean task"
	-rm gomeetfaker.pb.go examples/**/pb/*.pb.go

.PHONY: tools
tools:
	@echo "$(NAME): tools task"
ifeq ("$(wildcard _tools/src/github.com/twitchtv/retool)","")
	$(MAKE) tools-sync-retool
endif
	GOPATH=$(shell pwd)/_tools/ && \
		go install github.com/twitchtv/retool
	_tools/bin/retool build
ifeq ("$(wildcard $(PROTOC_BIN))","")
	$(MAKE) tools-sync-protoc
endif

.PHONY: tools-clean
tools-clean:
	@echo "$(NAME): tools-clean task"
	-rm -rf _tools/bin _tools/pkg _tools/manifest.json _tools/protoc

.PHONY: tools-sync
tools-sync: tools-sync-retool tools-sync-protoc
tools-sync:
	@echo "$(NAME): tools-sync task"

.PHONY: tools-sync-retool
tools-sync-retool:
	@echo "$(NAME): tools-sync-retool task"
	GOPATH=$(shell pwd)/_tools/ && \
		go get github.com/twitchtv/retool && \
		go install github.com/twitchtv/retool
	_tools/bin/retool sync

.PHONY: tools-sync-protoc
tools-sync-protoc:
	@echo "$(NAME): tools-sync-protoc task"
	@rm -rf _tools/protoc
	@mkdir -p _tools/protoc
	@mkdir -p _tools/bin
	@curl -L -o _tools/protoc/$(PROTOC_PKG_NAME) $(PROTOC_REPO_URL)/$(PROTOC_PKG_NAME)
	@cd _tools/protoc && unzip $(PROTOC_PKG_NAME)
	@cp _tools/protoc/bin/protoc $(PROTOC_BIN)
	@cp -r _tools/protoc/include/google/protobuf/* third_party/google/protobuf/

.PHONY: tools-upgrade
tools-upgrade: tools
	GOPATH=$(shell pwd)/_tools/ && \
		for tool in $(shell cat tools.json | grep "Repository" | awk '{print $$2}' | sed 's/,//g' | sed 's/"//g' ); do $$GOPATH/bin/retool upgrade $$tool origin/master ; done

.PHONY: dep
dep: tools
	_tools/bin/dep ensure

