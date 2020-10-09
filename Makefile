# Image URL to use all building/pushing image targets
IMG ?= ${KO_DOCKER_REPO}/karpenter:latest

all: generate verify build test
ci: generate verify build battletest

# Build controller binary
build:
	go build -o bin/karpenter karpenter/main.go

# Run tests
test:
	ginkgo -r

# Run stronger tests
battletest:
	ginkgo -r \
		-cover -coverprofile=coverage.out -outputdir=. -coverpkg=./pkg/... \
		--randomizeAllSpecs --randomizeSuites -race
	sed '$d' coverage.out
	go tool cover -func coverage.out

# Verify code. Includes dependencies, linting, formatting, etc
verify:
	go mod tidy
	go mod download
	go vet ./...
	go fmt ./...
	golangci-lint run

# Generate code. Must be run if changes are made to ./pkg/apis/...
generate:
	controller-gen \
		object:headerFile="hack/boilerplate.go.txt" \
		webhook \
		crd:trivialVersions=false \
		rbac:roleName=karpenter \
		paths="./pkg/..." \
		output:crd:artifacts:config=config/crd/bases \
		output:webhook:artifacts:config=config/webhook

	./hack/boilerplate.sh

	# Hack to remove v1.AdmissionReview until https://github.com/kubernetes-sigs/controller-runtime/issues/1161 is fixed
	perl -pi -e 's/^  - v1$$//g' config/webhook/manifests.yaml
	# CRDs don't currently jive with volatile time.
	# `properties[lastTransitionTime].type: Unsupported value: "Any": supported
	# values: "array", "boolean", "integer", "number", "object", "string"`
	perl -pi -e 's/Any/string/g' config/crd/bases/autoscaling.karpenter.sh_horizontalautoscalers.yaml
	perl -pi -e 's/Any/string/g' config/crd/bases/autoscaling.karpenter.sh_scalablenodegroups.yaml
	perl -pi -e 's/Any/string/g' config/crd/bases/autoscaling.karpenter.sh_metricsproducers.yaml


# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
apply:
	kubectl kustomize config/dev | ko apply -B -f -

delete:
	kubectl kustomize config/dev | ko delete -f -

# Build and release a container image
release:
	docker build . -t ${IMG}
	docker push ${IMG}

# Install developer toolchain
toolchain:
	./hack/toolchain.sh

.PHONY: all test build release run apply delete verify generate toolchain
