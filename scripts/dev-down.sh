#!/usr/bin/env bash

PROJECT_NAME="kube-apiserver-proxy"
FORCE="${1}"

echo "Deleting cluster..."
docker rm -f "${PROJECT_NAME}-kind-control-plane"

if [[ "${FORCE}" -eq "1" ]]; then
  echo "Deleting local registry..."
  docker rm -f "${PROJECT_NAME}-kind-registry"
else
  echo "Stopping local registry..."
  docker stop "${PROJECT_NAME}-kind-registry"
fi

echo "Done."
