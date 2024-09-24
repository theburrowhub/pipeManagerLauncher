# Local variables
APPS := dashboard.bin pipeline-converter.bin webhook-listener.bin cleaner.bin
IMAGES := dashboard.image pipeline-converter.image webhook-listener.image cleaner.image

# Local Development Environment
K3D_REGISTRY_NAME=k3d-registry
K3D_REGISTRY_PORT=5111
K3S_CLUSTER_NAME=k3s-default
KUBECONFIG=$(PWD)/kubeconfig

# Project name
PROJECT_NAME=github.com/sergiotejon/pipeManager

# TODO: add deploy targets

help: ## Display this help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Common targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "*** Run \033[36mmake setup\033[0m to set up local development environment. ***"
	@echo ""
	@echo "To build go applications:"
	@for app in $(APPS); do \
		echo "  \033[36mmake $$app\033[0m"; \
	done
	@echo ""
	@echo "and docker images:"
	@for image in $(IMAGES); do \
		echo "  \033[36mmake $$image\033[0m"; \
	done

setup: ## Set up local development environment
	@echo "Setting up local development environment..."
	k3d registry create registry -p ${K3D_REGISTRY_PORT}
	k3d cluster create --registry-use ${K3D_REGISTRY_NAME}:${K3D_REGISTRY_PORT} -a 3
	@echo "Local development environment set up"

retrieve-kubeconfig: ## Retrieve kubeconfig for local development environment
	@echo "Retrieving kubeconfig for local development environment"
	@k3d kubeconfig get ${K3S_CLUSTER_NAME} > ${KUBECONFIG}
	@echo "export KUBECONFIG=${KUBECONFIG}" > set-kubeconfig.sh
	@echo "Kubeconfig retrieved"
	@echo "Run 'source set-kubeconfig.sh' to set the kubeconfig environment variable for the current shell"

remove: ## Remove local development environment
	@echo "Removing local development environment..."
	k3d cluster delete ${K3S_CLUSTER_NAME}
	k3d registry delete ${K3D_REGISTRY_NAME}
	@echo "Local development environment removed"

all: $(APPS) $(IMAGES) ## Build all go applications and docker images

bin: $(APPS) ## Build all go applications

images: $(IMAGES) ## Build all docker images

deploy: ## Deploy applications to devel k8s cluster
	@echo "Deploying applications to devel k8s cluster..."
	helm upgrade --install --wait --timeout 300s \
		--create-namespace --namespace pipe-manager \
		-f env/devel/values.yaml \
		webhook-listener ./charts/webhook-listener

# Build go application
$(APPS):
	@echo "Building $(basename $@)"
	go build \
		-ldflags "-X ${PROJECT_NAME}/internal/pkg/version.Version=$(shell cz version -p)" \
		-o bin/$(basename $@) cmd/$(basename $@)/main.go

# Build docker image
$(IMAGES):
	docker build \
		-f build/Dockerfile \
		--build-arg APP_NAME=$(basename $@) \
		--build-arg APP_VERSION=$(shell cz version -p) \
		-t ${K3D_REGISTRY_NAME}:${K3D_REGISTRY_PORT}/$(basename $@):$(shell cz version -p) .
	docker push ${K3D_REGISTRY_NAME}:${K3D_REGISTRY_PORT}/$(basename $@):$(shell cz version -p)

clean: ## Clean up
	@echo "Cleaning up"
	rm -rf ${KUBECONFIG}
	rm -rf set-kubeconfig.sh
	rm -rf bin/*
	for image in $(IMAGES); do \
		docker rmi -f ${REGISTRY_NAME}/$${image%.image}:latest || true; \
	done
