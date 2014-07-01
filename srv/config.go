package main

type Config struct {
}

func ReadConfig(path string) Config {
	return Config{}
}

type Authenticator interface {
	Authenticate(login, pass string) Team
}

func (config *Config) Authenticate(login, pass string) Team {
	return Team(&login)
}
