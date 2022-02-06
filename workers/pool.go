// container/list is used instead of a map for preserving order of workers/listeners.
// even though order is not actually required it is better approach overall.

package workers

import (
	"container/list"
	"sync"
)

type Pool interface {
	Take(n int) bool
	Return()
	Available() int
	Wait(n uint) bool
	Done() bool
	Max() int
}

type workerPool struct {
	mx        sync.RWMutex
	max       int
	avail     int
	listeners *list.List
	lMx       sync.Mutex
}

func NewPool(maxAvail int) *workerPool {
	return &workerPool{max: maxAvail, avail: maxAvail, listeners: list.New()}
}

type listener struct {
	above int
	ch    chan bool
}

func (w *workerPool) Take(n int) bool {
	w.mx.Lock()
	ok := w.avail >= n
	if ok {
		w.avail -= n
		defer w.ping()
	}
	w.mx.Unlock()
	return ok
}

func (w *workerPool) Return() {
	w.mx.Lock()
	if w.avail < w.max {
		w.avail++
		defer w.ping()
	}
	w.mx.Unlock()
}

func (w *workerPool) Available() int {
	w.mx.RLock()
	av := w.avail
	w.mx.RUnlock()
	return av
}

// this always returns true, but it waits for n+1 available workers before doing so
func (w *workerPool) Wait(n uint) bool {
	return w.listen(int(n))
}

// same as Wait() but it will return only when all workers are available(all processing is done)
func (w *workerPool) Done() bool {
	return w.listen(-1)
}

func (w *workerPool) listen(n int) bool {
	w.lMx.Lock()

	worker := &listener{
		above: n,
		ch:    make(chan bool),
	}
	e := w.listeners.PushBack(worker)

	w.lMx.Unlock()

	go w.ping()

	<-worker.ch

	w.lMx.Lock()
	close(worker.ch)
	w.listeners.Remove(e)
	w.lMx.Unlock()

	return true
}

func (w *workerPool) Max() int {
	return w.max
}

func (w *workerPool) ping() {
	w.lMx.Lock()
	avail := w.Available()
	done := avail == w.Max()

	for e := w.listeners.Front(); e != nil; e = e.Next() {
		worker := e.Value.(*listener)
		if (worker.above == -1 && done) || (worker.above > -1 && worker.above < avail) {
			worker.ch <- true
		}
	}

	w.lMx.Unlock()
}
