package rqm

import (
	"math/rand"
	"time"
)

// RandomTicker is similar to time.Ticker but ticks at random intervals. The
// min and max values are the duration in milliseconds of the shortest and
// longest tick respectively.
type RandomTicker struct {
	C     chan time.Time
	stopc chan struct{}
	min   int
	max   int
}

// NewRandomTicker returns a pointer to an initialized instance of the
// RandomTicker. Min and max are durations of the shortest and longest allowed
// tick, in milliseconds.
func NewRandomTicker(min, max int) *RandomTicker {
	rt := &RandomTicker{
		C:     make(chan time.Time),
		stopc: make(chan struct{}),
		min:   min,
		max:   max,
	}
	go rt.loop()
	return rt
}

// Stop stops the instance of RandomTicker.
func (rt *RandomTicker) Stop() {
	rt.stopc <- struct{}{}
}

func (rt *RandomTicker) loop() {
	for {
		time.Sleep(rt.nextInterval())
		select {
		case <-rt.stopc:
			return
		case rt.C <- time.Now():
		default:
		}
	}
}

func (rt *RandomTicker) nextInterval() time.Duration {
	return time.Duration(rand.Intn(rt.max-rt.min)+rt.min) * time.Millisecond
}
