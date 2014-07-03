package main

import "io"
import "fmt"
import "bufio"
import "strings"

type Proto struct {
	conn   io.ReadWriteCloser
	reader *bufio.Reader
}

func NewProto(conn io.ReadWriteCloser) Proto {
	return Proto{conn: conn, reader: bufio.NewReader(conn)}
}

func (proto *Proto) readln() (string, error) {
	if line, err := proto.reader.ReadString('\n'); err == nil {
		return line[:len(line)-1], nil
	} else {
		return "", err
	}
}

func (proto *Proto) writeln(values ...interface{}) error {
	_, err := fmt.Fprintln(proto.conn, values...)
	return err
}

func (proto *Proto) Write(err *CommandError, lines ...[]interface{}) {
	if err == nil {
		proto.writeln("OK")
	} else {
		proto.writeln(fmt.Sprintf("FAILED %d %s", err.Id, err.Desc))
	}
	for _, values := range lines {
		proto.writeln(values...)
	}
}

func (proto *Proto) Command() *Command {
	for {
		line, err := proto.readln()
		if err != nil {
			return nil
		}
		words := strings.Fields(line)
		if len(words) > 0 {
			return &Command{strings.ToUpper(words[0]), words[1:]}
		}
	}
}

func (proto *Proto) Credentials() (login, pass string, err error) {
	proto.writeln("LOGIN")
	login, err = proto.readln()
	if err != nil {
		return "", "", err
	}
	proto.writeln("PASS")
	pass, err = proto.readln()
	if err != nil {
		return "", "", err
	}
	return login, pass, nil
}
