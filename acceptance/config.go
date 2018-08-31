package acceptance

import "time"

type Config struct {
	TimeoutMinutes time.Duration `json:"timeout_in_minutes"`
	Kubo           KuboConfig    `json:"kubo"`
}

type KuboConfig struct {
	ClusterName  string `json:"cluster_name"`
	APIServerURL string `json:"api_server_url"`
	CACert       string `json:"ca_cert"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}
