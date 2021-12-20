package tcpserver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_addConnection(t *testing.T) {
	t.Run("add and delete a connection", func(t *testing.T) {
		srv := New()
		reqs := make(chan Req)
		go srv.ReadRequests(reqs)
		conn5 := mockNetConn{remoteAddr: Addr{AddrString: "0005"}}
		reqs <- Req{C: Conn{Conn: &conn5}, Text: "HELLO 5"}
		time.Sleep(time.Millisecond * 30)
		assert.Len(t, srv.conns, 1)
		reqs <- Req{C: Conn{Conn: &conn5}, Text: "STOP"}
		time.Sleep(time.Millisecond * 30)
		assert.Len(t, srv.conns, 0)
	})
	t.Run("add many connections", func(t *testing.T) {
		srv := New()
		reqs := make(chan Req)
		go srv.ReadRequests(reqs)
		conn5 := mockNetConn{remoteAddr: Addr{AddrString: "0005"}}
		conn6 := mockNetConn{remoteAddr: Addr{AddrString: "0006"}}
		conn7 := mockNetConn{remoteAddr: Addr{AddrString: "0007"}}
		reqs <- Req{C: Conn{Conn: &conn5}, Text: "HELLO 5"}
		reqs <- Req{C: Conn{Conn: &conn6}, Text: "HELLO 6"}
		reqs <- Req{C: Conn{Conn: &conn7}, Text: "HELLO 7"}
		time.Sleep(time.Millisecond * 30)
		assert.Len(t, srv.conns, 3)
		reqs <- Req{C: Conn{Conn: &conn5}, Text: "STOP"}
		time.Sleep(time.Millisecond * 30)
		assert.Len(t, srv.conns, 2)
	})
}

func Test_broadcast(t *testing.T) {
	srv := New()
	reqs := make(chan Req)
	go srv.ReadRequests(reqs)
	conn5 := mockNetConn{remoteAddr: Addr{AddrString: "0005"}}
	conn6 := mockNetConn{remoteAddr: Addr{AddrString: "0006"}}
	conn7 := mockNetConn{remoteAddr: Addr{AddrString: "0007"}}
	reqs <- Req{C: Conn{Conn: &conn5}, Text: "HELLO 5"}
	reqs <- Req{C: Conn{Conn: &conn6}, Text: "HELLO 6"}
	reqs <- Req{C: Conn{Conn: &conn7}, Text: "HELLO 7"}
	reqs <- Req{C: Conn{Conn: &conn5}, Text: "HI GUYS"}

	time.Sleep(time.Millisecond * 30)
	assert.Equal(t, "HI GUYS", conn6.buf.String())
	assert.Equal(t, "HI GUYS", conn7.buf.String())
	assert.Equal(t, "", conn5.buf.String())
}
