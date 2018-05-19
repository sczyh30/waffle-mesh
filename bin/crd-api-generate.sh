#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

CODEGEN_PKG=./vendor/k8s.io/code-generator

${CODEGEN_PKG}/generate-groups.sh all \
  github.com/sczyh30/waffle-mesh/brain/k8s/gen github.com/sczyh30/waffle-mesh/brain/k8s \
  crd:v1