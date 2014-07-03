package main

import "io/ioutil"
import "gopkg.in/yaml.v1"

type Config struct {
	Port     int
	Path     string
	Teams    map[string]string
	Interval int
	Commands int
}

func ReadConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type Authenticator interface {
	Authenticate(login, pass string) *Team
}

func (config *Config) Authenticate(login, pass string) *Team {
	if config.Teams[login] == pass {
		return &Team{login}
	} else {
		return nil
	}
}
