package limit

import "log"
import "net"

type DoS interface {
	Accept(conn net.Conn) bool
	Release(conn net.Conn)
}

type dos struct {
	limit int
	used  map[string]int
}

func NewDoS(limit int) DoS {
	return &dos{limit: limit, used: make(map[string]int)}
}

func remoteIP(conn net.Conn) string {
	if host, _, err := net.SplitHostPort(conn.RemoteAddr().String()); err != nil {
		return ""
	} else {
		return host
	}
}

func (d *dos) Accept(conn net.Conn) bool {
	ip := remoteIP(conn)
	if d.used[ip] < d.limit {
		d.used[ip]++
		log.Printf("Accepted %s (%d)\n", ip, d.used[ip])
		return true
	}
	return false
}

func (d *dos) Release(conn net.Conn) {
	ip := remoteIP(conn)
	d.used[ip]--
	log.Printf("Released %s (%d)\n", ip, d.used[ip])
}
