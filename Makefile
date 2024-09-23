# Local variables
APPS := dashboard.bin pipeline-converter.bin webhook-listener.bin cleaner.bin
IMAGES := dashboard.image pipeline-converter.image webhook-listener.image cleaner.image

# Local Container Registry
REGISTRY_NAME=localhost:50601

# K3s cluster configuration
K3S_CLUSTER_NAME=k3s-default
KUBECONFIG=$(PWD)/kubeconfig

# TODO: add deploy targets

help: ## Display this help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Common targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
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

retrieve-kubeconfig: ## Retrieve kubeconfig for local development environment
	@echo "Retrieving kubeconfig for local development environment"
	@k3d kubeconfig get ${K3S_CLUSTER_NAME} > ${KUBECONFIG}
	@echo "export KUBECONFIG=${KUBECONFIG}" > set-kubeconfig.sh
	@echo "Kubeconfig retrieved"
	@echo "Run 'source set-kubeconfig.sh' to set the kubeconfig environment variable for the current shell"

all: $(APPS) $(IMAGES) ## Build all go applications and docker images

bin: $(APPS) ## Build all go applications

images: $(IMAGES) ## Build all docker images

deploy: ## Deploy applications to devel k8s cluster
	@echo "TODO: Deploying applications to devel k8s cluster"

# Build go application
$(APPS):
	@echo "Building $(basename $@)"
	go build -o bin/$(basename $@) cmd/$(basename $@)/main.go

# Build docker image
$(IMAGES):
	docker build \
		-f build/Dockerfile \
		--build-arg APP_NAME=$(basename $@) \
		-t ${REGISTRY_NAME}/$(basename $@):latest .
	docker push ${REGISTRY_NAME}/$(basename $@):latest

clean: ## Clean up
	@echo "Cleaning up"
	rm -rf ${KUBECONFIG}
	rm -rf set-kubeconfig.sh
	rm -rf bin/*
	for image in $(IMAGES); do \
		docker rmi -f ${REGISTRY_NAME}/$${image%.image}:latest || true; \
	done
