package pooled

import (
	"net"
	"sync"
	"time"

	"github.com/byuoitav/common/log"
)

// NewConnection .
type NewConnection func(key interface{}) (net.Conn, error)

// Work .
type Work func(net.Conn) error

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
		log.L.Infof("Opening new connection for %s", key)
		conn, err := m.newConn(key)
		if err != nil {
			m.mu.Unlock()
			return err
		}

		reqs = make(chan request, 10)
		m.m[key] = reqs
		m.mu.Unlock()

		go func() {
			defer func() {
				log.L.Infof("Closing connection for %s", key)
				// finish up remaining requests
				for req := range reqs {
					req.resp <- req.work(conn)
				}

				conn.Close()
			}()

			timer := time.NewTimer(m.ttl)
			for {
				select {
				case req := <-reqs:
					// reset the timer
					if !timer.Stop() {
						<-timer.C
					}
					timer.Reset(m.ttl)

					req.resp <- req.work(conn)
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
