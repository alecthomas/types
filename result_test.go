package types

import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestResultGet(t *testing.T) {
	r := Ok(1)
	v, ok := r.Get()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	r = Errorf[int]("foo")
	_, ok = r.Get()
	assert.False(t, ok)
}

func TestResultMarshalJSON(t *testing.T) {
	r := Ok(1)
	b, err := json.Marshal(r)
	assert.NoError(t, err)
	assert.Equal(t, `{"value":1}`, string(b))

	r = Errorf[int]("foo")
	b, err = json.Marshal(r)
	assert.NoError(t, err)
	assert.Equal(t, `{"error":"foo"}`, string(b))
}

func TestResultUnmarshalJSON(t *testing.T) {
	r := Result[int]{}
	err := json.Unmarshal([]byte(`{"value":1}`), &r)
	assert.NoError(t, err)
	assert.Equal(t, Ok(1), r)

	r = Errorf[int]("error")
	err = json.Unmarshal([]byte(`{"error":"error"}`), &r)
	assert.NoError(t, err)
}
