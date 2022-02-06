package conder

import (
	"sync"
)

func New() *Cond {
	mx := &sync.Mutex{}
	return &Cond{c: sync.NewCond(mx), mx: mx}
}

type Cond struct {
	mx *sync.Mutex
	c  *sync.Cond
}

// invoke one waiter
func (c *Cond) Next() {
	c.c.Signal()
}

// invoke all waiters
func (c *Cond) Done() {
	c.c.Broadcast()
}

// do work once Next() or Done() is called
func (c *Cond) Do(work func()) {
	c.mx.Lock()
	c.c.Wait()
	work()
	c.mx.Unlock()
}

// waits for Next() or Done() to be called before returning
func (c *Cond) Wait() {
	ch := make(chan struct{}, 1)
	c.Do(func() { ch <- struct{}{} })
	<-ch
	return
}
