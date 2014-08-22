package game

import "fmt"
import "github.com/szgut/www24/srv/core"
import "github.com/szgut/www24/srv/score"

type Params map[string]string
type Cons func(params Params, firstRount int, teams []core.Team, ss score.Storage) core.Game

func RegistryFind(name string) (Cons, error) {
	cons, ok := map[string]Cons{
		//"simple": newSimpleGame,
		//"fields": newFieldsGame,
		"star": newStarGame,
	}[name]
	if !ok {
		return nil, fmt.Errorf("Game %s not found in registry", name)
	}
	return cons, nil
}
