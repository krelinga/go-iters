package iters

import "iter"

type Pair[T1, T2 any] struct {
	One T1
	Two T2
}

// NewPair creates a new Pair with the given values.
// This is a convenience function to avoid having to specify the type parameters explicitly.
func NewPair[T1, T2 any](one T1, two T2) Pair[T1, T2] {
	return Pair[T1, T2]{One: one, Two: two}
}

// ToPairs converts an `iter.Seq2` into an `iter.Seq` of `Pair`.
func ToPairs[T1, T2 any](in iter.Seq2[T1, T2]) iter.Seq[Pair[T1, T2]] {
	return func(yield func(Pair[T1, T2]) bool) {
		for one, two := range in {
			if !yield(Pair[T1, T2]{One: one, Two: two}) {
				return
			}
		}
	}
}

// FromPairs converts an `iter.Seq` of `Pair` into an `iter.Seq2`.
func FromPairs[T1, T2 any](in iter.Seq[Pair[T1, T2]]) iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
		for pair := range in {
			if !yield(pair.One, pair.Two) {
				return
			}
		}
	}
}
