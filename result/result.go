// Package result provides a Result type that can contain a value or an error.
package result

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Stdlib interfaces types implement.
type stdlib interface {
	fmt.Stringer
	fmt.GoStringer
	json.Marshaler
	json.Unmarshaler
}

var _ stdlib = (*Result[int])(nil)

// Ok returns a Result that contains a value.
func Ok[T any](value T) Result[T] { return Result[T]{value: value} }

// Err returns a Result that contains an error.
func Err[T any](err error) Result[T] { return Result[T]{err: err} }

// Map a Result[L] to Result[R]. If Result[L] is an error mapper will not be called and the result will be an error.
func Map[L, R any](l Result[L], mapper func(L) (R, error)) Result[R] {
	value, err := l.Result()
	if err != nil {
		return Err[R](err)
	}
	return Outcome(mapper(value))
}

// From returns a Result that contains a value or an error.
//
// It can be used to convert a function that returns a value and an error into a
// Result.
func From[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// Outcome returns a Result that contains a value or an error.
//
// It can be used to convert a function that returns a value and an error into a
// Result.
//
// Deprecated: Use From instead.
func Outcome[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// Errf returns a Result that contains a formatted error.
func Errf[T any](format string, args ...any) Result[T] {
	return Err[T](fmt.Errorf(format, args...))
}

// Errorf returns a Result that contains a formatted error.
//
// Deprecated: Use Errf instead.
func Errorf[T any](format string, args ...any) Result[T] {
	return Err[T](fmt.Errorf(format, args...))
}

// A Result type is a type that can contain an error or a value.
type Result[T any] struct {
	value T
	err   error
}

// Get returns the value and a boolean indicating whether the value is present.
func (r Result[T]) Get() (T, bool) { return r.value, r.err == nil }

// Result returns the underlying value and error.
func (r Result[T]) Result() (T, error) { return r.value, r.err }

// Default returns the Result value if it is present, otherwise it returns the
// value passed.
func (r Result[T]) Default(value T) T {
	if r.err == nil {
		return r.value
	}
	return value
}

// Err returns the error, if any.
func (r Result[T]) Err() error { return r.err }

func (r Result[T]) String() string {
	if r.err == nil {
		return fmt.Sprintf("%v", r.value)
	}
	return fmt.Sprintf("error(%q)", r.err.Error())
}

func (r Result[T]) GoString() string {
	if r.err == nil {
		return fmt.Sprintf("Ok[%T](%#v)", r.value, r.value)
	}
	return fmt.Sprintf("Err[%T](%q)", r.value, r.err)
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	value := map[string]any{}
	if r.err == nil {
		value["value"] = r.value
		return json.Marshal(value)
	} else {
		value["error"] = r.err.Error()
	}
	return json.Marshal(value)
}

func (r *Result[T]) UnmarshalJSON(data []byte) error {
	value := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if v, ok := value["value"]; ok {
		return json.Unmarshal(v, &r.value)
	}
	if v, ok := value["error"]; ok {
		var errStr string
		if err := json.Unmarshal(v, &errStr); err != nil {
			return err
		}
		r.err = errors.New(errStr)
		return nil
	}
	return errors.New("invalid result")
}
