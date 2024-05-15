// Package once provides a way to call a function exactly once (memoisation).
package once

import (
	"context"
	"sync"
)

type Handle[T any] struct {
	once  sync.Once
	f     func(context.Context) (T, error)
	value T
	err   error
}

func (h *Handle[T]) Get(ctx context.Context) (T, error) {
	h.once.Do(func() { h.value, h.err = h.f(ctx) })
	return h.value, h.err
}

// Once returns a new Handle[T] that calls f exactly once.
//
// The returned Handle[T] is safe for concurrent use.
// If f returns an error, the error will be returned by Get on all subsequent calls.
func Once[T any](f func(context.Context) (T, error)) *Handle[T] {
	return &Handle[T]{f: f}
}
