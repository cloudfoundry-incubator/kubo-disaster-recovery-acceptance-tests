# kubo-disaster-recovery-acceptance-tests (k-DRATs)

Tests a given CFCR K8s cluster can be backed up and restored using [`bbr`](https://github.com/cloudfoundry-incubator/bosh-backup-and-restore).

The acceptance test suite provides hooks around `bbr director backup` and `bbr director restore`.

## Dependencies

1. Install [Golang](https://golang.org/doc/install)
1. Install [`ginkgo` CLI](https://github.com/onsi/ginkgo#set-me-up)
1. Install [`kubectl` CLI](https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl)
1. Download [`bbr` CLI](https://github.com/cloudfoundry-incubator/bosh-backup-and-restore/releases) and add it the `PATH`

## Running k-DRATs in your pipelines

We encourage you to use our [`run-k-drats-master` CI task](https://github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/tree/master/ci/run-k-drats) to run k-DRATS in your Concourse pipeline.

Please refer to our k-drats [pipeline definition](https://github.com/cloudfoundry-incubator/backup-and-restore-ci/blob/master/pipelines/k-drats/pipeline.yml) for a working example.

## Running k-DRATs locally

1. Spin up a CFCR deployment
   - [kubo-deployment](https://github.com/cloudfoundry-incubator/kubo-deployment) is supported. Make sure you apply the [enable-bbr opsfile](https://github.com/cloudfoundry-incubator/kubo-deployment/blob/master/manifests/ops-files/enable-bbr.yml) at deploy time to ensure the backup and restore scripts are enabled.
1. Clone this repo
    ```bash
    $ go get github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests
    $ cd $GOPATH/src/github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests
    ```
1. Create an `config.json` file, for example:
    ```json
    {
      "run_test_case_deployment": true,
      "run_test_case_etcd_cluster": true,
      "timeout_in_minutes": 30,
      "api_server_url": "<k8s_api_url>",
      "ca_cert": "<k8s_ca_cert>",
      "username": "<k8s_username>",
      "password": "<k8s_password>"
    }
    ```
1. Export `CONFIG_PATH` to be path to `config.json` file you just created.
1. Export the following BOSH environment variables
   - `BOSH_ENVIRONMENT` - URL of BOSH Director which has deployed the CFCR cluster
   - `BOSH_CLIENT` - BOSH Director username
   - `BOSH_CLIENT_SECRET` - BOSH Director password
   - `BOSH_CA_CERT` - BOSH Director's CA cert content
   - `BOSH_ALL_PROXY` - optional, set the proxy to be used in case the BOSH director is behind a jumpbox 
   - `BOSH_DEPLOYMENT` - name of the CFCR deployment to backup and restore
1. Run acceptance tests
    ```bash
    $ ./scripts/_run_acceptance_tests.sh
    ```

## Config Variables

* `api_server_url` - Url of K8s api server
* `ca_cert` - K8s CA cert
* `username` - K8s username
* `password` - K8s password
* `timeout_in_minutes` - default ginkgo `Eventually` timeout in minutes
* `run_test_case_<test-case-name>` - flag for whether to run a given testcase, if omitted defaults to `false`

## Contributing to k-DRATs

k-DRATS runs a collection of test cases against a CFCR cluster.

Test cases should be used for checking that K8s data has been backed up and restored correctly. For example, if two  workflows are deployed before `bbr director backup`, and the workflows are removed after taking the backup. Then after a successful `bbr director restore`, workflows will be restored back to their original state.

To add extra test cases, create a new test case that implements the [`TestCase` interface](https://github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/blob/master/acceptance/testcase.go).

The methods that need to be implemented are:
* `Name() string` - should return name of the test case.
* `BeforeBackup(Config)` - runs before the backup is taken, and should create state in the K8s cluster to be backed up.
* `AfterBackup(Config)` - runs after the backup is complete but before the restore is started.
* `AfterRestore(Config)` - runs after the restore is complete, and should assert that the state in the restored K8s cluster matches that created in `BeforeBackup(Config)`.
* `Cleanup(Config)` - should clean up the state created in the K8s cluster through the test.

`Config` contains the config for accessing the target K8s. *Note*: the use of this config is optional, `kubectl` is already configured to access the target K8s cluster when a testcase runs.

### Creating a new test case

1. Create a new test case in the [testcases package](https://github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/tree/master/testcases).
1. Add the newly created test case to the list of `availableTestCases` in [`acceptance_suite_test.go`](https://github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/blob/master/acceptance/acceptance_suite_test.go).
