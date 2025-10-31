# Build stage - use BUILDPLATFORM for cross-compilation
FROM --platform=$BUILDPLATFORM registry.access.redhat.com/ubi9/go-toolset:1.22 AS builder

# Build arguments for cross-compilation
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Switch to root to set up workspace
USER root

# Set working directory
WORKDIR /workspace

# Copy go mod files
COPY --chown=1001:0 go.mod go.sum ./

# Download dependencies
USER 1001
RUN go mod download

# Copy source code
USER root
COPY --chown=1001:0 main.go ./
COPY --chown=1001:0 internal/ internal/

# Build the plugin for the target architecture
USER 1001
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -v -o velero-pvc-remove-selector .

# Final stage
FROM registry.access.redhat.com/ubi9-minimal:latest

# Create plugins directory
RUN mkdir /plugins

# Copy the binary from builder
COPY --from=builder /workspace/velero-pvc-remove-selector /plugins/

# Use non-root user
USER 65534:65534

# Velero plugin entrypoint - copies plugin to target directory
ENTRYPOINT ["/bin/bash", "-c", "cp /plugins/* /target/."]
