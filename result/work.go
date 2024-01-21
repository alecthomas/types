package result

// Go runs a function in a goroutine and returns a channel that will receive the
// Result.
func Go[T any](f func() (T, error)) chan Result[T] {
	out := make(chan Result[T])
	go func() {
		defer close(out)
		out <- Outcome(f())
	}()
	return out
}
