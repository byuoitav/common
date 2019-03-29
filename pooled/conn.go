package pooled

import (
	"bufio"
	"net"
	"time"
)

// Conn .
type Conn interface {
	ReadWriter() *bufio.ReadWriter
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error

	netconn() net.Conn
}

type conn struct {
	rw   *bufio.ReadWriter
	conn net.Conn
}

// Wrap .
func Wrap(c net.Conn) Conn {
	return &conn{
		rw:   bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c)),
		conn: c,
	}
}

func (c *conn) ReadWriter() *bufio.ReadWriter {
	return c.rw
}

func (c *conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *conn) netconn() net.Conn {
	return c.conn
}
