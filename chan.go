package sama

import (
	"sync"
)

// Chan consumes ch with a pool and returns an output channel that yields do(v) results.
// Results are emitted as soon as they're ready (order is NOT guaranteed).
func Chan[Ti, To any](ch <-chan Ti, do func(v Ti) To, concurrency ...int) chan To {
	var (
		n  = limit(concurrency)
		g  = new(sync.WaitGroup)
		rs = make(chan To, n)
	)
	for i := 0; i < n; i++ {
		g.Go(func() {
			for v := range ch {
				rs <- do(v)
			}
		})
	}
	go func() {
		g.Wait()
		close(rs)
	}()

	return rs
}
