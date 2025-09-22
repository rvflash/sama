package sama

import "runtime"

// Tsu iterates from 0 to n (not included), calling “do” for each value and returning the results in order.
func Tsu[To any](n int, do func(v int) To, concurrency ...int) chan To {
	if n <= 0 {
		out := make(chan To, limit(concurrency))
		close(out)
		return out
	}
	in := make(chan int, n)
	for i := 0; i < n; i++ {
		in <- i
	}
	close(in)

	return San(in, do)
}

func limit(opts []int) int {
	var n int
	if len(opts) == 1 {
		n = opts[0]
	} else {
		n = runtime.NumCPU() * 2
	}
	return n
}
