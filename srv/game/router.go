package game

import "reflect"
import "strconv"
import "log"
import "github.com/szgut/www24/srv/core"

type CommandHandle func(team core.Team, params []string) core.CommandResult

type Router struct {
	commands map[string]CommandHandle
}

type parser func(param string) (interface{}, error)

var parse = map[reflect.Kind]parser{
	reflect.String: func(param string) (interface{}, error) {
		return param, nil
	},
	reflect.Int: func(param string) (interface{}, error) {
		i, err := strconv.ParseInt(param, 10, 0)
		return int(i), err
	},
	reflect.Float64: func(param string) (interface{}, error) {
		return strconv.ParseFloat(param, 64)
	},
}

func NewRouter(methods map[string]interface{}) *Router {
	commands := make(map[string]CommandHandle)
	for name, method := range methods {
		typ := reflect.TypeOf(method)
		val := reflect.ValueOf(method)
		paramc := typ.NumIn() - 1
		parsers := []parser{}
		for i := 0; i < paramc; i++ {
			kind := typ.In(i + 1).Kind()
			parser, ok := parse[kind]
			if !ok {
				log.Fatalf("Unknown argument kind %v in command handler %v\n", kind, name)
			}
			parsers = append(parsers, parser)
		}
		commands[name] = func(team core.Team, params []string) core.CommandResult {
			if len(params) != paramc {
				return core.NewErrResult(core.BadFormatError())
			}
			args := []reflect.Value{reflect.ValueOf(team)}
			for i := range params {
				arg, err := parsers[i](params[i])
				if err != nil {
					return core.NewErrResult(core.BadFormatError())
				}
				args = append(args, reflect.ValueOf(arg))
			}
			return val.Call(args)[0].Interface().(core.CommandResult)
		}
	}
	return &Router{commands: commands}
}

func (self *Router) Execute(team core.Team, cmd core.Command) core.CommandResult {
	handle, ok := self.commands[cmd.Name]
	if !ok {
		return core.NewErrResult(core.UnknownCommandError())
	}
	return handle(team, cmd.Params)
}
