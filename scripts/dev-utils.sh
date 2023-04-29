#!/usr/bin/env bash

# Variables
PROJECT_NAME="kube-apiserver-proxy"
PROJECT_DOMAIN="${PROJECT_DOMAIN}"

# This function will create the ctlptl registry (which is local registry used by Tilt) and the kind cluster.
function create {
  CLUSTER_VERSION=$1
  CONFIG_DIR=$2

  echo "Creating cluster and local registry..."
  ctlptl apply -f "${CONFIG_DIR}/ctlptl-registry.yml"
  KIND_IMAGE="kindest/node:v${CLUSTER_VERSION}" yq e '.kindV1Alpha4Cluster.nodes[0].image |= strenv(KIND_IMAGE)' < "${CONFIG_DIR}/kind-config.yml" | ctlptl apply -f - || true
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
  docker start ${PROJECT_NAME}-kind-registry ${PROJECT_NAME}-kind-control-plane
}
