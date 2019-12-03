package rqm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fwojciec/rqm"
)

func TestQueue(t *testing.T) {
	ctx := context.Background()
	q := &rqm.Queue{}
	_, err := q.Pop(ctx)
	if !errors.Is(err, rqm.ErrQueueEmpty) {
		t.Fatalf("Expected ErrQueueEmpty error, received: %s", err)
	}
	err = q.Push(ctx, &rqm.Rq{URL: "test 1"})
	if err != nil {
		t.Error(err)
	}
	err = q.Push(ctx, &rqm.Rq{URL: "test 2"})
	if err != nil {
		t.Error(err)
	}
	x1, err := q.Pop(ctx)
	if err != nil {
		t.Error(err)
	}
	if x1.URL != "test 1" {
		t.Fatalf(`expected the URL to be "test 1", received %q`, x1.URL)
	}
	x2, err := q.Pop(ctx)
	if err != nil {
		t.Error(err)
	}
	if x2.URL != "test 2" {
		t.Fatalf(`expected the URL to be "test 2", received %q`, x2.URL)
	}
	_, err = q.Pop(ctx)
	if !errors.Is(err, rqm.ErrQueueEmpty) {
		t.Fatalf("Expected ErrQueueEmpty error, received: %s", err)
	}
}
