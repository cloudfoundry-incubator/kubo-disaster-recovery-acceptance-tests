package testcase

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/command"
	. "github.com/onsi/gomega"
)

type EtcdCluster struct{}

func (EtcdCluster) Name() string {
	return "etcd_cluster"
}

func NewEtcdCluster() EtcdCluster {
	return EtcdCluster{}
}

func (EtcdCluster) BeforeBackup(Config) {}

func (EtcdCluster) AfterBackup(Config) {}

func (EtcdCluster) AfterRestore(Config) {
	checkEtcdClusterHealth()

	instances := getNumberOfEtcdInstances()

	leaderGUIDs := make(map[int]string)
	for index := 0; index < instances; index++ {
		leaderGUIDs[index] = getEtcdLeader(index)
	}
	for _, leaderGUID := range leaderGUIDs {
		Expect(leaderGUID).To(Equal(leaderGUIDs[0]), fmt.Sprintf("split brain etcd cluster: %v", leaderGUIDs))
	}
}

func (EtcdCluster) Cleanup(Config) {}

func checkEtcdClusterHealth() {
	command.RunSuccessfully(
		"bosh ssh to run etcdctl cluster-health",
		"bosh", "ssh", "master/0",
		"-c", "\"sudo ETCDCTL_API=2 /var/vcap/jobs/etcd/bin/etcdctl cluster-health\"",
	)
}

func getNumberOfEtcdInstances() int {
	manifestSession := command.RunSuccessfullyWithoutStream(
		"bosh manifest",
		"bosh", "manifest",
	)

	manifestPath, err := ioutil.TempFile("", "k-drats-etcd-cluster")
	Expect(err).NotTo(HaveOccurred())
	err = ioutil.WriteFile(manifestPath.Name(), manifestSession.Out.Contents(), 644)
	Expect(err).NotTo(HaveOccurred())

	interpolateSession := command.RunSuccessfullyWithoutStream(
		"bosh interpolate",
		"bosh", "int", "--path=/instance_groups/name=master/instances", manifestPath.Name(),
	)

	instancesLines := strings.Split(string(interpolateSession.Out.Contents()), "\n")
	Expect(err).NotTo(HaveOccurred())
	instances, err := strconv.Atoi(instancesLines[0])
	Expect(err).NotTo(HaveOccurred())

	return instances
}

func getEtcdLeader(index int) string {
	session := command.RunSuccessfully(
		"bosh ssh to run etcdctl member list",
		"bosh", "ssh", fmt.Sprintf("master/%d", index),
		"-c", "\"sudo ETCDCTL_API=2 /var/vcap/jobs/etcd/bin/etcdctl member list\"",
	)

	var leaderInstanceGUID string
	lines := strings.Split(string(session.Out.Contents()), "\n")
	for _, line := range lines {
		if strings.Contains(line, "isLeader=true") {
			re, err := regexp.Compile(`name=(\S+)`)
			Expect(err).NotTo(HaveOccurred())

			if matches := re.FindStringSubmatch(line); matches != nil {
				leaderInstanceGUID = matches[1]
				break
			}
		}
	}
	Expect(leaderInstanceGUID).NotTo(BeEmpty())

	return leaderInstanceGUID
}
