package workers

import "sync"

// monitor maps multiple "workers" by name and their statuses and collectively
// returns true if all workers are true or false if at least one worker is not.
type Monitor interface {
	Push(name string, status bool)
	Status() bool
	// listening on the channel is not required, pushing changes will not block
	Listen() <-chan bool
	// resets all workers to true and makes the Status() and Listen() active again,
	// Listen() will also receive true if the monitor has been stopped before.
	Reset()
	// stops the monitor so Status() will always return false and Listen() will once receive false,
	// then it will no longer receives changes, if the monitor has not been stopped before.
	Stop()
}

func NewMonitor() *monitor {
	return &monitor{s: make(map[string]bool), ch: make(chan bool)}
}

type monitor struct {
	mx      sync.Mutex
	stopped bool
	s       map[string]bool
	ch      chan bool
}

func (m *monitor) Push(name string, status bool) {
	m.mx.Lock()
	if m.stopped {
		m.mx.Unlock()
		return
	}
	before := m.status()
	m.s[name] = status
	after := m.status()
	m.mx.Unlock()

	if before != after {
		select {
		case m.ch <- after:
		default:
		}
	}
}

func (m *monitor) Status() bool {
	m.mx.Lock()
	s := m.status()
	m.mx.Unlock()
	return s
}

func (m *monitor) status() bool {
	if m.stopped {
		return false
	}
	s := true
	for k := range m.s {
		if m.s[k] == false {
			s = false
			break
		}
	}
	return s
}

func (m *monitor) Listen() <-chan bool {
	return m.ch
}

func (m *monitor) Stop() {
	m.mx.Lock()
	if m.stopped == false {
		m.stopped = true
		for k := range m.s {
			m.s[k] = false
		}
		select {
		case m.ch <- false:
		default:
		}
	}
	m.mx.Unlock()
}

func (m *monitor) Reset() {
	m.mx.Lock()
	if m.stopped {
		m.stopped = false
		for k := range m.s {
			m.s[k] = true
		}
		select {
		case m.ch <- true:
		default:
		}
	}
	m.mx.Unlock()
}
