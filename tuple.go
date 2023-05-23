package types

type Pair[T, U any] struct {
	First  T
	Second U
}

func PairOf[T, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{First: first, Second: second}
}

func (p Pair[T, U]) Get() (T, U) {
	return p.First, p.Second
}

type Triple[T, U, V any] struct {
	First  T
	Second U
	Third  V
}

func TripleOf[T, U, V any](first T, second U, third V) Triple[T, U, V] {
	return Triple[T, U, V]{First: first, Second: second, Third: third}
}

func (t Triple[T, U, V]) Get() (T, U, V) {
	return t.First, t.Second, t.Third
}
