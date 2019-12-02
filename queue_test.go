package rqm_test

import (
	"context"
	"testing"

	"github.com/fwojciec/rqm"
)

func TestQueue(t *testing.T) {
	ctx := context.Background()
	q := &rqm.Queue{}
	isEmpty1, err := q.IsEmpty(ctx)
	if err != nil {
		t.Error(err)
	}
	if !isEmpty1 {
		t.Fatalf("expected a freshly initialized queue to be empty")
	}
	err = q.Push(ctx, &rqm.Rq{URL: "test 1"})
	if err != nil {
		t.Error(err)
	}
	err = q.Push(ctx, &rqm.Rq{URL: "test 2"})
	if err != nil {
		t.Error(err)
	}
	isEmpty2, err := q.IsEmpty(ctx)
	if err != nil {
		t.Error(err)
	}
	if isEmpty2 {
		t.Fatalf("expected a queue with two items in it to be not empty")
	}
	x1, err := q.Pop(ctx)
	if err != nil {
		t.Error(err)
	}
	if x1.URL != "test 1" {
		t.Fatalf(`expected the URL to be "test 1", received %q`, x1.URL)
	}
	isEmpty3, err := q.IsEmpty(ctx)
	if err != nil {
		t.Error(err)
	}
	if isEmpty3 {
		t.Fatalf("expected a queue with one item in it to be not empty")
	}
	x2, err := q.Pop(ctx)
	if err != nil {
		t.Error(err)
	}
	if x2.URL != "test 2" {
		t.Fatalf(`expected the URL to be "test 2", received %q`, x2.URL)
	}
	isEmpty4, err := q.IsEmpty(ctx)
	if err != nil {
		t.Error(err)
	}
	if !isEmpty4 {
		t.Fatalf("expected a queue with no items in it to be empty")
	}
}
