package result

import (
	"encoding/json"
	"errors"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestResultGet(t *testing.T) {
	r := Ok(1)
	v, ok := r.Get()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	r = Errf[int]("foo")
	_, ok = r.Get()
	assert.False(t, ok)
}

func TestResultMarshalJSON(t *testing.T) {
	r := Ok(1)
	b, err := json.Marshal(r)
	assert.NoError(t, err)
	assert.Equal(t, `{"value":1}`, string(b))

	r = Errf[int]("foo")
	b, err = json.Marshal(r)
	assert.NoError(t, err)
	assert.Equal(t, `{"error":"foo"}`, string(b))
}

func TestResultUnmarshalJSON(t *testing.T) {
	r := Result[int]{}
	err := json.Unmarshal([]byte(`{"value":1}`), &r)
	assert.NoError(t, err)
	assert.Equal(t, Ok(1), r)

	r = Errf[int]("error")
	err = json.Unmarshal([]byte(`{"error":"error"}`), &r)
	assert.NoError(t, err)
}

func TestResultMap(t *testing.T) {
	err := errors.New("error")
	tests := []struct {
		name     string
		input    Result[string]
		expected Result[int64]
		called   bool
	}{
		{"Success", Ok("1234"), Ok[int64](1234), true},
		{"Error", Err[string](err), Err[int64](err), false},
		{"Invalid", Ok("hello"), Err[int64](errors.New(`strconv.ParseInt: parsing "hello": invalid syntax`)), true},
	}
	called := false
	curry := func(v string) (int64, error) {
		called = true
		return strconv.ParseInt(v, 10, 64)
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Map(test.input, curry)
			assert.Equal(t, test.expected, actual)
			assert.Equal(t, test.called, called)
			called = false
		})
	}
}
