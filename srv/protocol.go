package main

import "io"
import "fmt"
import "bufio"
import "strings"
import "github.com/szgut/www24/srv/core"

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

func (proto *Proto) Write(cmdErr *core.CommandError, lines ...[]interface{}) error {
	if cmdErr == nil {
		if err := proto.writeln("OK"); err != nil {
			return err
		}
	} else {
		if err := proto.writeln(fmt.Sprintf("FAILED %d %s", cmdErr.Id, cmdErr.Desc)); err != nil {
			return err
		}
	}
	for _, values := range lines {
		if err := proto.writeln(values...); err != nil {
			return err
		}
	}
	return nil
}

func (proto *Proto) Command() *core.Command {
	for {
		line, err := proto.readln()
		if err != nil {
			return nil
		}
		words := strings.Fields(line)
		if len(words) > 0 {
			return &core.Command{strings.ToUpper(words[0]), words[1:]}
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
