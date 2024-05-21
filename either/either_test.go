package either

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestEither(t *testing.T) {
	left := LeftOf[string](42)
	assert.Equal(t, 42, left.Get())
	right := RightOf[int]("foo")
	assert.Equal(t, "foo", right.Get())

	var either Either[int, string] = left
	if either, ok := either.(Left[int, string]); ok {
		assert.Equal(t, 42, either.Get())
	}
}
