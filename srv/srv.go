package main

import "fmt"
import "net"
import "os"
import "log"

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

func main() {
	initLogger()
	config, err := ReadConfig(configPath())
	check(err)

	l := listen("localhost", config.Port)
	defer l.Close()

	bch, wait := StartBackend(config)
	dos := NewDoS(config.Connections)
	for {
		conn, err := l.Accept()
		check(err)
		if dos.Accept(conn) {
			go handleConnection(conn, dos, config, bch, wait)
		} else {
			conn.Close()
		}
	}
}

func handleConnection(conn net.Conn, dos DoS, auth Authenticator, bch chan<- CommandMessage, wait func()) {
	defer conn.Close()
	defer dos.Release(conn)
	proto := NewProto(conn)

	login, pass, err := proto.Credentials()
	if err != nil {
		return
	}
	team := auth.Authenticate(login, pass)
	if team == nil {
		proto.Write(AuthenticationFailedError())
	} else {
		proto.Write(nil)
		log.Println("Team", team, conn.RemoteAddr(), "authenticated")
		authenticated(proto, *team, bch, wait)
	}
}

func authenticated(proto Proto, team Team, bch chan<- CommandMessage, wait func()) {
	defer log.Println("Team", team, "disconnected")

	waitOk := func(msg string) {
		proto.writeln(msg)
		wait()
		proto.Write(nil)
	}

	for cmd := proto.Command(); cmd != nil; cmd = proto.Command() {
		if cmd.Name == "WAIT" {
			proto.Write(nil)
			waitOk("WAITING")
		} else {
			rch := make(chan ResultMessage)
			bch <- CommandMessage{team, *cmd, rch}
			result := <-rch
			proto.Write(result.Err, result.Params)
			if result.Err.ShouldWait() {
				waitOk("FORCED_WAITING")
			}
		}
	}
}
