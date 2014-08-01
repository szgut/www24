package main

import "io/ioutil"
import "gopkg.in/yaml.v1"

import "github.com/szgut/www24/srv/core"
import "github.com/szgut/www24/srv/game"

type Config struct {
	Path        string `yaml:"db_path"`
	Teams       map[string]string
	Connections int `yaml:"max_connections"`
	Tasks       map[string]TaskConfig
}

type TaskConfig struct {
	Port         int
	Game         string
	TickInterval int `yaml:"tick_interval"`
	Commands     int `yaml:"commands_limit"`
	Params       game.Params
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
	Authenticate(login, pass string) *core.Team
}

func (config *Config) Authenticate(login, pass string) *core.Team {
	correctPass, ok := config.Teams[login]
	if ok && pass == correctPass {
		team := core.NewTeam(login)
		return &team
	} else {
		return nil
	}
}

func (config *Config) GetTeams() []core.Team {
	teams := make([]core.Team, 0)
	for login, _ := range config.Teams {
		teams = append(teams, core.NewTeam(login))
	}
	return teams
}
