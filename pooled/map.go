package pooled

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/byuoitav/common/log"
)

// NewConnection .
type NewConnection func(key interface{}) (Conn, error)

// Work .
type Work func(Conn) error

// Map .
type Map struct {
	newConn NewConnection
	ttl     time.Duration

	m  map[interface{}]chan request
	mu sync.Mutex
}

type request struct {
	work Work
	resp chan error
}

// NewMap .
func NewMap(ttl time.Duration, newConn NewConnection) *Map {
	return &Map{
		m:       make(map[interface{}]chan request),
		newConn: newConn,
		ttl:     ttl,
	}
}

// Do .
func (m *Map) Do(key interface{}, work Work) error {
	var reqs chan request
	var ok bool

	l := log.L.Named(fmt.Sprintf("%s", key))

	m.mu.Lock()
	if reqs, ok = m.m[key]; !ok {
		// open a new connection
		l.Infof("Opening new connection")
		conn, err := m.newConn(key)
		if err != nil {
			m.mu.Unlock()
			return fmt.Errorf("failed to open new connection for %s: %s", key, err)
		}

		if conn == nil {
			m.mu.Unlock()
			return fmt.Errorf("got nil connection from new connection function")
		}

		reqs = make(chan request, 10)
		m.m[key] = reqs
		m.mu.Unlock()

		conn.Log().Infof("Successfully opened new connection")

		go func() {
			defer func() {
				m.mu.Lock()
				delete(m.m, key)
				m.mu.Unlock()

				conn.Log().Infof("Closing connection")
				close(reqs)

				// finish up remaining requests
				for req := range reqs {
					req.resp <- req.work(conn)
				}

				conn.netconn().Close()
			}()

			timer := time.NewTimer(m.ttl)

			for {
				// reset the buffer by reading everything currently in it
				bytes, err := conn.EmptyReadBuffer(m.ttl)
				if err != nil {
					conn.Log().Warnf("failed to empty buffer: %s", err)
					return
				}
				if len(bytes) > 0 {
					conn.Log().Debugf("Read %v leftover bytes: 0x%x", len(bytes), bytes)
				}

				// reset the deadlines
				conn.netconn().SetDeadline(time.Time{})

				select {
				case req := <-reqs:
					err := req.work(conn)
					req.resp <- err
					if err, ok := err.(net.Error); ok && (!err.Temporary() || err.Timeout()) {
						// if it was a timeout error, close the connection
						conn.Log().Warnf("closing connection due to non-temporary or timeout error: %s", err.Error())
						return
					}

					// reset the timer
					if !timer.Stop() {
						<-timer.C
					}
					timer.Reset(m.ttl)
				case <-timer.C:
					return
				}
			}
		}()
	} else {
		l.Infof("Reusing already open connection")
		m.mu.Unlock()
	}

	req := request{
		work: work,
		resp: make(chan error),
	}

	reqs <- req
	return <-req.resp
}
