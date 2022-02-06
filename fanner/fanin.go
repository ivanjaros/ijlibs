package fanner

import (
	"sync"
)

// FanIn provides many-to-one pipeline.
type FanIn interface {
	Add(ch <-chan interface{})
	Remove(ch <-chan interface{})
	Listen() <-chan interface{}
	Close()
}

type in struct {
	out       chan interface{}
	listeners map[<-chan interface{}]chan struct{}
	mx        *sync.RWMutex
	wg        *sync.WaitGroup
}

func NewIn(size int) FanIn {
	if size < 1 {
		size = 1
	}
	return &in{
		listeners: make(map[<-chan interface{}]chan struct{}),
		out:       make(chan interface{}, size),
		mx:        new(sync.RWMutex),
		wg:        new(sync.WaitGroup),
	}
}

func (f *in) Add(ch <-chan interface{}) {
	f.mx.Lock()
	closer := make(chan struct{})
	f.listeners[ch] = closer

	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		for {
			select {
			case v, ok := <-ch:
				if ok == false {
					return
				}

				select {
				case f.out <- v:
				default:
				}

			case <-closer:
				return
			}
		}
	}()

	f.mx.Unlock()
}

func (f *in) Remove(ch <-chan interface{}) {
	f.mx.Lock()
	if closer, ok := f.listeners[ch]; ok {
		close(closer)
		delete(f.listeners, ch)
	}
	f.mx.Unlock()
}

func (f *in) Listen() <-chan interface{} {
	return f.out
}

func (f *in) Close() {
	f.mx.Lock()
	for k := range f.listeners {
		close(f.listeners[k])
		delete(f.listeners, k)
	}
	f.wg.Wait()
	f.mx.Unlock()
}
