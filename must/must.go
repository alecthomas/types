// Package must provides a way to panic if the result of a multi-value function
// call returns an error, or return just the values if there is no error.
package must

// Get panics if err is not nil, otherwise it returns value.
func Get[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// Get2 panics if err is not nil, otherwise it returns left and right.
func Get2[T, U any](left T, right U, err error) (T, U) {
	if err != nil {
		panic(err)
	}
	return left, right
}
