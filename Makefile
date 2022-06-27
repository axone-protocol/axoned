# ℹ Freely based on: https://gist.github.com/thomaspoignant/5b72d579bd5f311904d973652180c705

# Constants
BINARY_NAME             = okp4d
TARGET_FOLDER           = target
DIST_FOLDER             = $(TARGET_FOLDER)/dist
DOCKER_IMAGE_GOLANG_CI  = golangci/golangci-lint:v1.44.0
CMD_ROOT               :=./cmd/${BINARY_NAME}

# Some colors
COLOR_GREEN  = $(shell tput -Txterm setaf 2)
COLOR_YELLOW = $(shell tput -Txterm setaf 3)
COLOR_WHITE  = $(shell tput -Txterm setaf 7)
COLOR_CYAN   = $(shell tput -Txterm setaf 6)
COLOR_RESET  = $(shell tput -Txterm sgr0)

BUILD_TAGS += netgo
BUILD_TAGS := $(strip $(BUILD_TAGS))

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

LD_FLAGS := $(strip $(LD_FLAGS))

BUILD_FLAGS := -tags "$(BUILD_TAGS)" -ldflags '$(LD_FLAGS)' -trimpath

# Commands
GO_BUiLD := go build $(BUILD_FLAGS)

# Environments
ENVIRONMENTS = \
	darwin-amd64 \
	darwin-arm64 \
	linux-amd64 \
	windows-amd64
ENVIRONMENTS_TARGETS = $(addprefix build-go-, $(ENVIRONMENTS))

.PHONY: all lint lint-go build build-go help

all: help

## Lint:
lint: lint-go ## Lint all available linters

lint-go: ## Lint go source code
	@echo "${COLOR_CYAN}🔍 Inspecting go source code${COLOR_RESET}"
	@docker run --rm \
  		-v `pwd`:/app:ro \
  		-w /app \
  		${DOCKER_IMAGE_GOLANG_CI} \
  		golangci-lint run -v

## Build:
build: build-go ## Build all available artefacts (executable, docker image, etc.)

build-go: ## Build node executable for the current environment (default build)
	@echo "${COLOR_CYAN} 🏗️ Building project ${COLOR_RESET}${CMD_ROOT}${COLOR_CYAN}${COLOR_RESET} into ${COLOR_YELLOW}${DIST_FOLDER}${COLOR_RESET}"
	@$(call build-go,"","",${DIST_FOLDER}/${BINARY_NAME})

build-go-all: $(ENVIRONMENTS_TARGETS) ## Build node executables for all available environments

$(ENVIRONMENTS_TARGETS):
	@GOOS=$(word 3, $(subst -, ,$@)); \
    GOARCH=$(word 4, $(subst -, ,$@)); \
    FOLDER=${DIST_FOLDER}/$$GOOS/$$GOARCH; \
    if [ $$GOOS = "windows" ]; then \
      EXTENSION=".exe"; \
    fi; \
    FILENAME=$$FOLDER/${BINARY_NAME}$$EXTENSION; \
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
