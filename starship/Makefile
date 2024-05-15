NAME = axone-starship
FILE = configs/local.yaml

HELM_REPO = starship
HELM_CHART = devnet
HELM_VERSION = v0.1.38

###############################################################################
###                              Run tests                                  ###
###############################################################################

.PHONY: test
test:
	cd tests/ && go test -v ./...

###############################################################################
###                              All commands                               ###
###############################################################################

.PHONY: setup
setup: setup-deps setup-helm

.PHONY: stop
stop: stop-forward delete

.PHONY: clean
clean: stop clean-kind

###############################################################################
###                          Dependency check                               ###
###############################################################################

.PHONY: check
setup-deps:
	bash $(CURDIR)/scripts/dev-setup.sh


###############################################################################
###                          Docker build cmds                              ###
###############################################################################

DOCKER_BUILDX_NAME := axone

.PHONY: docker-setup
docker-setup:
	-@docker buildx create --use --name $(DOCKER_BUILDX_NAME)

.PHONY: docker-build
docker-build: docker-setup
	cd $(CURDIR)/.. && docker buildx build --platform linux/arm64,linux/amd64 -t anmol1696/axoned:latest . --push -f starship/Dockerfile

###############################################################################
###                              Helm Charts                                ###
###############################################################################

setup-helm:
	helm repo add $(HELM_REPO) https://cosmology-tech.github.io/starship/
	helm repo update
	helm search repo $(HELM_REPO)/$(HELM_CHART) --version $(HELM_VERSION)

install:
	helm install -f $(FILE) $(NAME) $(HELM_REPO)/$(HELM_CHART) --version $(HELM_VERSION)
	@echo "Please run \`$ kubectl get pods\` and wait till all pods are in running state"
	@echo "Then run \`$ make port-forward\` to forward all the ports locally"

install-dev:
	$(MAKE) install FILE=configs/devnet.yaml

install-local:
	$(MAKE) install FILE=configs/local.yaml

install-ci:
	$(MAKE) install FILE=configs/ci.yaml

delete:
	-helm delete $(NAME)

###############################################################################
###                                 Port forward                            ###
###############################################################################

.PHONY: port-forward
port-forward:
	bash $(CURDIR)/scripts/port-forward.sh --config=$(FILE)

.PHONY: port-forward-dev
port-forward-dev:
	$(MAKE) port-forward FILE=configs/devnet.yaml

.PHONY: port-forward-local
port-forward-local:
	$(MAKE) port-forward FILE=configs/local.yaml

.PHONY: port-forward-ci
port-forward-ci:
	$(MAKE) port-forward FILE=configs/ci.yaml


.PHONY: stop-forward
stop-forward:
	-pkill -f "port-forward"

###############################################################################
###                          Local Kind Setup                               ###
###############################################################################
KIND_CLUSTER=starship

.PHONY: setup-kind
setup-kind:
	kind create cluster --name $(KIND_CLUSTER)

.PHONY: clean-kind
clean-kind:
	kind delete cluster --name $(KIND_CLUSTER)
