package pooled

import (
	"bufio"
	"io"
	"net"
	"time"
)

// Conn .
type Conn interface {
	io.ReadWriter

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

func (c *conn) Write(p []byte) (int, error) {
	n, err := c.rw.Write(p)
	if err != nil {
		return n, err
	}

	return n, c.rw.Flush()
}

func (c *conn) Read(p []byte) (int, error) {
	return c.rw.Read(p)
}

func (c *conn) netconn() net.Conn {
	return c.conn
}
