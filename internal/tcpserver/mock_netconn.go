package tcpserver

import (
	"bytes"
	"net"
	"time"
)

type mockNetConn struct {
	remoteAddr net.Addr
	buf        bytes.Buffer
}

// Addr is a fake network interface which implements the net.Addr interface
type Addr struct {
	NetworkString string
	AddrString    string
}

func (a Addr) Network() string {
	return a.NetworkString
}

func (a Addr) String() string {
	return a.AddrString
}

func (c *mockNetConn) Read(b []byte) (n int, err error) {
	panic("implement me")
}

func (c *mockNetConn) Write(b []byte) (n int, err error) {
	c.buf.Write(b)
	return 0, err
}

func (c *mockNetConn) Close() error {
	panic("implement me")
}

func (c *mockNetConn) LocalAddr() net.Addr {
	panic("implement me")
}

func (c *mockNetConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *mockNetConn) SetDeadline(t time.Time) error {
	panic("implement me")
}

func (c *mockNetConn) SetReadDeadline(t time.Time) error {
	panic("implement me")
}

func (c *mockNetConn) SetWriteDeadline(t time.Time) error {
	panic("implement me")
}
