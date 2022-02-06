package listener_closer

import (
	"net"
	"sync"
)

func Wrap(l net.Listener) net.Listener {
	return &closer{Listener: l}
}

type closer struct {
	net.Listener
	once sync.Once
	e    error
}

func (c *closer) Close() error {
	c.once.Do(func() { c.e = c.Listener.Close() })
	return c.e
}
