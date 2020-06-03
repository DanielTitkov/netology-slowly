package configs

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config holds app configuration
type Config struct {
	Port            string
	ShutdownTimeout int `yaml:"shutdownTimeout"`
	MaxSlowTimeout  int `yaml:"maxSlowTimeout"`
}

// ReadConfigs reads config from yaml file
func ReadConfigs(path string) (Config, error) {
	var cfg Config
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
