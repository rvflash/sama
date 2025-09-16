package sama_test

import (
	"fmt"
	"sync/atomic"

	"github.com/rvflash/sama"
)

func ExampleKun() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	var sum atomic.Int32
	sama.Kun[int](ch, func(v int) {
		sum.Add(int32(v))
	})
	fmt.Println(sum.Load())
	// Output: 6
}
