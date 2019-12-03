package rqm

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrQueueEmpty is returned when trying to Pop an item from empty queue
	ErrQueueEmpty = errors.New("queue is empty")
)

// Queuer manages the queue of requests waiting to be processed.
type Queuer interface {
	Pop(context.Context) (*Rq, error)
	Push(context.Context, *Rq) error
}

// Queue is a basic queuer implementation.
type Queue struct {
	q  []*Rq
	mu sync.RWMutex
}

// Pop removes an element from the front of the queue and returns it.
func (q *Queue) Pop(ctx context.Context) (r *Rq, err error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.q) == 0 {
		return nil, ErrQueueEmpty
	}
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
