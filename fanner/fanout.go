package fanner

// @todo new version is untested

import "sync"

// FanOut provides one-to-many pipeline
type FanOut interface {
	// creates new listener for provided topic with optional buffering
	Register(topic string, buffer ...int) <-chan interface{}
	// removes listener from the queue and closes it
	Unregister(pipe <-chan interface{})
	// blocks until all listeners have received all messages.
	// can be called concurrently by multiple senders
	// but Register() and Unregister() will be blocked until its done.
	Send(topic string, msg ...interface{})
}

type out struct {
	mx     sync.RWMutex
	topics map[string]*topic
}

func NewOut() FanOut {
	return &out{topics: make(map[string]*topic)}
}

type topic struct {
	mx        sync.RWMutex
	listeners map[chan interface{}]chan struct{}
}

func (t *topic) add(buffer ...int) <-chan interface{} {
	t.mx.Lock()
	buff := 0
	if len(buffer) > 0 && buffer[0] > 0 {
		buff = buffer[0]
	}
	pipe := make(chan interface{}, buff)
	t.listeners[pipe] = make(chan struct{})
	t.mx.Unlock()
	return pipe
}

func (t *topic) remove(ch <-chan interface{}) {
	t.mx.Lock()
	for pipe, closer := range t.listeners {
		if pipe == ch {
			close(closer)
			delete(t.listeners, pipe)
			close(pipe)
			break
		}
	}
	t.mx.Unlock()
}

func (t *topic) send(v ...interface{}) {
	t.mx.RLock()
	wg := new(sync.WaitGroup)

	for pipe, closer := range t.listeners {
		wg.Add(1)
		go func(recv chan interface{}, closer chan struct{}, values []interface{}, wg *sync.WaitGroup) {
			for k := range values {
				select {
				case <-closer:
				case recv <- values[k]:
				}
			}
			wg.Done()
		}(pipe, closer, v, wg)
	}

	wg.Wait()
	t.mx.RUnlock()
}

func (n *out) Register(t string, buffer ...int) <-chan interface{} {
	n.mx.Lock()
	if _, ok := n.topics[t]; ok == false {
		n.topics[t] = &topic{listeners: make(map[chan interface{}]chan struct{})}
	}
	fo := n.topics[t]
	pipe := fo.add(buffer...)
	n.mx.Unlock()
	return pipe
}

func (n *out) Unregister(pipe <-chan interface{}) {
	n.mx.Lock()

	for topic, fo := range n.topics {
		fo.remove(pipe)
		if len(n.topics[topic].listeners) == 0 {
			delete(n.topics, topic)
		}
	}

	n.mx.Unlock()
}

func (n *out) Send(topic string, messages ...interface{}) {
	if len(messages) == 0 {
		return
	}

	n.mx.RLock()
	if f, ok := n.topics[topic]; ok {
		f.send(messages...)
	}
	n.mx.RUnlock()
}
