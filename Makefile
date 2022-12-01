# ℹ Freely based on: https://gist.github.com/thomaspoignant/5b72d579bd5f311904d973652180c705

# Constants
BINARY_NAME             = okp4d
TARGET_FOLDER           = target
DIST_FOLDER             = $(TARGET_FOLDER)/dist
RELEASE_FOLDER          = $(TARGET_FOLDER)/release
DOCKER_IMAGE_GOLANG		= golang:1.19-alpine3.16
DOCKER_IMAGE_GOLANG_CI  = golangci/golangci-lint:v1.49
DOCKER_IMAGE_BUF  		= okp4/buf-cosmos:0.3.1
DOCKER_BUILDX_BUILDER   = okp4-builder
CMD_ROOT               :=./cmd/${BINARY_NAME}
LEDGER_ENABLED ?= true

# Some colors
COLOR_GREEN  = $(shell tput -Txterm setaf 2)
COLOR_YELLOW = $(shell tput -Txterm setaf 3)
COLOR_WHITE  = $(shell tput -Txterm setaf 7)
COLOR_CYAN   = $(shell tput -Txterm setaf 6)
COLOR_RED    = $(shell tput -Txterm setaf 1)
COLOR_RESET  = $(shell tput -Txterm sgr0)

BUILD_TAGS += netgo
BUILD_TAGS := $(strip $(BUILD_TAGS))
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

# Flags
WHITESPACE := $(subst ,, )
COMMA := ,
BUILD_TAGS_COMMA_SEP := $(subst $(WHITESPACE),$(COMMA),$(BUILD_TAGS))
VERSION  := $(shell cat version)
COMMIT   := $(shell git log -1 --format='%H')
LD_FLAGS  = \
	-X github.com/cosmos/cosmos-sdk/version.Name=okp4d         \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=okp4d   \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=okp4d   \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
    -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(BUILD_TAGS_COMMA_SEP)

ifeq ($(LINK_STATICALLY),true)
	LD_FLAGS += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif

LD_FLAGS := $(strip $(LD_FLAGS))

BUILD_FLAGS := -tags "$(BUILD_TAGS)" -ldflags '$(LD_FLAGS)' -trimpath

# Commands
GO_BUiLD := go build $(BUILD_FLAGS)

# Environments
ENVIRONMENTS = \
	darwin-amd64 \
	darwin-arm64 \
	linux-amd64 \
	linux-arm64
ENVIRONMENTS_TARGETS = $(addprefix build-go-, $(ENVIRONMENTS))

# Release binaries
RELEASE_BINARIES = \
	linux-amd64 \
	linux-arm64
RELEASE_TARGETS = $(addprefix release-binary-, $(RELEASE_BINARIES))

.PHONY: all lint lint-go build build-go help

all: help

## Lint:
lint: lint-go lint-proto ## Lint all available linters

lint-go: ## Lint go source code
	@echo "${COLOR_CYAN}🔍 Inspecting go source code${COLOR_RESET}"
	@docker run --rm \
  		-v `pwd`:/app:ro \
  		-w /app \
  		${DOCKER_IMAGE_GOLANG_CI} \
  		golangci-lint run -v

lint-proto: ## Lint proto files
	@echo "${COLOR_CYAN}🔍️ lint proto${COLOR_RESET}"
	@docker run --rm \
		-v ${HOME}/.cache:/root/.cache \
  		-v `pwd`:/proto \
  		-w /proto \
  		${DOCKER_IMAGE_BUF} \
  		lint proto -v

## Format:
format: format-go ## Run all available formatters

format-go: ## Format go files
	@echo "${COLOR_CYAN}📐 Formatting go source code${COLOR_RESET}"
	@docker run --rm \
  		-v `pwd`:/app:rw \
  		-w /app \
  		${DOCKER_IMAGE_GOLANG} \
  		sh -c \
		"go install mvdan.cc/gofumpt@v0.4.0; gofumpt -w -l ."

## Build:
build: build-go ## Build all available artefacts (executable, docker image, etc.)

build-go: ## Build node executable for the current environment (default build)
	@echo "${COLOR_CYAN} 🏗️ Building project ${COLOR_RESET}${CMD_ROOT}${COLOR_CYAN}${COLOR_RESET} into ${COLOR_YELLOW}${DIST_FOLDER}${COLOR_RESET}"
	@$(call build-go,"","",${DIST_FOLDER}/${BINARY_NAME})

build-go-all: $(ENVIRONMENTS_TARGETS) ## Build node executables for all available environments

$(ENVIRONMENTS_TARGETS):
	@GOOS=$(word 3, $(subst -, ,$@)); \
    GOARCH=$(word 4, $(subst -, ,$@)); \
    if [ $$GOARCH = "amd64" ]; then \
    	TARGET_ARCH="x86_64"; \
    elif [ $$GOARCH = "arm64" ]; then \
    	TARGET_ARCH="aarch64"; \
    fi; \
    HOST_OS=`uname -s | tr A-Z a-z`; \
    HOST_ARCH=`uname -m`; \
    if [ $$GOOS != $$HOST_OS ] || [ $$TARGET_ARCH != $$HOST_ARCH ]; then \
      echo "${COLOR_RED} ❌ Cross compilation impossible${COLOR_RESET}" >&2; \
      exit 1; \
    fi; \
    FOLDER=${DIST_FOLDER}/$$GOOS/$$GOARCH; \
    FILENAME=$$FOLDER/${BINARY_NAME}; \
	echo "${COLOR_CYAN} 🏗️ Building project ${COLOR_RESET}${CMD_ROOT}${COLOR_CYAN} for environment ${COLOR_YELLOW}$$GOOS ($$GOARCH)${COLOR_RESET} into ${COLOR_YELLOW}$$FOLDER${COLOR_RESET}" && \
	$(call build-go,$$GOOS,$$GOARCH,$$FILENAME)


## Install:
install: ## Install node executable
	@echo "${COLOR_CYAN} 🚚 Installing project ${BINARY_NAME}${COLOR_RESET}"
	@go install ${BUILD_FLAGS} ${CMD_ROOT}

## Test:
test: test-go ## Pass all the tests

test-go: build ## Pass the test for the go source code
	@echo "${COLOR_CYAN} 🧪 Passing go tests${COLOR_RESET}"
	@go test -v -covermode=count -coverprofile ./target/coverage.out ./...

## Clean:
clean: ## Remove all the files from the target folder
	@echo "${COLOR_CYAN} 🗑 Cleaning folder $(TARGET_FOLDER)${COLOR_RESET}"
	@rm -rf $(TARGET_FOLDER)/

## Proto:
proto-format: ## Format Protobuf files
	@echo "${COLOR_CYAN} 📐 Formatting Protobuf files${COLOR_RESET}"
	@docker run --rm \
    		-v ${HOME}/.cache:/root/.cache \
    		-v `pwd`:/proto \
    		-w /proto \
    		${DOCKER_IMAGE_BUF} \
    		format -w -v

proto-build: ## Build all Protobuf files
	@echo "${COLOR_CYAN} 🔨️Build Protobuf files${COLOR_RESET}"
	@docker run --rm \
		-v ${HOME}/.cache:/root/.cache \
		-v `pwd`:/proto \
		-w /proto \
		${DOCKER_IMAGE_BUF} \
		build proto -v

proto-gen: proto-build ## Generate all the code from the Protobuf files
	@echo "${COLOR_CYAN} 📝 Generating code from Protobuf files${COLOR_RESET}"
	@docker run --rm \
		-v ${HOME}/.cache:/root/.cache \
		-v `pwd`:/proto \
		-w /proto \
		${DOCKER_IMAGE_BUF} \
		generate proto --template buf.gen.proto.yaml -v
	@cp -r github.com/okp4/okp4d/x/* x/
	@rm -rf github.com

## Documentation:
doc-proto: proto-gen ## Generate the documentation from the Protobuf files
	@for MODULE in $(shell find proto -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq | xargs dirname) ; do \
		echo "${COLOR_CYAN} 📖 Generate documentation for $${MODULE} module${COLOR_RESET}" ; \
  		docker run --rm \
        		-v ${HOME}/.cache:/root/.cache \
        		-v `pwd`:/proto \
        		-w /proto \
        		${DOCKER_IMAGE_BUF} \
        		generate --path $${MODULE} --template buf.gen.doc.yaml -v ; \
        mv docs/proto/docs.md docs/$${MODULE}.md ; \
	done

## Release:
release-assets: release-binary-all release-checksums ## Generate release assets

release-binary-all: $(RELEASE_TARGETS)

$(RELEASE_TARGETS): ensure-buildx-builder
	@GOOS=$(word 3, $(subst -, ,$@)); \
    GOARCH=$(word 4, $(subst -, ,$@)); \
    BINARY_NAME="okp4d-${VERSION}-$$GOOS-$$GOARCH"; \
	echo "${COLOR_CYAN} 🎁 Building ${COLOR_GREEN}$$GOOS $$GOARCH ${COLOR_CYAN}release binary${COLOR_RESET} into ${COLOR_YELLOW}${RELEASE_FOLDER}${COLOR_RESET}"; \
	docker buildx use ${DOCKER_BUILDX_BUILDER}; \
	docker buildx build \
		--platform $$GOOS/$$GOARCH \
		-t $$BINARY_NAME \
		--load \
		.; \
	mkdir -p ${RELEASE_FOLDER}; \
	docker rm -f tmp-okp4d || true; \
	docker create -ti --name tmp-okp4d $$BINARY_NAME; \
	docker cp tmp-okp4d:/usr/bin/okp4d ${RELEASE_FOLDER}/$$BINARY_NAME; \
	docker rm -f tmp-okp4d; \
	tar -zcvf ${RELEASE_FOLDER}/$$BINARY_NAME.tar.gz ${RELEASE_FOLDER}/$$BINARY_NAME;

release-checksums:
	@echo "${COLOR_CYAN} 🐾 Generating release binary checksums${COLOR_RESET} into ${COLOR_YELLOW}${RELEASE_FOLDER}${COLOR_RESET}"
	@rm ${RELEASE_FOLDER}/sha256sum.txt; \
	for asset in `ls ${RELEASE_FOLDER}`; do \
		shasum -a 256 ${RELEASE_FOLDER}/$$asset >> ${RELEASE_FOLDER}/sha256sum.txt; \
	done;

ensure-buildx-builder:
	@echo "${COLOR_CYAN} 👷 Ensuring docker buildx builder${COLOR_RESET}"
	@docker buildx ls | sed '1 d' | cut -f 1 -d ' ' | grep -q ${DOCKER_BUILDX_BUILDER} || \
	docker buildx create --name ${DOCKER_BUILDX_BUILDER}

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${COLOR_YELLOW}make${COLOR_RESET} ${COLOR_GREEN}<target>${COLOR_RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${COLOR_YELLOW}%-20s${COLOR_GREEN}%s${COLOR_RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${COLOR_CYAN}%s${COLOR_RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
	@echo ''
	@echo 'This Makefile depends on ${COLOR_CYAN}docker${COLOR_RESET}. To install it, please follow the instructions:'
	@echo '- for ${COLOR_YELLOW}macOS${COLOR_RESET}: https://docs.docker.com/docker-for-mac/install/'
	@echo '- for ${COLOR_YELLOW}Windows${COLOR_RESET}: https://docs.docker.com/docker-for-windows/install/'
	@echo '- for ${COLOR_YELLOW}Linux${COLOR_RESET}: https://docs.docker.com/engine/install/'

# Build go executable
# $1: operating system (GOOS)
# $2: architecture (GOARCH)
# $3: filename of the executable generated
define build-go
	GOOS=$1 GOARCH=$2 $(GO_BUiLD) -o $3 ${CMD_ROOT}
endef
