package config

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

// Config ..
type Config struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

// New ..
func New() *Config {
	return &Config{
		Port: ":9090",
		Host: "localhost",
	}
}

// ParseFile ..
func (c *Config) ParseFile(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return err
	}
	return nil
}
