package workers

import (
	"context"
	"golang.org/x/sync/semaphore"
)

type Semaphore interface {
	Take(n uint) bool
	Try(n uint) bool
	Return(n uint)
}

type weighted struct {
	sem *semaphore.Weighted
}

func NewSemaphore(n uint) Semaphore {
	return &weighted{sem: semaphore.NewWeighted(int64(n))}
}

func (s *weighted) Take(n uint) bool {
	return s.sem.Acquire(context.Background(), int64(n)) == nil
}

func (s *weighted) Try(n uint) bool {
	return s.sem.TryAcquire(int64(n))
}

func (s *weighted) Return(n uint) {
	s.sem.Release(int64(n))
}
