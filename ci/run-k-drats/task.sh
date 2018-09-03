#!/usr/bin/env bash

set -eu

pushd bbr-binary-release
  tar xvf ./*.tar
  export PATH="$PATH:$PWD/releases"
popd

pushd "bbl-state/$BBL_STATE_DIR"
  eval "$(bbl print-env)"
popd

credhub login
./kubo-deployment/bin/set_kubeconfig "${BOSH_DIRECTOR_NAME}/${BOSH_DEPLOYMENT}" $(jq -r .api_server_url <"k-drats-config/config.json")

export GOPATH="$PWD"
export PATH="$PATH:$GOPATH/bin"
export CONFIG_PATH="$PWD/k-drats-config/$CONFIG_PATH"

pushd src/github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests
  scripts/_run_acceptance_tests.sh
popd
