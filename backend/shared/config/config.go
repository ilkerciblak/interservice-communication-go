// Package config
package config

import "os"

type config struct {
	DataDir string
}

func Config() *config {
	return &config{
		DataDir: getEnvorDefault("DATA_DIR", "./data"),
	}
}

func getEnvorDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def

}
