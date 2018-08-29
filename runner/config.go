package runner

import "time"

type Config struct {
	Timeout time.Duration
}

func NewConfig() (Config, error) {
	return Config{
		Timeout: 30 * time.Minute,
	}, nil
}
