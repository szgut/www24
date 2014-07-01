package main

import (
	"fmt"
	"net"
	"os"
)

const (
	LISTEN_HOST = "localhost"
	LISTEN_PORT = 3333
)

func listen(host string, port int) net.Listener {
	hostport := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", hostport)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on " + hostport)
	return listener
}

func main() {
	l := listen(LISTEN_HOST, LISTEN_PORT)
	defer l.Close()

	config := ReadConfig("")
	bch, wait := StartBackend(&SimpleGame{})
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn, &config, bch, wait)
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
		proto.Write(&CommandError{Desc: "authentication failed"})
	} else {
		proto.Write(nil)
		authenticated(proto, team, bch, wait)
	}
	fmt.Println(*team)
}

func authenticated(proto Proto, team Team, bch chan<- CommandMessage, wait func()) {
	fmt.Println(*team, "connected")
	defer fmt.Println(*team, "disconnected")

	for cmd := proto.Command(); cmd != nil; cmd = proto.Command() {
		if cmd.Name == "WAIT" {
			proto.Write(nil, []interface{}{"WAITING"})
			wait()
			proto.Write(nil)
		} else {
			rch := make(chan ResultMessage)
			bch <- CommandMessage{team, *cmd, rch}
			result := <-rch
			proto.Write(result.Err, result.Params)
		}
	}
}
