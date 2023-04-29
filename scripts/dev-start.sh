#!/usr/bin/env bash

PROJECT_NAME="kube-apiserver-proxy"

echo "Starting cluster and local registry..."

docker start "${PROJECT_NAME}-kind-control-plane" "${PROJECT_NAME}-kind-registry"

echo "Done."
