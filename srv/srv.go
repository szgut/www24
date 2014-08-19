package main

import "fmt"
import "net"
import "log"
import "flag"
import "github.com/szgut/www24/srv/backend"
import "github.com/szgut/www24/srv/core"
import "github.com/szgut/www24/srv/game"
import "github.com/szgut/www24/srv/score"
import "github.com/szgut/www24/srv/limit"

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func flags() Args {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	var args Args
	flag.StringVar(&args.configPath, "conf", "conf.yml", "config file path")
	flag.StringVar(&args.task, "task", "task", "task name")
	flag.BoolVar(&args.startNew, "start-new", false, "try to recover")
	flag.Parse()
	return args
}

type Args struct {
	configPath string
	task       string
	startNew   bool
}

func main() {
	args := flags()

	config, taskConfig := getConfigs(args.configPath, args.task)
	ss := score.NewStorage(config.Path, args.task)
	if args.startNew {
		ss.Initialize(config.GetTeams())
		return
	} else {
		ss.ReadScores()
	}
	game := createGame(taskConfig, ss.StartRound(), ss)
	bend := backend.StartNew(taskConfig.TickInterval, game)
	dos := limit.NewDoS(config.Connections)

	l := listen("0.0.0.0", taskConfig.Port)
	defer l.Close()
	for {
		conn, err := l.Accept()
		check(err)
		if !dos.Accept(conn) {
			conn.Close()
			continue
		}
		go func() {
			defer conn.Close()
			defer dos.Release(conn)
			handleConnection(NewProto(conn), config, bend)
		}()
	}
}

func getConfigs(path string, task string) (*Config, *TaskConfig) {
	config, err := ReadConfig(path)
	check(err)
	taskConfig, ok := config.Tasks[task]
	if !ok {
		log.Fatalf("task \"%s\" not found in config", task)
	}
	return config, &taskConfig
}

func listen(host string, port int) net.Listener {
	hostport := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", hostport)
	check(err)
	log.Println("Listening on " + hostport)
	return listener
}

func createGame(taskConfig *TaskConfig, startRound int, ss score.Storage) core.Game {
	cons, err := game.RegistryFind(taskConfig.Game)
	check(err)
	game := cons(taskConfig.Params, startRound, ss)
	return limit.Throttler(taskConfig.Commands, game)
}

func handleConnection(proto Proto, auth Authenticator, bend backend.Backend) {
	login, pass, err := proto.Credentials()
	if err != nil {
		return
	}
	team := auth.Authenticate(login, pass)
	if team == nil {
		err := core.AuthenticationFailedError()
		proto.Write(&err)
	} else {
		proto.Write(nil)
		log.Println("Team", team, "authenticated")
		authenticated(proto, *team, bend)
	}
}

func authenticated(proto Proto, team core.Team, bend backend.Backend) {
	defer log.Println("Team", team, "disconnected")

	waitOk := func(msg string) {
		proto.writeln(msg)
		bend.Wait()
		proto.Write(nil)
	}

	for cmd := proto.Command(); cmd != nil; cmd = proto.Command() {
		if cmd.Name == "WAIT" {
			proto.Write(nil)
			waitOk("WAITING")
		} else {
			result := bend.Command(team, *cmd)
			if proto.Write(result.Err, result.Params...) != nil {
				return
			}
			if result.Err.ShouldWait() {
				waitOk("FORCED_WAITING")
			}
		}
	}
}
