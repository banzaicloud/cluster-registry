# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"
REPO_ROOT=$(shell git rev-parse --show-toplevel)

.PHONY: ensure-tools
ensure-tools:
	@echo "ensure tools"
	@scripts/download-deps.sh

.PHONY: manifests
manifests: ensure-tools ## Generate manifests
	${REPO_ROOT}/bin/controller-gen $(CRD_OPTIONS) paths="./..." object:headerFile="hack/boilerplate.go.txt" output:crd:artifacts:config=crd

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: manifests
