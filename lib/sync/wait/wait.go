package wait

import (
	"sync"
	"time"
)

// implement the wait group with timeout

type Wait struct {
	w sync.WaitGroup
}

func (wt *Wait) Add(delta int) {
	wt.w.Add(delta)
}

func (wt *Wait) Done() {
	wt.w.Done()
}

func (wt *Wait) Wait() {
	wt.w.Wait()
}

func (wt *Wait) WaitWithTimeout(timeout time.Duration) bool {
	c := make(chan struct{}, 1)
	go func() {
		wt.w.Wait()
		c <- struct{}{}
	}()

	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}
