package client

import (
	"context"
	"sync"

	"github.com/nats-io/nats.go"
)

type conn struct {
	mux      sync.Mutex
	conn     *nats.EncodedConn
	url      string
	encType  string
	name     string
	maxRecon int
}

// Close - closes connection
func (c *conn) Close() {
	c.mux.Lock()
	if c.conn == nil {
		c.mux.Unlock()
		return
	}

	c.conn.Close()
	c.conn = nil
	c.mux.Unlock()
}

func (c *conn) publish(sub string, v interface{}) (err error) {
	if err = c.connect(); err != nil {
		return err
	}

	return c.conn.Publish(sub, v)
}

func (c *conn) subscribe(sub string, cb nats.Handler) (s *nats.Subscription, err error) {
	if err = c.connect(); err != nil {
		return nil, err
	}

	return c.conn.Subscribe(sub, cb)
}

func (c *conn) request(ctx context.Context, sub string, v interface{}, vPtr interface{}) (err error) {
	if err = c.connect(); err != nil {
		return err
	}

	return c.conn.RequestWithContext(ctx, sub, v, vPtr)
}

func (c *conn) connect() error {
	c.mux.Lock()
	if c.conn != nil {
		c.mux.Unlock()
		return nil
	}

	sc, err := nats.Connect(c.url, nats.Name(c.name), nats.MaxReconnects(c.maxRecon))
	if err != nil {
		c.mux.Unlock()
		return err
	}

	c.conn, err = nats.NewEncodedConn(sc, c.encType)
	c.mux.Unlock()

	return err
}

// newConn - creates connector for auto connecting
func newConn(url, encType, name string, maxReconnects int) *conn {
	return &conn{
		url:      url,
		encType:  encType,
		name:     name,
		maxRecon: maxReconnects,
		conn:     nil,
	}
}
