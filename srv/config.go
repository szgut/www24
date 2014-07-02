package main

import "fmt"
import "io/ioutil"
import "gopkg.in/yaml.v1"

type Config struct {
	Path string
	Teams map[string]string
	Interval int
}

func ReadConfig(path string) (*Config, error) {
	fmt.Println(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", config)
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
