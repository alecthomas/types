package once

import (
	"context"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/alecthomas/types/must"
)

func TestOnce(t *testing.T) {
	ctx := context.Background()
	calls := 0
	h := Once(func(context.Context) (int, error) { calls++; return 42, nil })
	assert.Equal(t, 42, must.Get(h.Get(ctx)))
	assert.Equal(t, 42, must.Get(h.Get(ctx)))
	assert.Equal(t, 42, must.Get(h.Get(ctx)))
	assert.Equal(t, 1, calls)
}
