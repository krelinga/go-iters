package iters

import "iter"

// Concat combines multiple `iter.Seq`s into a single `iter.Seq`.
// It yields elements from each sequence in the order they are provided.
func Concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range seqs {
			for val := range seq {
				if !yield(val) {
					return
				}
			}
		}
	}
}

// Concat2 combines multiple `iter.Seq2`s into a single `iter.Seq2`.
// It yields pairs of elements from each sequence in the order they are provided.
func Concat2[T1, T2 any](seqs ...iter.Seq2[T1, T2]) iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
		for _, seq := range seqs {
			for one, two := range seq {
				if !yield(one, two) {
					return
				}
			}
		}
	}
}