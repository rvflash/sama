package sama_test

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/rvflash/sama"
)

func ExampleTsu() {
	buf := new(bytes.Buffer)
	for v := range sama.Tsu(3, func(v int) string {
		if v == 0 {
			time.Sleep(100 * time.Millisecond)
		}
		return strconv.Itoa(v)
	}, 2) {
		_, _ = buf.WriteString(v)
	}
	fmt.Println(buf.String())
	// Output: 012
}
