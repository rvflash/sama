package sama_test

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/rvflash/sama"
)

func ExampleSan() {
	ch := make(chan int, 3)
	ch <- 3
	ch <- 2
	ch <- 1
	close(ch)

	buf := new(bytes.Buffer)
	for v := range sama.San[int, string](ch, func(v int) string {
		time.Sleep(time.Duration(v) * 100 * time.Millisecond)
		return strconv.Itoa(v)
	}) {
		_, _ = buf.WriteString(v)
	}
	fmt.Println(buf.String())
	// Output: 321
}
