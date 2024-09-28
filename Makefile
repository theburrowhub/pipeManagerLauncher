.PHONY: deploy release clean shell

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

# TODO: add commitlint

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

setup-cluster: ## Set up local development environment
	@echo "Setting up local development environment..."
	k3d registry create registry -p ${K3D_REGISTRY_PORT}
	k3d cluster create --registry-use ${K3D_REGISTRY_NAME}:${K3D_REGISTRY_PORT} -a 3
	@echo "Local development environment set up"

get-kubeconfig: ## Retrieve kubeconfig for local development environment
	@echo "Retrieving kubeconfig for local development environment"
	@k3d kubeconfig get ${K3S_CLUSTER_NAME} > ${KUBECONFIG}
	@chmod 600 ${KUBECONFIG}
	@echo "export KUBECONFIG=${KUBECONFIG}" > set-kubeconfig.sh
	@echo "Kubeconfig retrieved"
	@echo "Run 'source set-kubeconfig.sh' to set the kubeconfig environment variable for the current shell"

remove-cluster: ## Remove local development environment
	@echo "Removing local development environment..."
	k3d cluster delete ${K3S_CLUSTER_NAME}
	k3d registry delete ${K3D_REGISTRY_NAME}
	@echo "Local development environment removed"

shell: ## Open a shell in the devbox
	devbox shell

all: $(APPS) $(IMAGES) ## Build all go applications and docker images

bin: $(APPS) ## Build all go applications

images: $(IMAGES) ## Build all docker images

deploy: ## Deploy applications to devel k8s cluster
	@echo "Deploying applications to devel k8s cluster..."
	helm upgrade --install --wait --timeout 300s \
		--kubeconfig ${KUBECONFIG} --create-namespace --namespace pipe-manager \
		-f configs/devel/values.yaml \
		-f configs/devel/config.yaml \
		webhook-listener ./deploy/charts/webhook-listener

uninstall: ## Delete applications from devel k8s cluster
	@echo "Deleting applications from devel k8s cluster..."
	helm delete --kubeconfig ${KUBECONFIG} --namespace pipe-manager webhook-listener

port-forward: ## Port forward to devel k8s cluster
	@echo "Port forwarding to devel k8s cluster..."
	kubectl --kubeconfig ${KUBECONFIG} port-forward svc/webhook-listener 8080:80 --namespace pipe-manager

tunnel: ## Tunnel with ngrok to devel k8s cluster
	@echo "Tunneling with ngrok to devel k8s cluster..."
	ngrok http 8080

release: ## Release applications to prod k8s cluster
	@echo "TODO: Release applications and helm charts"
	@echo goreleaser build --snapshot
	@echo goreleaser release --snapshot

helm-docs: ## Generate helm documentation
	@echo "Generating helm documentation..."
	helm-docs

create-git-secret: ## Create git secret in devel k8s cluster using local ssh key
	@echo "Creating git secret..."
	kubectl --kubeconfig ${KUBECONFIG} create secret generic git-secret \
		--namespace pipe-manager \
		--from-file=id_rsa=${HOME}/.ssh/id_rsa

# Build go application
$(APPS):
	@echo "Building $(basename $@)"
	go build \
		-ldflags "-X ${PROJECT_NAME}/internal/pkg/version.Version=$(shell cz version -p)" \
		-o bin/$(basename $@) cmd/$(basename $@)/main.go

# Build docker image
$(IMAGES):
	docker build \
		-f Dockerfile \
		--build-arg APP_NAME=$(basename $@) \
		--build-arg APP_VERSION=$(shell cz version -p) \
		-t ${K3D_REGISTRY_NAME}:${K3D_REGISTRY_PORT}/$(basename $@):$(shell cz version -p) .
	docker push ${K3D_REGISTRY_NAME}:${K3D_REGISTRY_PORT}/$(basename $@):$(shell cz version -p)

clean: ## Clean up
	@echo "Cleaning up"
	rm -rf ${KUBECONFIG}
	rm -rf set-kubeconfig.sh
	rm -rf bin/*
	rm -rf dist
	rm -rf vendor
	for image in $(IMAGES); do \
		docker rmi -f ${REGISTRY_NAME}/$${image%.image}:latest || true; \
	done
