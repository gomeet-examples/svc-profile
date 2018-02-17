DEFAULT_PREFIXES=gomeet-svc-,svc-
SVC_PREFIX=svc-
NAME = $(SVC_PREFIX)profile
GO_PACKAGE_BASE = github.com/gomeet-examples
GOMEET_GO_PACKAGE_NAME = github.com/gomeet/gomeet
GO_PACKAGE_NAME = $(GO_PACKAGE_BASE)/$(NAME)
GO_PROTO_PACKAGE_ALIAS = pb
GO_PROTO_PACKAGE = $(GO_PACKAGE_NAME)/$(GO_PROTO_PACKAGE_ALIAS)
VERSION = $(shell cat VERSION)

OS_NAME=$(shell go env GOOS)
OS_ARCH=$(shell go env GOARCH)

DB_TYPES=mysql

DOCKER_TAG = $(shell cat VERSION | tr +- __)
DOCKER_IMAGE_NAME = gomeetexamples/$(NAME)
DOCKER_REGISTRY?=docker.io
DOCKER_NETWORK = network-grpc-gomeetexamples
DOCKER_SVC_CONTAINER = svc-$(NAME)
DOCKER_CONSOLE_CONTAINER = console-$(NAME)

SUB_SERVICES=
# SUB-SERVICES DEFINITION : make-tag-docker-compose
# DOCKER_TAG_SVC_{{SubServiceNameUpperSnakeCase}} = $(shell cat ./vendor/github.com/gomeet-examples/svc-{{SubServiceNameKebabCase}}/VERSION | tr +- __)
# END SUB-SERVICES DEFINITION : make-tag-docker-compose

PACKAGE_DIR = _build/packaged
PACKAGE_PROTO_NAME = proto.tar.gz

PROTOC_VERSION=3.5.1
PROTOC_REPO_URL=https://github.com/google/protobuf/releases/download/v$(PROTOC_VERSION)
PROTOC_BIN=_tools/bin/protoc

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
ifeq ($(OS_NAME),openbsd)
  CMD_SHASUM = sha256 -r
endif

# release arguments mangement
# usage :
#    make release <Git flow option : start|finish> <Release version : major|minor|patch> [Release version metadata (optional)]
ifeq (release,$(firstword $(MAKECMDGOALS)))
  RELEASE_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  ifneq ($(filter $(firstword $(RELEASE_ARGS)),start finish),)
    ifneq ($(filter $(word 2,$(RELEASE_ARGS)),major minor patch),)
        $(eval $(RELEASE_ARGS):;@:)
      else
        $(error unknow release version - usage : "make release <Git flow option : start|finish> <Release version : major|minor|patch> [Release version metadata (optional)]")
    endif
  else
    $(error unknow release state - usage : "make release <Git flow option : start|finish> <Release version : major|minor|patch> [Release version metadata (optional)]")
  endif
endif

# dep-update-gomeetexamples arguments mangement
# usage :
#    make dep-update-gomeetexamples <SUB_SERVICES_NAME>
ifeq (dep-update-gomeetexamples,$(firstword $(MAKECMDGOALS)))
  ASKED_SVC := $(word 2,$(MAKECMDGOALS))
  ifeq ($(ASKED_SVC),)
    SUB_SVC := $(ALL_SVC),
  else
    ifneq ($(filter $(ASKED_SVC),$(ALL_SVC)),)
      SUB_SVC := $(ASKED_SVC)
      $(eval $(ASKED_SVC):;@:)
    else
      $(error unknow $(ASKED_SVC) in ($(ALL_SVC)))
    endif
  endif
endif

define package_builder
	@echo "$(NAME): build $(1) - $(2) on $(3)"
	-mkdir -p $(3)
	cp VERSION $(3)/VERSION
	cp LICENSE $(3)/LICENSE
	cp CHANGELOG.md $(3)/CHANGELOG.md
	$(eval FILE_NAME := $(shell if [ "$(1)" = "windows" ]; then echo "$(NAME).exe"; else echo "$(NAME)"; fi))
	GOOS=$(1) GOARCH=$(2) \
		CGO_ENABLED=0 go build \
			-ldflags '-extldflags "-lm -lstdc++ -static"' \
			-ldflags "-X $(GO_PACKAGE_NAME)/service.version=$(VERSION) -X $(GO_PACKAGE_NAME)/service.name=$(NAME)" \
			-o $(3)/$(FILE_NAME) \
		main.go
	cd $(3) && $(CMD_SHASUM) $(FILE_NAME) > SHA256SUM
	@echo ""
endef

define package_shasum
	@echo "$(NAME): package_shasum"
	cd $(PACKAGE_DIR) \
	  &&	find */* '!' -name 'VERSION' '!' -name 'LICENSE' '!' -name '*SUM' '!' -name '*SUMS' \
		| xargs $(CMD_SHASUM) > SHA256SUMS
	cd $(PACKAGE_DIR) && [ -f "$(PACKAGE_PROTO_NAME)" ] && $(CMD_SHASUM) $(PACKAGE_PROTO_NAME) >> SHA256SUMS
endef

.PHONY: build
build: proto
	@echo "$(NAME): build task"
	-mkdir -p _build
	CGO_ENABLED=0 go build \
		-ldflags '-extldflags "-lm -lstdc++ -static"' \
		-ldflags "-X $(GO_PACKAGE_NAME)/service.version=$(VERSION) -X $(GO_PACKAGE_NAME)/service.name=$(NAME)" \
		-o _build/$(NAME) \
	main.go

.PHONY: clean
clean: tools-clean package-clean proto-clean
	@echo "$(NAME): clean task"

.PHONY: proto
proto: tools proto-gen-doc
	@echo "$(NAME): proto task"
	PATH=$(shell pwd)/_tools/bin:${PATH} && \
		for f in $(GO_PROTO_PACKAGE_ALIAS)/*.proto; do \
			protoc -I. \
				-I$(shell pwd)/third_party \
				--gogo_out=plugins="grpc:${GOPATH}/src" \
				--govalidators_out="${GOPATH}/src" \
				--grpc-gateway_out="logtostderr=true:${GOPATH}/src" \
				--swagger_out="logtostderr=true:." \
				--gomeetfaker_out="${GOPATH}/src" \
				--gomeet-service_out="project_pkg=${GO_PACKAGE_NAME};db_types=${DB_TYPES};default_prefixes=${DEFAULT_PREFIXES};sub_services=${SUB_SERVICES}:${GOPATH}/src" \
			$$f; \
		echo "$(NAME): compiled proto file - " $$f; \
		echo ""; \
	done
	cd $(GO_PROTO_PACKAGE_ALIAS) && go generate .

.PHONY: proto-gen-doc
proto-gen-doc:
	@echo "$(NAME): proto task"
	PATH=$(shell pwd)/_tools/bin:${GOPATH}/src/github.com/gomeet/gomeet/_build:${PATH} && \
		for f in $(GO_PROTO_PACKAGE_ALIAS)/*.proto; do \
			for t in "docbook,$(NAME).docbook" "html,$(NAME).html" "markdown,$(NAME).md"; do \
				protoc -I. \
					-I$(shell pwd)/third_party \
					--plugin=$(shell pwd)/_tools/bin/protoc-gen-doc \
					--doc_out=./docs/grpc-services/ \
					--doc_opt="$$t:google/*,github.com/*" \
				$$f; \
			done; \
		echo "$(NAME): gen doc files from " $$f; \
		echo ""; \
	done
	-mv docs/grpc-services/$(NAME).md docs/grpc-services/README.md
	-mv docs/grpc-services/$(NAME).html docs/grpc-services/index.html

.PHONY: proto-clean
proto-clean:
	@echo "$(NAME): proto-clean task"
	-rm $(GO_PROTO_PACKAGE_ALIAS)/*.pb.* $(GO_PROTO_PACKAGE_ALIAS)/*.swagger.json

.PHONY: package
package: clean proto package-proto docker
	@echo "$(NAME): package task"
	$(call package_builder,linux,amd64,$(PACKAGE_DIR)/linux-amd64)
	$(call package_builder,linux,arm,$(PACKAGE_DIR)/linux-arm32)
	$(call package_builder,linux,arm64,$(PACKAGE_DIR)/linux-arm64)
	$(call package_builder,openbsd,amd64,$(PACKAGE_DIR)/openbsd-amd64)
	$(call package_builder,darwin,amd64,$(PACKAGE_DIR)/darwin-amd64)
	$(call package_builder,windows,amd64,$(PACKAGE_DIR)/windows-amd64)

	cp VERSION $(PACKAGE_DIR)/VERSION
	cp LICENSE $(PACKAGE_DIR)/LICENSE
	cp CHANGELOG.md $(PACKAGE_DIR)/CHANGLOG.md
	$(call package_shasum)

.PHONY: package-clean
package-clean:
	@echo "$(NAME): package-clean task"
	-rm -r $(PACKAGE_DIR)

.PHONY: package-proto
package-proto:
	@echo "$(NAME): package-proto task"
	cp -r third_party $(PACKAGE_DIR)/proto
	cp $(GO_PROTO_PACKAGE_ALIAS)/*.proto $(PACKAGE_DIR)/proto/
	cp VERSION $(PACKAGE_DIR)/proto/
	cd $(PACKAGE_DIR) && tar -C ./proto/ -czvf $(PACKAGE_PROTO_NAME) . && rm -rf ./proto/
	$(call package_shasum)

.PHONY: release
release: tools
	$(eval RELEASE_META := $(word 3, $(RELEASE_ARGS)))
	$(eval RELEASE_META_FULL := $(if $(RELEASE_META),"+$(RELEASE_META)",""))
	$(eval RELEASE_VERSION := $(shell if [ "$(word 2,$(RELEASE_ARGS))" = "patch" ]; then echo "`sed 's/+dev//g' VERSION`$(RELEASE_META_FULL)" ; else _tools/bin/semver -$(word 2,$(RELEASE_ARGS)) -build "$(RELEASE_META)" $(VERSION); fi))
	echo "$(NAME): new $(word 2,$(RELEASE_ARGS)) release -> $(RELEASE_VERSION)"
	git flow release start "v$(RELEASE_VERSION)"
	echo "$(RELEASE_VERSION)" > VERSION
	git add VERSION
	git commit -m "Bump version - v$(RELEASE_VERSION)"
	awk \
		-v \
		log_title="## Unreleased\n\n- Nothing\n\n## $(RELEASE_VERSION) - $$(date +%Y-%m-%d)" \
		'{gsub(/## Unreleased/,log_title)}1' \
		CHANGELOG.md > CHANGELOG.md.tmp && \
			mv CHANGELOG.md.tmp CHANGELOG.md
	git add CHANGELOG.md
	git commit -m "Improved CHANGELOG.md"
	# TODO don't push binaries in git repository but use github release process
	#@$(MAKE) package
	#@$(MAKE) docker-push
	#git add _build/packaged/
	#git commit -m "Added v$(RELEASE_VERSION) packages"
	#git add .env
	#git commit -m "Added docker-compose .env"
	git flow release publish "v$(RELEASE_VERSION)"
ifeq (finish,$(firstword $(RELEASE_ARGS)))
	git flow release finish "v$(RELEASE_VERSION)"
	$(eval DEV_RELEASE_VERSION := $(shell _tools/bin/semver -patch -build "dev" $(RELEASE_VERSION)))
	echo "$(DEV_RELEASE_VERSION)" > VERSION
	git add VERSION
	git commit -m "Bump version - $(DEV_RELEASE_VERSION)"
	# TODO don't push binaries in git repository
	#@$(MAKE) package
	#@$(MAKE) docker-push
	#git add .env
	#git add $(PACKAGE_DIR)
	#git commit -m "Added v$(DEV_RELEASE_VERSION) packages"
	git push --tag
	git push origin develop
	git push origin master
endif

.PHONY: docker
docker:
	echo "TAG=$(DOCKER_TAG)" > .env
	# SUB-SERVICES DEFINITION : make-tag-docker-compose-to-env
	# echo "TAG_SVC_{{SubServiceNameUpperSnakeCase}}=$(DOCKER_TAG_SVC_{{SubServiceNameUpperSnakeCase}})" >> .env
	# END SUB-SERVICES DEFINITION : make-tag-docker-compose-to-env
	-docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .

.PHONY: docker-push
docker-push: docker
	-docker tag $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_TAG)
	-docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

.PHONY: start
start: docker
	docker-compose up -d
	@echo "\n\
	Grafana:\n\
		http://localhost:3000/\n\
		For admin password see GF_SECURITY_ADMIN_PASSWORD in $(shell pwd)/infra/config.monitoring file\n\
	Prometheus:\n\
		http://localhost:9090/\n\
	Alert manager:\n\
		http://localhost:9093/\n\
	Swagger docs:\n\
		http://localhost:8000/docs\n\
	Godoc:\n\
		http://localhost:8001/pkg/github.com/gomeet-examples/$(NAME)\n\
	Documentation: \n\
		$(shell pwd)/docs/docker/README.md \n\
		$(shell pwd)/docs/docker-compose/README.md\n\
		$(shell pwd)/docs/grafana/README.md\n\
		"

.PHONY: stop
stop:
	docker-compose down
	@echo "\nTo remove volume use :\n        docker volume prune -f"

.PHONY: install
install:
	go install .

.PHONY: dep
dep: tools
	_tools/bin/dep ensure

.PHONY: dep-prune
dep-prune: tools
	_tools/bin/dep prune

.PHONY: dep-update-gomeetexamples
dep-update-gomeetexamples: tools
	@echo "dep-update-gomeetexamples $(ASKED_SVC)"
	$(foreach svc, $(SUB_SVC), _tools/bin/dep ensure -update "$(GO_PACKAGE_BASE)/$(SVC_PREFIX)$(svc)";)

.PHONY: dep-update-gomeet-utils
dep-update-gomeet-utils: tools
	@echo "dep-update-gomeet-utils $(ASKED_SVC)"
	_tools/bin/dep ensure -update "$(GOMEET_GO_PACKAGE_NAME)"

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

.PHONY: tools-upgrade-gomeet
tools-upgrade-gomeet: tools
	@echo "dep-update-gomeet-utils $(ASKED_SVC)"
	_tools/bin/retool upgrade github.com/gomeet/gomeet-tools-markdown-server origin/master
	_tools/bin/retool upgrade github.com/gomeet/go-proto-gomeetfaker/protoc-gen-gomeetfaker origin/master
	_tools/bin/retool upgrade github.com/gomeet/gomeet/protoc-gen-gomeet-service origin/master
	_tools/bin/retool upgrade github.com/gomeet/gomeet/gomeet origin/master

.PHONY: tools-upgrade
tools-upgrade: tools
	GOPATH=$(shell pwd)/_tools/ && \
		for tool in $(shell cat tools.json | grep "Repository" | awk '{print $$2}' | sed 's/,//g' | sed 's/"//g' ); do $$GOPATH/bin/retool upgrade $$tool origin/master ; done

.PHONY: test
test: build
	cd service && go test
	if [ -f $(GO_PROTO_PACKAGE_ALIAS)/*_test.go ]; then cd pb && go test; fi
	_build/$(NAME) functest -e --random-port

.PHONY: docker-test
docker-test: docker
	docker run -v $(shell pwd):/go/src/$(GO_PACKAGE_NAME) --rm golang:1.8.3-alpine3.6 /bin/sh -c "cd /go/src/$(GO_PACKAGE_NAME)/service && go test"
	docker run --rm $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) functest -e

.PHONY: doc-server
doc-server: tools
	_tools/bin/gomeet-tools-markdown-server

.PHONY: gomeet-regenerate-project
gomeet-regenerate-project: tools
	_tools/bin/gomeet new ${GO_PACKAGE_NAME} \
		--default-prefixes=${DEFAULT_PREFIXES} \
		--proto-name=${GO_PROTO_PACKAGE_ALIAS} \
		--sub-services=${SUB_SERVICES} \
		--db-types=${DB_TYPES} \
		--force

