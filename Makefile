# ℹ Freely based on: https://gist.github.com/thomaspoignant/5b72d579bd5f311904d973652180c705

# Constants
BINARY_NAME             = axoned
TARGET_FOLDER           = target
DIST_FOLDER             = $(TARGET_FOLDER)/dist
RELEASE_FOLDER          = $(TARGET_FOLDER)/release
CMD_ROOT               :=./cmd/${BINARY_NAME}
LEDGER_ENABLED         ?= true

# Docker images
DOCKER_IMAGE_GOLANG		  = golang:1.23-alpine3.20
DOCKER_IMAGE_GOLANG_CI    = golangci/golangci-lint:v1.61
DOCKER_IMAGE_PROTO        = ghcr.io/cosmos/proto-builder:0.14.0
DOCKER_IMAGE_BUF          = bufbuild/buf:1.4.0
DOCKER_PROTO_RUN          := docker run --rm --user $(id -u):$(id -g) -v $(HOME)/.cache:/root/.cache -v $(PWD):/workspace --workdir /workspace $(DOCKER_IMAGE_PROTO)
DOCKER_BUF_RUN            := docker run --rm -v $(HOME)/.cache:/root/.cache -v $(PWD):/workspace --workdir /workspace $(DOCKER_IMAGE_BUF)
DOCKER_BUILDX_BUILDER     = axone-builder
DOCKER_IMAGE_MARKDOWNLINT = thegeeklab/markdownlint-cli:0.32.2
DOCKER_IMAGE_GOTEMPLATE   = hairyhenderson/gomplate:v3.11.3-alpine

# Some colors
COLOR_GREEN  = $(shell tput -Txterm setaf 2)
COLOR_YELLOW = $(shell tput -Txterm setaf 3)
COLOR_WHITE  = $(shell tput -Txterm setaf 7)
COLOR_CYAN   = $(shell tput -Txterm setaf 6)
COLOR_RED    = $(shell tput -Txterm setaf 1)
COLOR_RESET  = $(shell tput -Txterm sgr0)

# Blockchain constants
CHAIN     		:= localnet
CHAIN_HOME    	:= ./target/deployment/${CHAIN}
CHAIN_MONIKER 	:= local-node
CHAIN_BINARY 	:= ./${DIST_FOLDER}/${BINARY_NAME}

DAEMON_NAME 	:= axoned
DAEMON_HOME 	:= `pwd`/${CHAIN_HOME}

BUILD_TAGS += netgo
BUILD_TAGS := $(strip $(BUILD_TAGS))
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      BUILD_TAGS += ledger
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
        BUILD_TAGS += ledger
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
    -X github.com/cosmos/cosmos-sdk/version.AppName=axoned      \
	-X github.com/cosmos/cosmos-sdk/version.Name=axoned         \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=axoned   \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=axoned   \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)   \
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

# Handle sed -i on Darwin
SED_FLAG=
SHELL_NAME := $(shell uname -s)
ifeq ($(SHELL_NAME),Darwin)
    SED_FLAG := ""
endif

.PHONY: all
all: help

## Lint:
.PHONY: lint
lint: lint-go lint-proto ## Lint all available linters

.PHONY: lint-go
lint-go: ## Lint go source code
	@echo "${COLOR_CYAN}🔍 Inspecting go source code${COLOR_RESET}"
	@docker run --rm \
  		-v `pwd`:/app:ro \
  		-w /app \
  		${DOCKER_IMAGE_GOLANG_CI} \
  		golangci-lint run -v

.PHONY: lint-proto
lint-proto: ## Lint proto files
	@echo "${COLOR_CYAN}🔍️ lint proto${COLOR_RESET}"
	@$(DOCKER_BUF_RUN) lint proto -v

## Format:
.PHONY: format
format: format-go ## Run all available formatters

.PHONY: format-go
format-go: ## Format go files
	@echo "${COLOR_CYAN}📐 Formatting go source code${COLOR_RESET}"
	@docker run --rm \
  		-v `pwd`:/app:rw \
  		-w /app \
  		${DOCKER_IMAGE_GOLANG} \
  		sh -c \
		"go install mvdan.cc/gofumpt@v0.7.0; gofumpt -w -l ."

.PHONY: format-proto
format-proto: ## Format proto files
	@echo "${COLOR_CYAN}📐 Formatting proto files${COLOR_RESET}"
	@$(DOCKER_BUF_RUN) format -w

## Build:
.PHONY: build
build: build-go ## Build all available artefacts (executable, docker image, etc.)

.PHONY: build-go
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
.PHONY: install
install: ## Install node executable
	@echo "${COLOR_CYAN} 🚚 Installing project ${BINARY_NAME}${COLOR_RESET}"
	@go install ${BUILD_FLAGS} ${CMD_ROOT}

## Test:
.PHONY: test
test: test-go ## Pass all the tests

.PHONY: test-go
test-go: build ## Pass the test for the go source code
	@echo "${COLOR_CYAN} 🧪 Passing go tests${COLOR_RESET}"
	@go test -v -coverprofile ./target/coverage.txt ./...

## Chain:
chain-init: build ## Initialize the blockchain with default settings.
	@echo "${COLOR_CYAN} 🛠️ Initializing chain ${COLOR_RESET}${CHAIN}${COLOR_CYAN} under ${COLOR_YELLOW}${CHAIN_HOME}${COLOR_RESET}"

	@rm -rf "${CHAIN_HOME}"; \
	${CHAIN_BINARY} init axone-node \
	  --chain-id=axone-${CHAIN} \
	  --home "${CHAIN_HOME}"; \
	\
	sed -i $(SED_FLAG) "s/\"stake\"/\"uaxone\"/g" "${CHAIN_HOME}/config/genesis.json"; \
	\
	MNEMONIC_VALIDATOR="island position immense mom cross enemy grab little deputy tray hungry detect state helmet \
	  tomorrow trap expect admit inhale present vault reveal scene atom"; \
	echo $$MNEMONIC_VALIDATOR \
	  | ${CHAIN_BINARY} keys add validator \
	      --recover \
	      --keyring-backend test \
	      --home "${CHAIN_HOME}"; \
	\
	${CHAIN_BINARY} genesis add-genesis-account validator 1000000000uaxone \
	  --keyring-backend test \
	  --home "${CHAIN_HOME}"; \
	\
	NODE_ID=`${CHAIN_BINARY} tendermint show-node-id --home ${CHAIN_HOME}`; \
	${CHAIN_BINARY} genesis gentx validator 1000000uaxone \
	  --node-id $$NODE_ID \
	  --chain-id=axone-${CHAIN} \
	  --keyring-backend test \
      --home "${CHAIN_HOME}"; \
	\
	${CHAIN_BINARY} genesis collect-gentxs \
	  --home "${CHAIN_HOME}"

chain-start: build ## Start the blockchain with existing configuration (see chain-init)
	@echo "${COLOR_CYAN} 🛠️ Starting chain ${COLOR_RESET}${CHAIN}${COLOR_CYAN} with configuration ${COLOR_YELLOW}${CHAIN_HOME}${COLOR_RESET}"; \
	${CHAIN_BINARY} start --moniker ${CHAIN_MONIKER} \
	  --home ${CHAIN_HOME}

chain-stop: ## Stop the blockchain
	@echo "${COLOR_CYAN} ✋️ Stopping chain ${COLOR_RESET}${CHAIN}${COLOR_CYAN} with configuration ${COLOR_YELLOW}${CHAIN_HOME}${COLOR_RESET}"
	@killall axoned

chain-upgrade: build ## Test the chain upgrade from the given FROM_VERSION to the given TO_VERSION.
	@echo "${COLOR_CYAN} ⬆️ Upgrade the chain ${COLOR_RESET}${CHAIN}${COLOR_CYAN} from ${COLOR_YELLOW}${FROM_VERSION}${COLOR_RESET}${COLOR_CYAN} to ${COLOR_YELLOW}${TO_VERSION}${COLOR_RESET}"
	@killall cosmovisor || \
	rm -rf ${TARGET_FOLDER}/${FROM_VERSION}; \
	git clone -b ${FROM_VERSION} https://github.com/axone-protocol/axoned.git ${TARGET_FOLDER}/${FROM_VERSION}; \
	echo "${COLOR_CYAN} 🏗 Build the ${COLOR_YELLOW}${FROM_VERSION}${COLOR_RESET}${COLOR_CYAN} binary...${COLOR_RESET}"; \
	cd ${TARGET_FOLDER}/${FROM_VERSION}; \
	make build; \
	BINARY_OLD=${TARGET_FOLDER}/${FROM_VERSION}/${DIST_FOLDER}/${DAEMON_NAME}; \
	cd ../../; \
	echo $$BINARY_OLD; \
	make chain-init CHAIN_BINARY=$$BINARY_OLD; \
	\
	echo "${COLOR_CYAN} 👩‍🚀 Preparing cosmovisor ${COLOR_RESET}"; \
	export DAEMON_NAME=${DAEMON_NAME}; \
	export DAEMON_HOME=${DAEMON_HOME}; \
	\
    cat <<< $$(jq '.app_state.gov.params.voting_period = "20s"' ${CHAIN_HOME}/config/genesis.json) > ${CHAIN_HOME}/config/genesis.json; \
	\
 	cosmovisor init $$BINARY_OLD; \
 	cosmovisor run start --moniker ${CHAIN_MONIKER} \
 		--home ${CHAIN_HOME} \
 		--log_level debug & \
	sleep 10; \
	echo "${COLOR_CYAN} 🗳️ Submitting software-upgrade tx ${COLOR_RESET}"; \
	$$BINARY_OLD tx upgrade software-upgrade ${TO_VERSION}\
	    --title "Axoned upgrade" \
	    --summary "⬆️ Upgrade the chain from ${FROM_VERSION} to ${TO_VERSION}" \
	    --upgrade-height 20 \
	    --upgrade-info "{}" \
	    --deposit 10000000uaxone \
	    --no-validate \
	    --yes \
	    --from validator \
	    --keyring-backend test \
	    --chain-id axone-${CHAIN} \
	    --home ${CHAIN_HOME}; \
	sleep 5;\
	\
	sleep 5;\
 	$$BINARY_OLD tx gov vote 1 yes \
     		--from validator \
     		--yes \
     		--home ${CHAIN_HOME} \
     		--chain-id axone-${CHAIN} \
     		--keyring-backend test; \
	mkdir -p ${DAEMON_HOME}/cosmovisor/upgrades/${TO_VERSION}/bin && cp ${CHAIN_BINARY} ${DAEMON_HOME}/cosmovisor/upgrades/${TO_VERSION}/bin; \
	wait

## Clean:
.PHONY: clean
clean: ## Remove all the files from the target folder
	@echo "${COLOR_CYAN} 🗑 Cleaning folder $(TARGET_FOLDER)${COLOR_RESET}"
	@rm -rf $(TARGET_FOLDER)/

## Proto:
.PHONY: proto
proto: format-proto lint-proto proto-gen doc-proto ## Generate all resources for proto files (go, doc, etc.)

.PHONY: proto-gen
proto-gen: ## Generate all the code from the Protobuf files
	@echo "${COLOR_CYAN} 📝 Generating code from Protobuf files${COLOR_RESET}"
	@$(DOCKER_PROTO_RUN) sh ./scripts/protocgen-code.sh


## Documentation:
.PHONY: doc
doc: doc-proto doc-command doc-predicate ## Generate all the documentation

.PHONY: doc-proto
doc-proto: proto-gen ## Generate the documentation from the Protobuf files
	@echo "${COLOR_CYAN} 📝 Generating doc from Protobuf files${COLOR_RESET}"
	@$(DOCKER_PROTO_RUN) sh ./scripts/protocgen-doc.sh
	@for MODULE in $(shell find proto -name '*.proto' -maxdepth 3 -print0 | xargs -0 -n1 dirname | sort | uniq | xargs dirname) ; do \
		echo "${COLOR_CYAN} 📖 Generate documentation for $${MODULE} module${COLOR_RESET}" ; \
        DEFAULT_DATASOURCE="./docs/proto/templates/default.yaml" ; \
        MODULE_DATASOURCE="merge:./$${MODULE}/docs.yaml|$${DEFAULT_DATASOURCE}" ; \
        DATASOURCE="docs=`[ -f $${MODULE}/docs.yaml ] && echo $$MODULE_DATASOURCE || echo $${DEFAULT_DATASOURCE}`" ; \
		docker run --rm \
				-v ${HOME}/.cache:/root/.cache \
				-v `pwd`:/usr/src/axoned \
				-w /usr/src/axoned \
				${DOCKER_IMAGE_GOTEMPLATE} \
				-d $$DATASOURCE -f docs/$${MODULE}.md -o docs/$${MODULE}.md ; \
	done
	@docker run --rm \
	  -v `pwd`:/usr/src/axoned \
	  -w /usr/src/axoned/docs \
	  ${DOCKER_IMAGE_MARKDOWNLINT} -f proto

.PHONY: doc-command
doc-command: ## Generate markdown documentation for the command
	@echo "${COLOR_CYAN} 📖 Generate markdown documentation for the command${COLOR_RESET}"
	@OUT_FOLDER="docs/command"; \
	rm -rf $$OUT_FOLDER; \
	go get ./scripts; \
	go run ./scripts/. command; \
	sed -i $(SED_FLAG) 's/(default \"\/.*\/\.axoned\")/(default \"\/home\/john\/\.axoned\")/g' $$OUT_FOLDER/*.md; \
	sed -i $(SED_FLAG) 's/node\ name\ (default\ \".*\")/node\ name\ (default\ \"my-machine\")/g' $$OUT_FOLDER/*.md; \
	sed -i $(SED_FLAG) 's/IP\ (default\ \".*\")/IP\ (default\ \"127.0.0.1\")/g' $$OUT_FOLDER/*.md; \
	sed -i $(SED_FLAG) 's/&lt;appd&gt;/axoned/g' $$OUT_FOLDER/*.md; \
	sed -i $(SED_FLAG) 's/<appd>/axoned/g' $$OUT_FOLDER/*.md; \
	sed -i $(SED_FLAG) -E 's| (https?://[a-zA-Z0-9\.\/_=%-]+)| [\1](\1) |g' $$OUT_FOLDER/*.md; \
	docker run --rm \
	  -v `pwd`:/usr/src/docs \
	  -w /usr/src/docs \
	  ${DOCKER_IMAGE_MARKDOWNLINT} -f $$OUT_FOLDER -c docs/.markdownlint.yaml

.PHONY: doc-predicate
doc-predicate: ## Generate markdown documentation for all the predicates (module logic)
	@echo "${COLOR_CYAN} 📖 Generate markdown documentation for the predicates${COLOR_RESET}"
	@OUT_FOLDER="docs/predicate"; \
	rm -rf $$OUT_FOLDER; \
	mkdir -p $$OUT_FOLDER; \
	go get ./scripts; \
	go run ./scripts/. predicate; \
	docker run --rm \
	  -v `pwd`:/usr/src/docs \
	  -w /usr/src/docs \
	  ${DOCKER_IMAGE_MARKDOWNLINT} -f $$OUT_FOLDER -c docs/.markdownlint.yaml


## Mock:
.PHONY: mock
mock: ## Generate all the mocks (for tests)
	@echo "${COLOR_CYAN} 🧱 Generating all the mocks${COLOR_RESET}"
	@go install go.uber.org/mock/mockgen@v0.5.0
	@mockgen -source=x/mint/types/expected_keepers.go -package testutil -destination x/mint/testutil/expected_keepers_mocks.go
	@mockgen -source=x/vesting/types/expected_keepers.go -package testutil -destination x/vesting/testutil/expected_keepers_mocks.go
	@mockgen -source=x/logic/types/expected_keepers.go -package testutil -destination x/logic/testutil/expected_keepers_mocks.go
	@mockgen -destination x/logic/testutil/gas_mocks.go -package testutil cosmossdk.io/store/types GasMeter
	@mockgen -destination x/logic/testutil/fs_mocks.go -package testutil io/fs FS
	@mockgen -destination x/logic/testutil/read_file_fs_mocks.go -package testutil io/fs ReadFileFS
	@mockgen -source "$$(go list -f '{{.Dir}}' github.com/cosmos/cosmos-sdk/codec/types)/interface_registry.go" \
    	-package testutil \
    	-destination x/logic/testutil/interface_registry_mocks.go

## Release:
.PHONY: release-assets
release-assets: release-binary-all release-checksums ## Generate release assets

release-binary-all: $(RELEASE_TARGETS)

$(RELEASE_TARGETS): ensure-buildx-builder
	@GOOS=$(word 3, $(subst -, ,$@)); \
    GOARCH=$(word 4, $(subst -, ,$@)); \
    BINARY_NAME="axoned-${VERSION}-$$GOOS-$$GOARCH"; \
	echo "${COLOR_CYAN} 🎁 Building ${COLOR_GREEN}$$GOOS $$GOARCH ${COLOR_CYAN}release binary${COLOR_RESET} into ${COLOR_YELLOW}${RELEASE_FOLDER}${COLOR_RESET}"; \
	docker buildx use ${DOCKER_BUILDX_BUILDER}; \
	docker buildx build \
		--platform $$GOOS/$$GOARCH \
		-t $$BINARY_NAME \
		--load \
		.; \
	mkdir -p ${RELEASE_FOLDER}; \
	docker rm -f tmp-axoned || true; \
	docker create -ti --name tmp-axoned $$BINARY_NAME; \
	docker cp tmp-axoned:/usr/bin/axoned ${RELEASE_FOLDER}/$$BINARY_NAME; \
	docker rm -f tmp-axoned; \
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
.PHONY: help
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
