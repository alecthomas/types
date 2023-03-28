package types

import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestOptionGet(t *testing.T) {
	o := Some(1)
	v, ok := o.Get()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	o = None[int]()
	_, ok = o.Get()
	assert.False(t, ok)
}

func TestOptionMarshalJSON(t *testing.T) {
	o := Some(1)
	b, err := o.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "1", string(b))

	o = None[int]()
	b, err = o.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "null", string(b))
}

func TestOptionUnmarshalJSON(t *testing.T) {
	o := Option[int]{}
	err := json.Unmarshal([]byte("1"), &o)
	assert.NoError(t, err)
	b, ok := o.Get()
	assert.True(t, ok)
	assert.Equal(t, 1, b)
}

func TestOptionString(t *testing.T) {
	o := Some(1)
	assert.Equal(t, "1", o.String())

	o = None[int]()
	assert.Equal(t, "None", o.String())
}

func TestOptionGoString(t *testing.T) {
	o := Some(1)
	assert.Equal(t, "Some[int](1)", o.GoString())

	o = None[int]()
	assert.Equal(t, "None[int]()", o.GoString())
}
