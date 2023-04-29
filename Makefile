ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
SHELL := /bin/bash
PROJECT_NAME := kube-apiserver-proxy

.DEFAULT_GOAL := help

# ----------------------------------------------------------------------------------------------------------------------
# Private variables
# ----------------------------------------------------------------------------------------------------------------------

_DOCKER_FILELINT_IMAGE=cytopia/file-lint:latest-0.8
_DOCKER_HADOLINT_IMAGE=hadolint/hadolint:v2.12.0
_DOCKER_JSONLINT_IMAGE=cytopia/jsonlint:1.6
_DOCKER_MAKEFILELINT_IMAGE=cytopia/checkmake:latest-0.5
_DOCKER_MARKDOWNLINT_IMAGE=davidanson/markdownlint-cli2:v0.6.0
_DOCKER_SHELLCHECK_IMAGE=koalaman/shellcheck-alpine:v0.9.0
_DOCKER_SHFMT_IMAGE=mvdan/shfmt:v3-alpine
_DOCKER_YAMLLINT_IMAGE=cytopia/yamllint:1

# TODO: replace this image with a kube-apiserver-proxy-specific one
_DOCKER_TOOLS_IMAGE=omissis/go-jsonschema:tools-latest

_PROJECT_DIRECTORY=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

# ----------------------------------------------------------------------------------------------------------------------
# Utility functions
# ----------------------------------------------------------------------------------------------------------------------

#1: docker image
#2: script name
define run-script-docker
	@docker run --rm \
		-v ${_PROJECT_DIRECTORY}:/data \
		-w /data \
		--entrypoint /bin/sh \
		$(1) scripts/$(2).sh
endef

# check-variable-%: Check if the variable is defined.
check-variable-%:
	@[[ "${${*}}" ]] || (echo '*** Please define variable `${*}` ***' && exit 1)

# ----------------------------------------------------------------------------------------------------------------------
# Linting Targets
# ----------------------------------------------------------------------------------------------------------------------

.PHONY: lint lint-docker
lint: lint-markdown lint-shell lint-yaml lint-dockerfile lint-makefile lint-json lint-file
lint-docker: lint-markdown-docker lint-shell-docker lint-yaml-docker lint-dockerfile-docker lint-makefile-docker lint-json-docker lint-file-docker

.PHONY: lint-markdown lint-markdown-docker
lint-markdown:
	@scripts/lint-markdown.sh

lint-markdown-docker:
	$(call run-script-docker,${_DOCKER_MARKDOWNLINT_IMAGE},lint-markdown)

.PHONY: lint-shell lint-shell-docker
lint-shell:
	@scripts/lint-shell.sh

lint-shell-docker:
	$(call run-script-docker,${_DOCKER_SHELLCHECK_IMAGE},lint-shell)

.PHONY: lint-yaml lint-yaml-docker
lint-yaml:
	@scripts/lint-yaml.sh

lint-yaml-docker:
	$(call run-script-docker,${_DOCKER_YAMLLINT_IMAGE},lint-yaml)

.PHONY: lint-dockerfile lint-dockerfile-docker
lint-dockerfile:
	@scripts/lint-dockerfile.sh

lint-dockerfile-docker:
	$(call run-script-docker,${_DOCKER_HADOLINT_IMAGE},lint-dockerfile)

.PHONY: lint-makefile lint-makefile-docker
lint-makefile:
	@scripts/lint-makefile.sh

lint-makefile-docker:
	$(call run-script-docker,${_DOCKER_MAKEFILELINT_IMAGE},lint-makefile)

.PHONY: lint-json lint-json-docker
lint-json:
	@scripts/lint-json.sh

lint-json-docker:
	$(call run-script-docker,${_DOCKER_JSONLINT_IMAGE},lint-json)

.PHONY: lint-file lint-file-docker
lint-file:
	@scripts/lint-file.sh

lint-file-docker:
	$(call run-script-docker,${_DOCKER_FILELINT_IMAGE},lint-file)

# ----------------------------------------------------------------------------------------------------------------------
# Formatting Targets
# ----------------------------------------------------------------------------------------------------------------------

.PHONY: format format-docker
format: format-file format-shell format-markdown
format-docker: format-file-docker format-shell-docker format-markdown-docker

.PHONY: format-file format-file-docker
format-file:
	@scripts/format-file.sh

format-file-docker:
	$(call run-script-docker,${_DOCKER_FILELINT_IMAGE},format-file)

.PHONY: format-markdown format-markdown-docker
format-markdown:
	@scripts/format-markdown.sh

format-markdown-docker:
	$(call run-script-docker,${_DOCKER_MARKDOWNLINT_IMAGE},format-markdown)

.PHONY: format-shell format-shell-docker
format-shell:
	@scripts/format-shell.sh

format-shell-docker:
	$(call run-script-docker,${_DOCKER_SHFMT_IMAGE},format-shell)

.PHONY: format-yaml format-yaml-docker

format-yaml:
	@scripts/format-yaml.sh

format-yaml-docker: docker-tools
	$(call run-script-docker,${_DOCKER_TOOLS_IMAGE},format-yaml)

.PHONY: format-json format-json-docker

format-json:
	@scripts/format-json.sh

format-json-docker: docker-tools
	$(call run-script-docker,${_DOCKER_TOOLS_IMAGE},format-json)

# -------------------------------------------------------------------------------------------------
# Development Targets
# -------------------------------------------------------------------------------------------------

# dev-up: Start the local development environment with Tilt.
# Use FORCE=1 to recreate the self-signed TLS certificates
.PHONY: dev-up
dev-up: check-variable-CLUSTER_VERSION
	@./scripts/dev-up.sh $(CLUSTER_VERSION) $(FORCE)

# dev-down: Tears down the local development environment.
# Use FORCE=1 to delete local kind registry
.PHONY: dev-down
dev-down:
	@./scripts/dev-down.sh $(FORCE)

# dev-start: Resumes the local development environment.
.PHONY: dev-start
dev-start:
	@./scripts/dev-start.sh

# dev-stop: Pauses the local development environment.
.PHONY: dev-stop
dev-stop:
	@./scripts/dev-stop.sh

# Deploy ------------------------------------------------------------------------------------------

# deploy-helm: Install or update deployment using helm
.PHONY: deploy-helm
deploy-helm: check-variable-IMAGE_TAG_NAME check-variable-CLUSTER_NAME check-variable-NAMESPACE
	@helm upgrade -i ${PROJECT_NAME} helm_chart -f helm_chart/values.yaml \
	 --namespace ${NAMESPACE} \
	 --set image.repository=${PROJECT_NAME} \
	 --set image.tag=${IMAGE_TAG_NAME}
