package sama_test

import (
	"testing"

	"github.com/rvflash/sama"
)

func TestTsu(t *testing.T) {
	t.Parallel()

	t.Run("Default", func(t *testing.T) {
		var n int
		for i := range sama.Tsu(0, func(v int) int {
			return v
		}) {
			n += i
		}
		if n != 0 {
			t.Errorf("got %d, want 0", n)
		}
	})

	t.Run("OK", func(t *testing.T) {
		const size = 3
		var n int
		for i := range sama.Tsu(size, func(v int) int {
			return v
		}) {
			n += i
		}
		if n != size {
			t.Errorf("got %d, want %d", n, size)
		}
	})
}
