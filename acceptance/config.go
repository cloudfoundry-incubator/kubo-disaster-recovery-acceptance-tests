package acceptance

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
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
	config := Config{
		TimeoutMinutes: 5,
	}
	if path == "" {
		return config
	}
	rawConfig, err := ioutil.ReadFile(filepath.Clean(path))
	Expect(err).NotTo(HaveOccurred())

	err = json.Unmarshal(rawConfig, &config)
	Expect(err).NotTo(HaveOccurred())

	return config
}
