# Velero PVC Remove Selector Plugin

> **⚠️ WARNING: Example Implementation Only**
>
> This is an example implementation of a Velero plugin that has **never been tested in production or testing deployments**.
> It should be treated as a **reference example** and not as a production-ready plugin to be used.
> Use at your own risk and ensure thorough testing in your specific environment before any production use.

A Velero RestoreItemAction (RIA) plugin that removes the `selector` field from PersistentVolumeClaim objects during restore operations.

## Overview

This plugin is designed to work with Velero/OADP and automatically removes the `spec.selector` field from PVCs during restore. This is useful when restoring PVCs that have selectors that might conflict with the target environment.

## Project Structure

```
.
├── internal/
│   └── plugin/
│       └── pvc_restore.go             # RIA implementation
├── main.go                            # Plugin entry point
├── Containerfile                      # Multi-arch container build
├── Makefile                           # Build automation
├── go.mod                             # Go module definition
├── go.sum                             # Go dependencies checksum
└── README.md                          # This file
```

## Building

### Prerequisites

- Go 1.22 or later
- Podman or Docker (automatically detected by Makefile)

The Makefile automatically detects whether Docker or Podman is available on your system. You can override this by setting the `CONTAINER_TOOL` environment variable:

```bash
export CONTAINER_TOOL=podman  # or docker
```

### Build Binary

Build for your current architecture (amd64 by default):
```bash
make build
```

Build for specific architectures:
```bash
make build-amd64  # Build for amd64
make build-arm64  # Build for arm64
```

### Build Container Images

The default `container` target builds multi-arch images for both amd64 and arm64 using cross-compilation:

```bash
make container
```

The Containerfile uses `--platform=$BUILDPLATFORM` to enable efficient cross-compilation. This allows building for different target architectures (amd64, arm64) from any build platform.

This will create two platform-specific images:
- `${IMAGE}-amd64` for linux/amd64
- `${IMAGE}-arm64` for linux/arm64

To build for a specific architecture:
```bash
make container-amd64  # Build only amd64 image
make container-arm64  # Build only arm64 image
```

To build with custom image settings:

```bash
make container IMAGE_REGISTRY=quay.io IMAGE_REPO=yourname/oadp-pvc-remove-selector IMAGE_TAG=v1.0.0
```

### Push Container Images and Manifest

The `push` target builds both architecture images, pushes them, and creates a multi-arch manifest:

```bash
make push
```

This will:
1. Build images for amd64 and arm64
2. Push both images to the registry
3. Create a manifest list that references both images
4. Push the manifest list

Users can then pull the appropriate image for their architecture automatically:
```bash
podman pull quay.io/yourname/oadp-pvc-remove-selector:v1.0.0
```

## Installation

### 1. Build and Push the Plugin Image

```bash
make push IMAGE_REGISTRY=quay.io IMAGE_REPO=yourname/oadp-pvc-remove-selector IMAGE_TAG=v1.0.0
```

### 2. Install the Plugin

#### Using OADP DataProtectionApplication

Add the plugin to your DataProtectionApplication custom resource:

```yaml
apiVersion: oadp.openshift.io/v1alpha1
kind: DataProtectionApplication
metadata:
  name: dpa-custom-plugin-sample
spec:
  configuration:
    velero:
      defaultPlugins:
      # Your current plugins stay...
      customPlugins:
      - name: pvc-remove-selector
        image: quay.io/migi/oadp-pvc-remove-selector:latest
```

## Usage

Once installed, the plugin will automatically process all PersistentVolumeClaim resources during restore operations. No additional configuration is required.

### How It Works

1. During a Velero restore operation, when a PVC is being restored, the plugin intercepts it
2. The plugin checks if the PVC has a `spec.selector` field
3. If present, the selector is removed
4. The modified PVC is then restored to the cluster

## Development

### Running Tests

```bash
make test
```

### Code Formatting and Vetting

```bash
make fmt
make vet
```

### Run All Checks

```bash
make check
```

### Clean Build Artifacts

```bash
make clean
```

## References

- [Velero Plugin Documentation](https://velero.io/docs/main/custom-plugins/)
- [Velero RestoreItemAction Interface](https://velero.io/docs/main/restore-item-action/)
- [Velero Example Plugins](https://github.com/vmware-tanzu/velero-plugin-example)
