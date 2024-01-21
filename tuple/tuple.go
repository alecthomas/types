// Package tuple provides pseudo-tuple types.
package tuple

type Pair[A, B any] struct {
	A A
	B B
}

func PairOf[A, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{A: a, B: b}
}

func (p Pair[A, B]) Get() (A, B) {
	return p.A, p.B
}

type Triple[A, B, C any] struct {
	A A
	B B
	C C
}

func TripleOf[A, B, C any](a A, b B, c C) Triple[A, B, C] {
	return Triple[A, B, C]{A: a, B: b, C: c}
}

func (t Triple[A, B, C]) Get() (A, B, C) {
	return t.A, t.B, t.C
}

type Quad[A, B, C, W any] struct {
	A A
	B B
	C C
	D W
}

func QuadOf[A, B, C, D any](a A, b B, c C, d D) Quad[A, B, C, D] {
	return Quad[A, B, C, D]{A: a, B: b, C: c, D: d}
}

func (q Quad[A, B, C, W]) Get() (A, B, C, W) {
	return q.A, q.B, q.C, q.D
}
