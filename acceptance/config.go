package acceptance

import (
	"encoding/json"
	"io/ioutil"
	"time"

	. "github.com/onsi/gomega"
)

type Config struct {
	TimeoutMinutes time.Duration `json:"timeout_in_minutes"`
	APIServerURL   string        `json:"api_server_url"`
	CACert         string        `json:"ca_cert"`
	Username       string        `json:"username"`
	Password       string        `json:"password"`
}

func NewConfig(path string) Config {
	rawConfig, err := ioutil.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())

	var config Config
	err = json.Unmarshal(rawConfig, &config)
	Expect(err).NotTo(HaveOccurred())

	if config.TimeoutMinutes == 0 {
		config.TimeoutMinutes = 5
	}

	return config
}
