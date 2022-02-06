//  This is altered copy from github.com/dgraph-io/badger/y/y.go

package workers

import "sync"

type Throttle interface {
	Do() error
	Done() error
	Finish() error
}

type throttled struct {
	wg    sync.WaitGroup
	ch    chan struct{}
	errCh chan error
}

func NewThrottle(max int) *throttled {
	return &throttled{
		ch:    make(chan struct{}, max),
		errCh: make(chan error, max),
	}
}

func (t *throttled) Do() error {
	select {
	case t.ch <- struct{}{}:
		t.wg.Add(1)
		return nil
	case err := <-t.errCh:
		return err
	}
}

func (t *throttled) Done(err error) {
	if err != nil {
		t.errCh <- err
	}
	select {
	case <-t.ch:
	default:
		panic("Do/Done flow mismatch")
	}
	t.wg.Done()
}

func (t *throttled) Finish() error {
	t.wg.Wait()
	close(t.ch)
	close(t.errCh)
	for err := range t.errCh {
		if err != nil {
			return err
		}
	}
	return nil
}
