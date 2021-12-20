package tcpserver

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

const debug = true

type server struct {
	conns map[string]Conn
	l     sync.Mutex
}

func New() *server {
	return &server{
		conns: map[string]Conn{},
	}
}

type Req struct {
	C      Conn
	Text   string
	Closed bool
}

type Conn struct {
	net.Conn
}

func (s *server) ReadRequests(reqs chan Req) {
	for req := range reqs {
		if req.Text == "STOP" {
			s.deleteConnection(req.C)
			fmt.Printf("Stopping %s\n", req.C.RemoteAddr().String())
			close(reqs)
			break
		}
		if strings.HasPrefix(req.Text, "HELLO ") {
			s.addConnection(req.C, req.Text)
			continue
		}
		if req.Closed {
			s.deleteConnection(req.C)
			close(reqs)
		}
		s.broadcast(req)
	}
}

func (s *server) deleteConnection(c net.Conn) {
	s.l.Lock()
	defer s.l.Unlock()
	name := s.findClientNameByRemoteAddr(c.RemoteAddr().String())
	if name != "" {
		delete(s.conns, name)
	}
	s.printDebug()
}

func (s *server) findClientNameByRemoteAddr(remoteAddr string) string {
	for name, ci := range s.conns {
		if ci.RemoteAddr().String() == remoteAddr {
			return name
		}
	}
	return ""
}

func (s *server) broadcast(msg Req) {
	if strings.TrimSpace(msg.Text) == "" {
		return
	}
	s.l.Lock()
	defer s.l.Unlock()
	for _, c := range s.conns {
		if c.RemoteAddr().String() != msg.C.RemoteAddr().String() {
			_, err := c.Write([]byte(msg.Text))
			if err != nil {
				fmt.Println("broadcast: ", err)
			}
		}
	}
}

// adds a connection and return its name if succeed
func (s *server) addConnection(c net.Conn, temp string) string {
	clientName := strings.TrimSpace(strings.TrimPrefix(temp, "HELLO "))
	s.l.Lock()
	defer s.l.Unlock()
	connWithThisClientName, ok := s.conns[clientName]
	if ok && connWithThisClientName.RemoteAddr().String() != c.RemoteAddr().String() {
		_, _ = c.Write([]byte(fmt.Sprintf(`Name "%s" is busy by other connection: %s`,
			clientName, connWithThisClientName.RemoteAddr().String())))
		return ""
	}
	clientNameWithThisAddr := s.findClientNameByRemoteAddr(c.RemoteAddr().String())
	if clientNameWithThisAddr != "" {
		_, _ = c.Write([]byte(fmt.Sprintf(`Your connection always has a name "%s"`, clientNameWithThisAddr)))
		return ""
	}
	s.conns[clientName] = Conn{Conn: c}
	s.printDebug()
	return clientName
}

func (s *server) printDebug() {
	if debug {
		fmt.Println("Server's connection state:", s.conns)
	}
}
