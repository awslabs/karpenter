
# Image URL to use all building/pushing image targets
IMG ?= ${KO_DOCKER_REPO}/karpenter:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=false"
GOLINT_OPTIONS ?= "--set_exit_status=1"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: build test

# Run tests
test: generate fmt vet manifests
	go test ./... -cover

# Build controller binary
build: generate fmt vet tidy
	go build -o bin/karpenter karpenter/main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run karpenter/main.go --enable-leader-election=false --enable-webhook=false

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kustomize build config/dev | ko apply -B -f -

undeploy: manifests
	kustomize build config/dev | ko delete -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	controller-gen $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	golint $(GOLINT_OPTIONS) ./...
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Tidy up modules
tidy:
	go mod tidy

# Generate code
generate:
	controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./pkg/apis/..."
	./hack/boilerplate.sh

# Build the docker image
docker-build:
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

docker-release: docker-build docker-push
