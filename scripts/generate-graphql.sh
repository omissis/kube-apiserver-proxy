#!/bin/sh -x

set -e
set -o errexit -o nounset

KUBERNETES_VERSION="v1.27.1"
TMPDIR="$(mktemp -d)"
FILENAMES=(
    api__v1_openapi
    # api_openapi
    # apis__admissionregistration.k8s.io__v1_openapi
    # apis__admissionregistration.k8s.io__v1alpha1_openapi
    # apis__admissionregistration.k8s.io_openapi
    # apis__apiextensions.k8s.io__v1_openapi
    # apis__apiextensions.k8s.io_openapi
    # apis__apps__v1_openapi
    # apis__apps_openapi
    # apis__authentication.k8s.io__v1_openapi
    # apis__authentication.k8s.io__v1alpha1_openapi
    # apis__authentication.k8s.io__v1beta1_openapi
    # apis__authentication.k8s.io_openapi
    # apis__authorization.k8s.io__v1_openapi
    # apis__authorization.k8s.io_openapi
    # apis__autoscaling__v1_openapi
    # apis__autoscaling__v2_openapi
    # apis__autoscaling_openapi
    # apis__batch__v1_openapi
    # apis__batch_openapi
    # apis__certificates.k8s.io__v1_openapi
    # apis__certificates.k8s.io__v1alpha1_openapi
    # apis__certificates.k8s.io_openapi
    # apis__coordination.k8s.io__v1_openapi
    # apis__coordination.k8s.io_openapi
    # apis__discovery.k8s.io__v1_openapi
    # apis__discovery.k8s.io_openapi
    # apis__events.k8s.io__v1_openapi
    # apis__events.k8s.io_openapi
    # apis__flowcontrol.apiserver.k8s.io__v1beta2_openapi
    # apis__flowcontrol.apiserver.k8s.io__v1beta3_openapi
    # apis__flowcontrol.apiserver.k8s.io_openapi
    # apis__internal.apiserver.k8s.io__v1alpha1_openapi
    # apis__internal.apiserver.k8s.io_openapi
    # apis__networking.k8s.io__v1_openapi
    # apis__networking.k8s.io__v1alpha1_openapi
    # apis__networking.k8s.io_openapi
    # apis__node.k8s.io__v1_openapi
    # apis__node.k8s.io_openapi
    # apis__policy__v1_openapi
    # apis__policy_openapi
    # apis__rbac.authorization.k8s.io__v1_openapi
    # apis__rbac.authorization.k8s.io_openapi
    # apis__resource.k8s.io__v1alpha2_openapi
    # apis__resource.k8s.io_openapi
    # apis__scheduling.k8s.io__v1_openapi
    # apis__scheduling.k8s.io_openapi
    # apis__storage.k8s.io__v1_openapi
    # apis__storage.k8s.io_openapi
    # apis_openapi
    # logs_openapi
    # openid__v1__jwks_openapi
    # version_openapi
)

for FILENAME in ${FILENAMES[@]}; do
    # Download the OpenAPI spec file
    curl "https://raw.githubusercontent.com/kubernetes/kubernetes/${KUBERNETES_VERSION}/api/openapi-spec/v3/${FILENAME}.json" | \

    # Add missing "type" fields
    jq '(select(.components.schemas != null) | .components.schemas[].properties | select(. != null) | .[] | select(.allOf != null) | .type) = "object"' | \
    jq '(select(.components.schemas != null) | .components.schemas[].properties | select(. != null) | .[].items | select(.allOf != null) | .type) = "object"' | \
    jq '(select(.components.schemas != null) | .components.schemas[].properties | select(. != null) | .[].additionalProperties | select(.allOf != null) | .type) = "object"' | \

    # Remove extra response mime types
    jq '(.paths[] | .[] | select(. | type =="object") | .responses | .. | .content? | select(. != null)) |= with_entries(select([.key] | inside(["application/json", "*/*"])))' | \

    # Remove watch endpoints
    jq '(.paths) |= with_entries(select( .key | test("/api/v1/watch") | not))' | \

    # Remove group prefix
    sed -e 's/io\.k8s\.api\.core\.v1\.//g' | \

    # Save the modified OpenAPI spec file
    cat >> "${TMPDIR}/${FILENAME}.json"

    # Convert the OpenAPI spec file to GraphQL schema
    node_modules/.bin/openapi-to-graphql "${TMPDIR}/${FILENAME}.json" --save "internal/graph/${FILENAME}.graphqls"
done
