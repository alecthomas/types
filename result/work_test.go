package result

import (
	"errors"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
)

func TestGo(t *testing.T) {
	out := Go(func() (int, error) {
		return 1, nil
	})
	assertResult(t, Ok(1), out)

	out = Go(func() (int, error) {
		return 0, errors.New("error")
	})
	assertResult(t, Errorf[int]("error"), out)
}

func assertResult[T any](t *testing.T, expected Result[T], actual chan Result[T]) {
	t.Helper()
	select {
	case actual := <-actual:
		assert.Equal(t, expected, actual)

	case <-time.After(time.Millisecond * 100):
		t.Fatal("timeout")
	}
}
