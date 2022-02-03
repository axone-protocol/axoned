# ℹ Freely based on: https://gist.github.com/thomaspoignant/5b72d579bd5f311904d973652180c705

# Constants
BINARY_NAME             = okp4d
TARGET_FOLDER           = target
DOCKER_IMAGE_GOLANG_CI  = golangci/golangci-lint:v1.44.0
CMD_ROOT               :=./cmd/${BINARY_NAME}

# Some colors
COLOR_GREEN  = $(shell tput -Txterm setaf 2)
COLOR_YELLOW = $(shell tput -Txterm setaf 3)
COLOR_WHITE  = $(shell tput -Txterm setaf 7)
COLOR_CYAN   = $(shell tput -Txterm setaf 6)
COLOR_RESET  = $(shell tput -Txterm sgr0)

# Flags
VERSION  :=$(shell cat version)
COMMIT   :=$(shell git log -1 --format='%H')
LD_FLAGS  = -X github.com/cosmos/cosmos-sdk/version.Name=okp4d         \
		    -X github.com/cosmos/cosmos-sdk/version.ServerName=okp4d   \
		    -X github.com/cosmos/cosmos-sdk/version.ClientName=okp4d   \
		    -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		    -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)
BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

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

build-go: ## Build node executable
	@echo "${COLOR_CYAN} 🏗️ Building project ${CMD_ROOT} into ${TARGET_FOLDER}/${COLOR_RESET}"
	@go build -o ${TARGET_FOLDER}/${BINARY_NAME} ${BUILD_FLAGS} ${CMD_ROOT}

## Install:
install: ## Install node executable
	@echo "${COLOR_CYAN} 🚚 Installing project ${BINARY_NAME}${COLOR_RESET}"
	@go build ${BUILD_FLAGS} ${CMD_ROOT}

## Start:
start: install ## Start the blockchain node
	@echo "${COLOR_CYAN} 🚀 Starting project ${COLOR_RESET}"
	@okp4d start

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
