// Package either provides a simple implementation of a sum type, Either, that can be either a Left or a Right.
//
// Usage:
//
//	left := LeftOf[string](42)
//	fmt.Println(left.Get()) // 42
//
//	right := RightOf[int]("foo")
//	fmt.Println(right.Get()) // foo
//
//	var either Either[int, string] = left
//	either = right
package either

import "fmt"

// Either is a "sum type" that can be either Left or Right.
type Either[L, R any] interface {
	String() string
	GoString() string
	either(L, R) //nolint:inamedparam
}

type Left[L, R any] struct{ value L }

var _ Either[int, string] = Left[int, string]{}

func (Left[L, R]) either(L, R)        {}
func (e Left[L, R]) Get() L           { return e.value }
func (e Left[L, R]) String() string   { return fmt.Sprintf("%v", e.value) }
func (e Left[L, R]) GoString() string { var r R; return fmt.Sprintf("LeftOf[%T](%#v)", r, e.value) }

type Right[L, R any] struct{ value R }

var _ Either[int, string] = Right[int, string]{}

func (Right[L, R]) either(L, R)        {}
func (e Right[L, R]) Get() R           { return e.value }
func (e Right[L, R]) String() string   { return fmt.Sprintf("%v", e.value) }
func (e Right[L, R]) GoString() string { var l L; return fmt.Sprintf("RightOf[%T](%#v)", l, e.value) }

// LeftOf creates an Either[L, R] with a left value.
//
// Note that the L and R type parameters are flipped so that we can use type
// inference to avoid having to specify the L type.
func LeftOf[R, L any](value L) Left[L, R] { return Left[L, R]{value} }

// RightOf creates an Either[L, R] with a right value.
func RightOf[L, R any](value R) Right[L, R] { return Right[L, R]{value} }
