package acceptance

import "time"

type Config struct {
	TimeoutMinutes     time.Duration `json:"timeout_in_minutes"`
	KuboDeploymentName string        `json:"kubo_deployment_name"`
	Kubo               KuboConfig    `json:"kubo"`
}

func NewConfig() (Config, error) {
	return Config{
		TimeoutMinutes: 30,
	}, nil
}

type KuboConfig struct {
	ClusterName  string `json:"cluster_name"`
	APIServerURL string `json:"api_server_url"`
	CACert       string `json:"ca_cert"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}
