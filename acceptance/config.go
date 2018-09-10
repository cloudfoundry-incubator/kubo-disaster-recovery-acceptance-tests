package acceptance

import "time"

type Config struct {
	TimeoutMinutes time.Duration `json:"timeout_in_minutes"`
	APIServerURL   string        `json:"api_server_url"`
	CACert         string        `json:"ca_cert"`
	Username       string        `json:"username"`
	Password       string        `json:"password"`
}
