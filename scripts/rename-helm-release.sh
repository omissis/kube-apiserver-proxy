#!/bin/sh

PROJECT_NAME="kube-apiserver-proxy"
HELM_TAG=$(yq e '.version' ./helm_chart/Chart.yaml)
CR_TOKEN=$1

RELEASE_INFO=$(
    curl -sSL \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer ${CR_TOKEN}" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    "https://api.github.com/repos/omissis/${PROJECT_NAME}/releases/tags/helm-chart-v${HELM_TAG}"
)

RELEASE_ID=$(echo "${RELEASE_INFO}" | jq -r '.id')

curl -sSL \
-X PATCH \
-H "Accept: application/vnd.github+json" \
-H "Authorization: Bearer ${CR_TOKEN}" \
-H "X-GitHub-Api-Version: 2022-11-28" \
"https://api.github.com/repos/omissis/${PROJECT_NAME}/releases/${RELEASE_ID}" \
-d "{\"name\":\"Helm Chart v${HELM_TAG}\"}"
