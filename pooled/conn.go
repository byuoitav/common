package pooled

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/byuoitav/common/log"
	"go.uber.org/zap"
)

// Conn .
type Conn interface {
	io.ReadWriter

	Log() *zap.SugaredLogger
	ReadWriter() *bufio.ReadWriter
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	EmptyReadBuffer(timeout time.Duration) ([]byte, error)
	ReadUntil(delim byte, timeout time.Duration) ([]byte, error)
	netconn() net.Conn
}

type conn struct {
	log  *zap.SugaredLogger
	rw   *bufio.ReadWriter
	conn net.Conn
}

// Wrap .
func Wrap(c net.Conn) Conn {
	return &conn{
		log:  log.L.Named(c.RemoteAddr().String()),
		rw:   bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c)),
		conn: c,
	}
}

func (c *conn) Log() *zap.SugaredLogger {
	return c.log
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

func (c *conn) EmptyReadBuffer(timeout time.Duration) ([]byte, error) {
	c.conn.SetReadDeadline(time.Now().Add(timeout))

	total := c.rw.Reader.Buffered()
	bytes := make([]byte, 0, total)

	for len(bytes) < total {
		buf := make([]byte, total-len(bytes))
		_, err := c.rw.Read(buf)
		if err != nil {
			return bytes, fmt.Errorf("unable to empty read buffer: %s", err)
		}

		bytes = append(bytes, buf...)
	}

	return bytes, nil
}

func (c *conn) ReadUntil(delim byte, timeout time.Duration) ([]byte, error) {
	c.SetReadDeadline(time.Now().Add(timeout))
	return c.rw.ReadBytes(delim)
}

func (c *conn) netconn() net.Conn {
	return c.conn
}
