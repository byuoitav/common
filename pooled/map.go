package pooled

import (
	"fmt"
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

	m.mu.Lock()
	if reqs, ok = m.m[key]; !ok {
		// open a new connection
		l := log.L.Named(fmt.Sprintf("%s", key))
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

		l.Infof("Successfully opened new connection")

		go func() {
			defer func() {
				l.Infof("Closing connection")
				close(reqs)

				// finish up remaining requests
				for req := range reqs {
					req.resp <- req.work(conn)
				}

				conn.netconn().Close()
			}()

			timer := time.NewTimer(m.ttl)
			for {
				select {
				case req := <-reqs:
					req.resp <- req.work(conn)

					// reset the timer
					if !timer.Stop() {
						<-timer.C
					}
					timer.Reset(m.ttl)
				case <-timer.C:
					m.mu.Lock()
					delete(m.m, key)
					m.mu.Unlock()
					return
				}
			}
		}()
	} else {
		log.L.Infof("Reusing already open connection for %s", key)
		m.mu.Unlock()
	}

	req := request{
		work: work,
		resp: make(chan error),
	}

	reqs <- req
	return <-req.resp
}
