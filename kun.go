package sama

import "sync"

// Kun consumes ch with a pool and calls do(v) for each item.
// It does not produce any output.
func Kun[T any](ch chan T, do func(v T), concurrency ...int) {
	var (
		g = new(sync.WaitGroup)
		n = limit(concurrency)
	)
	for i := 0; i < n; i++ {
		g.Go(func() {
			for v := range ch {
				do(v)
			}
		})
	}
	g.Wait()
}
