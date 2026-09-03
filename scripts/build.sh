#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

IMAGE_NAME="${IMAGE_NAME:-ddavideccia/my-vehicles-an-backup-agent}"
IMAGE_TAG="${IMAGE_TAG:-$(git rev-parse --short HEAD 2>/dev/null || echo latest)}"
IMAGE_REF="${IMAGE_NAME}:${IMAGE_TAG}"
DOCKERFILE="${DOCKERFILE:-Dockerfile.prod}"
PUSH="${PUSH:-1}"
PLATFORMS="${PLATFORMS:-linux/amd64,linux/arm64}"
BUILDER="${BUILDER:-multiarch-builder}"

if [ "$PUSH" = "1" ]; then
  # Multi-platform output needs the docker-container driver (the default
  # "docker" driver can only build for the host platform) and can only be
  # assembled in the registry: buildx pushes each platform's manifest as it
  # builds, it can't --load a multi-platform image into the local daemon.
  if ! docker buildx inspect "$BUILDER" >/dev/null 2>&1; then
    docker buildx create --name "$BUILDER" --driver docker-container
  fi

  echo "Building and pushing ${IMAGE_REF} (${PLATFORMS}) from ${DOCKERFILE}"
  docker buildx build \
    --builder "$BUILDER" \
    --pull \
    --platform "${PLATFORMS}" \
    -f "${ROOT_DIR}/${DOCKERFILE}" \
    -t "${IMAGE_REF}" \
    -t "${IMAGE_NAME}:latest" \
    --push \
    "$ROOT_DIR"
  echo "Pushed ${IMAGE_REF} and ${IMAGE_NAME}:latest to Docker Hub for ${PLATFORMS}"
else
  # No push: build a regular single-platform image for the host arch and
  # load it into the local daemon so it's immediately runnable.
  echo "Building ${IMAGE_REF} from ${DOCKERFILE} (local, host platform only)"
  docker build \
    --pull \
    -f "${ROOT_DIR}/${DOCKERFILE}" \
    -t "${IMAGE_REF}" \
    -t "${IMAGE_NAME}:latest" \
    "$ROOT_DIR"
  echo "Built ${IMAGE_REF} (PUSH=0, skipped push and multi-arch build)"
fi
