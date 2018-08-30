#!/usr/bin/env bash

set -eu

pushd bbr-binary-release
  tar xvf ./*.tar
  export BBR_BINARY_PATH="$PWD/releases/bbr"
popd

pushd "bbl-state/$BBL_STATE_DIR"
  eval "$(bbl print-env)"
popd

export GOPATH="$PWD"
export PATH="$PATH:$GOPATH/bin"
export CONFIG_PATH="$PWD/k-drats-config/$CONFIG_PATH"

pushd src/github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests
  scripts/_run_acceptance_tests.sh
popd
