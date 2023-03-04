package config

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

// Config ..
type Config struct {
	Port           string `yaml:"port"`
	Host           string `yaml:"host"`
	LogLevel       string `yaml:"log_level"`
	TraceDirectory string `yaml:"trace_directory"`
}

// New config ..
func New() *Config {
	return &Config{
		Port:           ":8080",
		Host:           "localhost",
		LogLevel:       "debug",
		TraceDirectory: "/tmp",
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
