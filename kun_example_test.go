package sama_test

import (
	"fmt"

	"github.com/rvflash/sama"
)

func ExampleKun() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	var sum int
	sama.Kun[int](ch, func(v int) {
		sum += v
	})
	fmt.Println(sum)
	// Output: 6
}
