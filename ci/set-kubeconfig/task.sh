#!/usr/bin/env bash

set -eu
set -o pipefail

pushd "bbl-state/$BBL_STATE_DIR"
  eval "$(bbl print-env)"
popd

api_server_ip="$(jq -r .lb < terraform/metadata)"
ca_cert="$(bosh int <(credhub get -n "${BOSH_DIRECTOR_NAME}/${BOSH_DEPLOYMENT}/tls-kubernetes" --output-json) --path=/value/ca)"
password="$(bosh int <(credhub get -n "${BOSH_DIRECTOR_NAME}/${BOSH_DEPLOYMENT}/kubo-admin-password" --output-json) --path=/value)"
cluster_name=k-drats
echo "$ca_cert" > ca_cert
chmod 600 ca_cert

kubectl config set-cluster $cluster_name --server=https://$api_server_ip:8443 --certificate-authority=ca_cert --embed-certs=true
kubectl config set-credentials $KUBO_USERNAME --token=$password
kubectl config set-context $cluster_name --cluster=$cluster_name --user=$KUBO_USERNAME
kubectl config use-context $cluster_name

cp ~/.kube/config kubeconfig/config