package sama

import (
	"sync"
)

// Chan
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
