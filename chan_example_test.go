package sama_test

import (
	"fmt"

	"github.com/rvflash/sama"
)

func ExampleChan() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	var sum int
	for v := range sama.Chan[int, int](ch, func(v int) int {
		return v * 3
	}) {
		sum += v
	}
	fmt.Println(sum)
	// Output: 18
}
