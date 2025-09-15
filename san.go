package sama

import (
	"sync"
)

func San[Ti, To any](ch <-chan Ti, do func(v Ti) To, concurrency ...int) chan To {
	type (
		req struct {
			k int
			v Ti
		}
		res struct {
			k int
			v To
		}
	)
	var (
		i, next int
		n       = limit(concurrency)
		g       = new(sync.WaitGroup)
		rq      = make(chan req, n)
		rs      = make(chan res, n)
		out     = make(chan To, n)
		buf     = make(map[int]To)
	)
	go func() {
		defer close(rq)
		for v := range ch {
			rq <- req{
				k: i,
				v: v,
			}
			i++
		}
	}()
	for w := 0; w < n; w++ {
		g.Go(func() {
			for v := range rq {
				rs <- res{
					k: v.k,
					v: do(v.v),
				}
			}
		})
	}
	go func() {
		g.Wait()
		close(rs)
	}()
	go func() {
		defer close(out)
		var (
			bv To
			ok bool
		)
		for v := range rs {
			if v.k == next {
				out <- v.v
				next++
				for {
					bv, ok = buf[next]
					if !ok {
						break
					}
					out <- bv
					delete(buf, next)
					next++
				}
			} else {
				buf[v.k] = v.v
			}
		}
	}()
	return out
}
