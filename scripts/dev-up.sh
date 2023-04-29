#!/usr/bin/env bash

# Variables
PROJECT_NAME="kube-apiserver-proxy"
PROJECT_DOMAIN="kube-apiserver-proxy.dev"
WORKING_DIR=$(pwd)
CLUSTER_VERSION="${1}"
FORCE="${2}"

# Utilities
# This function will create the ctlptl registry (which is local registry used by Tilt) and the kind cluster.
function create {
  CLUSTER_VERSION=$1
  CONFIG_DIR=$2

  echo "Creating cluster and local registry..."
  ctlptl apply -f "${CONFIG_DIR}/ctlptl-registry.yml"
  KIND_IMAGE="kindest/node:v${CLUSTER_VERSION}" yq e '.kindV1Alpha4Cluster.nodes[0].image |= strenv(KIND_IMAGE)' < "${CONFIG_DIR}/kind-config.yml" | ctlptl apply -f - || true

  docker update --restart=no "${PROJECT_NAME}-kind-control-plane" "${PROJECT_NAME}-kind-registry"
}

# This function will create and inject self-signed TLS certificates into the cluster, using mkcert(it install the CA certificate to your browser's trusted certificates). This is necessary for the ingress to work.
# Remember to add "0.0.0.0 ${PROJECT_DOMAIN}" line to your /etc/hosts file
function setup_certs {
  MANIFESTS_DIR=$1
  CERTS_DIR=$2
  FORCE=$3

  mkdir -p "${CERTS_DIR}"

  if [[ "${FORCE}" -eq "1" ]]; then
    rm "${MANIFESTS_DIR}/${PROJECT_DOMAIN}-tls.yaml"
    rm "${MANIFESTS_DIR}/wildcard.${PROJECT_DOMAIN}-tls.yaml"
  fi

  # setup self-signed tls certificates
  if [[ ! -f "${MANIFESTS_DIR}/${PROJECT_DOMAIN}-tls.yaml" ]] ||
    [[ ! -f "${MANIFESTS_DIR}/wildcard.${PROJECT_DOMAIN}-tls.yaml" ]]; then
    echo "Creating TLS certificates.."

    (cd "${CERTS_DIR}" && mkcert "*.${PROJECT_DOMAIN}" && mkcert "${PROJECT_DOMAIN}" && mkcert -install) &&
      kubectl create secret tls wildcard.${PROJECT_DOMAIN}-tls \
        --cert="${CERTS_DIR}/_wildcard.${PROJECT_DOMAIN}.pem" \
        --key="${CERTS_DIR}/_wildcard.${PROJECT_DOMAIN}-key.pem" \
        -o yaml --dry-run=client > "${MANIFESTS_DIR}/wildcard.${PROJECT_DOMAIN}-tls.yaml"

    kubectl create secret tls ${PROJECT_DOMAIN}-tls \
      --cert="${CERTS_DIR}/${PROJECT_DOMAIN}.pem" \
      --key="${CERTS_DIR}/${PROJECT_DOMAIN}-key.pem" \
      -o yaml --dry-run=client > "${MANIFESTS_DIR}/${PROJECT_DOMAIN}-tls.yaml"
  fi
}

# This function will start the ctlptl registry (which is local registry used by Tilt) and the kind cluster.
function start {
  echo "Starting cluster and local registry..."
  docker start "${PROJECT_NAME}-kind-registry" "${PROJECT_NAME}-kind-control-plane"
}

# Setup

# Setup self-signed TLS certificates for the ingress to work (it install the CA certificate to your browser's trusted certificates).
# Remember to add 0.0.0.0 ${PROJECT_DOMAIN} line to your /etc/hosts file.
setup_certs "${WORKING_DIR}/configs/kubernetes-manifests" "${WORKING_DIR}/configs/certs" "${FORCE}"

# Exec

if [[ -n "$(docker ps -q -a -f name=${PROJECT_NAME}-kind-registry || true)" ]] &&
  [[ -n "$(docker ps -q -a -f name=${PROJECT_NAME}-kind-control-plane || true)" ]]; then
  start
else
  create "${CLUSTER_VERSION}" "${WORKING_DIR}/configs/kubernetes-manifests"
fi

echo "Starting Tilt's dev environment..."
kind get kubeconfig --name "${PROJECT_NAME}-kind" > "${WORKING_DIR}/configs/.kubeconfig"
kubectl --kubeconfig "${WORKING_DIR}/configs/.kubeconfig" config set-context "${PROJECT_NAME}-kind"
export KUBECONFIG="${WORKING_DIR}/configs/.kubeconfig" && tilt up -f "${WORKING_DIR}/Tiltfile"

echo "Done."
