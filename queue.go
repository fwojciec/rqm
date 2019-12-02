package rqm

import (
	"context"
	"sync"
)

// Queue is a basic queuer implementation.
type Queue struct {
	q  []*Rq
	mu sync.RWMutex
}

// Pop removes an element from the front of the queue and returns it.
func (q *Queue) Pop(ctx context.Context) (r *Rq, err error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	r, q.q = q.q[0], q.q[1:]
	return
}

// Push adds an element to the front of the queue.
func (q *Queue) Push(ctx context.Context, r *Rq) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.q = append(q.q, r)
	return nil
}

// IsEmpty returns true if the queue is empty and false if not.
func (q *Queue) IsEmpty(ctx context.Context) (bool, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.q) == 0, nil
}
