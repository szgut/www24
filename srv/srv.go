package main

import "fmt"
import "net"
import "os"
import "log"
import "github.com/szgut/www24/srv/backend"
import "github.com/szgut/www24/srv/core"
import "github.com/szgut/www24/srv/game"
import "github.com/szgut/www24/srv/score"

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func listen(host string, port int) net.Listener {
	hostport := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", hostport)
	check(err)
	log.Println("Listening on " + hostport)
	return listener
}

func configPath() string {
	if len(os.Args) != 2 {
		fmt.Printf("%s <config path>\n", os.Args[0])
		os.Exit(1)
	}
	return os.Args[1]
}

func initLogger() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

func createGame(config *core.Config, ss score.Storage) game.Game {
	cons, err := game.RegistryFind(config.Game)
	check(err)
	game := cons(config, ss)
	return Throttler(config.Commands, game)
}

func main() {
	initLogger()
	config, err := core.ReadConfig(configPath())
	check(err)
	log.Println("Teams:", config.ListTeams())

	ss := score.NewStorage(config.Path, config.Task)
	fmt.Println(ss)

	l := listen("localhost", config.Port)
	defer l.Close()

	bend := backend.StartNew(config.Interval, createGame(config, ss))
	dos := NewDoS(config.Connections)
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

func handleConnection(proto Proto, auth core.Authenticator, bend backend.Backend) {
	login, pass, err := proto.Credentials()
	if err != nil {
		return
	}
	team := auth.Authenticate(login, pass)
	if team == nil {
		proto.Write(core.AuthenticationFailedError())
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
			proto.Write(result.Err, result.Params)
			if result.Err.ShouldWait() {
				waitOk("FORCED_WAITING")
			}
		}
	}
}
