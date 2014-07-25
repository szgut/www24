package core

import "io/ioutil"
import "gopkg.in/yaml.v1"

type Config struct {
	Port        int
	Path        string
	Teams       map[string]string
	Interval    int
	Commands    int
	Connections int
	Game        string
	Task        string
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
	correctPass, ok := config.Teams[login]
	if ok && pass == correctPass {
		return &Team{login}
	} else {
		return nil
	}
}

func (config *Config) ListTeams() []Team {
	teams := make([]Team, 0)
	for login, _ := range config.Teams {
		teams = append(teams, Team{login})
	}
	return teams
}
