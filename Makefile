
# Image URL to use all building/pushing image targets
REGISTRY ?= carolynvs/
IMG ?= ${REGISTRY}cloudkinds
TAG ?= latest

all: test build

# Run tests
test: build fmt vet
	go test ./pkg/... ./cmd/... -coverprofile cover.out

# Build cloudkinds binary
build: generate
	go build -o bin/cloudkinds github.com/carolynvs/cloudkinds/cmd/cloudkinds

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate
	go run ./cmd/cloudkinds/main.go

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy:
	helm upgrade --install cloudkinds --namespace cloudkinds charts/cloudkinds \
	 --recreate-pods --set sampleProvider.include=false \
	 --set image.registry="${IMG}",image.tag="${TAG}" \
	 --set imagePullPolicy="Always",deploymentStrategy="Recreate"

# Generate kubebuilder manifests e.g. CRD, RBAC etc
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all
	@echo "Copy the changes from config/* to charts/*"

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
	docker build -t ${IMG}:${TAG} -f cmd/cloudkinds/Dockerfile .
	docker build -t ${IMG}-sampleprovider:${TAG} -f cmd/sampleprovider/Dockerfile .

# Push the docker image
docker-push: docker-build
	docker push ${IMG}:${TAG}
	docker push ${IMG}-sampleprovider:${TAG}
