package main

import (
	"fmt"
	"net"
	"os"
)

func check(err error) {
	if err != nil {
		fmt.Println("Fatal error:", err.Error())
		os.Exit(1)
	}
}

func listen(host string, port int) net.Listener {
	hostport := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", hostport)
	check(err)
	fmt.Println("Listening on " + hostport)
	return listener
}

func configPath() string {
	if len(os.Args) != 2 {
		fmt.Printf("%s <config path>\n", os.Args[0])
		os.Exit(1)
	}
	return os.Args[1]
}

func main() {
	config, err := ReadConfig(configPath())
	check(err)

	l := listen("localhost", config.Port)
	defer l.Close()

	bch, wait := StartBackend(config)
	for {
		conn, err := l.Accept()
		check(err)
		go handleConnection(conn, config, bch, wait)
	}
}

func handleConnection(conn net.Conn, auth Authenticator, bch chan<- CommandMessage, wait func()) {
	defer conn.Close()
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
		authenticated(proto, *team, bch, wait)
	}
	fmt.Println(team)
}

func authenticated(proto Proto, team Team, bch chan<- CommandMessage, wait func()) {
	fmt.Println(team, "connected")
	defer fmt.Println(team, "disconnected")

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
