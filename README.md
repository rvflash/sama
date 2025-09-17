# sama

[![GoDoc](https://godoc.org/github.com/rvflash/sama?status.svg)](https://godoc.org/github.com/rvflash/sama)
[![Build Status](https://github.com/rvflash/sama/workflows/build/badge.svg)](https://github.com/rvflash/sama/actions?workflow=build)
[![Code Coverage](https://codecov.io/gh/rvflash/sama/branch/main/graph/badge.svg)](https://codecov.io/gh/rvflash/sama)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/sama?)](https://goreportcard.com/report/github.com/rvflash/sama)


Tiny, fast, generic worker pools for Go.

`sama` exposes three functions with the same spirit and different guarantees:
* `Kun`, fire-and-forget: process items from a channel with a worker pool, no output.
* `Chan`, unordered: process items from a channel with a worker pool, returns results as they are ready.
* `San`, ordered: same as `Chan`, but preserves input order in the output stream.

All three stop naturally when the input channel is closed and all in-flight work finishes.

Go ≥ 1.25 (no external dependency, only uses generics and the new `Go` function from `sync.WaitGroup`).

```go
// Kun consumes ch with a pool and calls do(v) for each item.
// It does not produce any output.
func Kun[T any](ch chan T, do func(v T), concurrency ...int)

// Chan consumes ch with a pool and returns an output channel that yields do(v) results.
// Results are emitted as soon as they're ready (order is NOT guaranteed).
func Chan[Ti, To any](ch <-chan Ti, do func(v Ti) To, concurrency ...int) chan To

// San is like Chan but guarantees that the output preserves the input order.
// The i-th value read from ch produces the i-th value on the returned channel.
func San[Ti, To any](ch <-chan Ti, do func(v Ti) To, concurrency ...int) chan To
```

> `concurrency` is optional. If omitted, `sama` uses an arbitrary default (2x`runtime.NumCPU()`).
> 
> Close the input channel to finish; the output channel (for `Chan`/`San`) will close automatically when all work is done.

## Installation

```bash
go get github.com/rvflash/sama
```

## Patterns & tips

### Sample use-case

```go
var (
    in = make(chan string)
    // San preserves input order in the output.
    out = sama.San(in, func(s string) string {
        // Simulate variable latency to show ordering guarantee.
        if s == "bravo" {
            time.Sleep(50 * time.Millisecond)
        }
        return strings.ToUpper(s)
    })
)

go func() {
    for _, v := range []string{"alpha", "bravo", "charlie", "delta"} {
        in <- v
    }
    close(in)
}()

// Order matches the input exactly.
for res := range out {
    fmt.Println(res)
}
// Output: 
// ALPHA
// BRAVO
// CHARLIE
// DELTA
```

### Backpressure & buffering

The input channel acts as backpressure. If producers outpace consumers, either:
* Increase `concurrency`, or
* Use a buffered `chan` for the input.

For `Chan`/`San`, the returned output channel is buffered enough for good throughput.
Still, if a downstream consumer is slow, overall speed will match the slowest stage.

### Cancellation & shutdown

Since the API does not accept `context.Context`, cancellation is cooperative:
- Close the input channel to signal completion.
- Ensure your `do` function returns promptly (check your own contexts inside do if it performs I/O).

## Errors

Use sum types or tuples to propagate errors:
```go
type Out struct {
    Val string
    Err error
}
out := sama.Chan(in, func(s string) Out {
    v, err := doWork(s)
    return Out{Val: v, Err: err}
}, 8)
```

When to choose which:
* `Kun`: side effects only (DB writes, HTTP calls where you handle errors internally).
* `Chan`: maximum throughput, order doesn’t matter (idempotent / commutative workloads).
* `San`: streaming and order matters (like re-sequencing responses for a client).


## Guarantees

* No leaks: all goroutines exit once the input channel is closed and all work completes.
* Ordering:
  - `San` stable (input index order).
  - `Chan` none (as-completed).
* Throughput:
  - `Chan` tends to be the fastest.
  - `San` adds a tiny reorder buffer proportional to “out-of-order window”.