package sama

import "runtime"

func limit(opts []int) int {
	var n int
	if len(opts) == 1 {
		n = opts[0]
	} else {
		n = runtime.NumCPU() * 2
	}
	return n
}
