#!/usr/bin/env bash

set -eu

pushd bbr-binary-release
  tar xvf ./*.tar
  export BBR_BINARY_PATH="$PWD/releases/bbr"
popd

export GOPATH="$PWD"
export PATH="$PATH:$GOPATH/bin"
INTEGRATION_CONFIG_PATH="$PWD/$INTEGRATION_CONFIG_PATH"

pushd src/github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests
  scripts/_run_acceptance_tests.sh
popd
