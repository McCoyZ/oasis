#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

rm -rf ./pkg/client

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
bash ./hack/generate-groups.sh "deepcopy,client,informer,lister" \
  zmc.io/oasis/pkg/client zmc.io/oasis/pkg/apis \
  "servicemesh:v1alpha1" \
  --output-base "../.." \
  --go-header-file "./hack/boilerplate.go.txt"

# To use your own boilerplate text append:
#   --go-header-file "${SCRIPT_ROOT}"/hack/custom-boilerplate.go.txt