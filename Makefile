
# Image URL to use all building/pushing image targets
IMG ?= carolynvs/cloudkinds
TAG ?= latest

all: test manager

# Run tests
test: generate fmt vet manifests
	go test ./pkg/... ./cmd/... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager github.com/carolynvs/cloudkinds/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/manager/main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crds

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: docker-push
	kubectl apply -f config/crds
	kustomize build config/default | kubectl apply -f -
	helm upgrade --install cloudkinds charts/cloudkinds \
	 --recreate-pods --set sampleProvider.include=true,imagePullPolicy="Always",deploymentStrategy="Recreate"
	helm upgrade --install cloudkinds-svcat charts/cloudkinds-servicecatalog \
	  --recreate-pods --set imagePullPolicy="Always",deploymentStrategy="Recreate"

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Generate code
generate:
	go generate ./pkg/... ./cmd/...

# Build the docker image
docker-build:
	docker build -t ${IMG}:${TAG} -f cmd/manager/Dockerfile .
	docker build -t ${IMG}-sampleprovider:${TAG} -f cmd/sampleprovider/Dockerfile .
	docker build -t ${IMG}-servicecatalog:${TAG} -f cmd/servicecatalog/Dockerfile .
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push: docker-build
	docker push ${IMG}:${TAG}
	docker push ${IMG}-sampleprovider:${TAG}
	docker push ${IMG}-servicecatalog:${TAG}
