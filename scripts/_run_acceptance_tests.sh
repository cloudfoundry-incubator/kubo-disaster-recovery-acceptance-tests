#!/usr/bin/env bash

set -euo pipefail

SCRIPTS_DIR="$(dirname $0)"

ginkgo -v --trace "${SCRIPTS_DIR}/../acceptance"
