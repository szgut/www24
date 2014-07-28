package game

import "fmt"
import "github.com/szgut/www24/srv/core"
import "github.com/szgut/www24/srv/score"

type Cons func(config *core.Config, ss score.Storage) core.Game

func RegistryFind(name string) (Cons, error) {
	cons, ok := map[string]Cons{
		"simple": newSimpleGame,
	}[name]
	if !ok {
		return nil, fmt.Errorf("Game %s not found in registry", name)
	}
	return cons, nil
}
