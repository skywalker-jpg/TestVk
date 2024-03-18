package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type DBConfig struct {
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	SSLMode  string `yaml:"ssl_mode"`
	User     string `yaml:"user"`
	Port     string `yaml:"port"`
}

type LoggerConfig struct {
	Sink  string `yaml:"sink"`
	Level string `yaml:"level"`
}

type AuthToken struct {
	Admin string `yaml:"admin"`
}

type AppConfig struct {
	DB     DBConfig     `yaml:"db"`
	Logger LoggerConfig `yaml:"logger"`
	Auth   AuthToken    `yaml:"auth_token"`
}

func NewConfig(path string) (*AppConfig, error) {
	yamlConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var appConfig AppConfig

	if err := yaml.Unmarshal(yamlConfig, &appConfig); err != nil {
		return nil, err
	}
	return &appConfig, nil
}
