.PHONY: build build-amd64 build-arm64 container container-amd64 container-arm64 push clean test manifest-create manifest-push

# CONTAINER_TOOL defines the container tool to be used for building images.
# By default, this Makefile uses docker, as the target commands have been tested primarily with it.
# However, if docker is not available, the Makefile will attempt to use podman if it's installed.
# You may also set CONTAINER_TOOL directly as an environment variable to specify a different tool.
# If neither docker nor podman is found, or if the specified tool is unavailable, the Makefile will exit with an error.

# Set CONTAINER_TOOL to Docker or Podman if not already defined by the user
CONTAINER_TOOL ?= $(shell \
  if command -v docker >/dev/null 2>&1; then echo docker; \
  elif command -v podman >/dev/null 2>&1; then echo podman; \
  else echo ""; \
  fi \
)
ifeq ($(shell command -v $(CONTAINER_TOOL) >/dev/null 2>&1 && echo found),)
  $(error The selected container tool '$(CONTAINER_TOOL)' is not available on this system. Please install it or choose a different tool.)
endif
$(info Using Container Tool: $(CONTAINER_TOOL))

# Image URL to use all building/pushing image targets
IMAGE_REGISTRY ?= quay.io
IMAGE_REPO ?= mpryc/oadp-pvc-remove-selector
IMAGE_TAG ?= latest
IMAGE ?= $(IMAGE_REGISTRY)/$(IMAGE_REPO):$(IMAGE_TAG)

# Platform-specific image tags
IMAGE_AMD64 ?= $(IMAGE)-amd64
IMAGE_ARM64 ?= $(IMAGE)-arm64

# Binary name
BINARY_NAME = velero-pvc-remove-selector

# Supported platforms
PLATFORMS ?= linux/amd64,linux/arm64

# Build the binary for amd64
build-amd64:
	@echo "Building $(BINARY_NAME) for amd64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/$(BINARY_NAME)-amd64 .

# Build the binary for arm64
build-arm64:
	@echo "Building $(BINARY_NAME) for arm64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o bin/$(BINARY_NAME)-arm64 .

# Build the binary (default to amd64)
build: build-amd64

# Build container image for amd64
container-amd64:
	@echo "Building container image $(IMAGE_AMD64)..."
	$(CONTAINER_TOOL) build --platform linux/amd64 -t $(IMAGE_AMD64) -f Containerfile .

# Build container image for arm64
container-arm64:
	@echo "Building container image $(IMAGE_ARM64)..."
	$(CONTAINER_TOOL) build --platform linux/arm64 -t $(IMAGE_ARM64) -f Containerfile .

# Build both architecture images
container: container-amd64 container-arm64
	@echo "Built multi-arch container images"

# Push both architecture images
push: container
	@echo "Pushing container images..."
	$(CONTAINER_TOOL) push $(IMAGE_AMD64)
	$(CONTAINER_TOOL) push $(IMAGE_ARM64)
	@echo "Creating and pushing manifest..."
	$(MAKE) manifest-create
	$(MAKE) manifest-push

# Create manifest list
manifest-create:
	@echo "Creating manifest list $(IMAGE)..."
	$(CONTAINER_TOOL) manifest create $(IMAGE) \
		$(IMAGE_AMD64) \
		$(IMAGE_ARM64)

# Push manifest list
manifest-push:
	@echo "Pushing manifest list $(IMAGE)..."
	$(CONTAINER_TOOL) manifest push $(IMAGE)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Run go fmt
fmt:
	@echo "Running go fmt..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run all checks
check: fmt vet test

# Help target
help:
	@echo "Available targets:"
	@echo "  build           - Build the binary for amd64 (default)"
	@echo "  build-amd64     - Build the binary for amd64"
	@echo "  build-arm64     - Build the binary for arm64"
	@echo "  container       - Build multi-arch container images (amd64 and arm64)"
	@echo "  container-amd64 - Build container image for amd64"
	@echo "  container-arm64 - Build container image for arm64"
	@echo "  push            - Build, push images, and create/push manifest"
	@echo "  manifest-create - Create manifest list for multi-arch"
	@echo "  manifest-push   - Push manifest list"
	@echo "  test            - Run tests"
	@echo "  deps            - Download and tidy dependencies"
	@echo "  clean           - Clean build artifacts"
	@echo "  fmt             - Run go fmt"
	@echo "  vet             - Run go vet"
	@echo "  check           - Run fmt, vet, and test"
	@echo "  help            - Show this help message"
	@echo ""
	@echo "Container tool: $(CONTAINER_TOOL)"
	@echo "Image: $(IMAGE)"
	@echo "Platforms: $(PLATFORMS)"
