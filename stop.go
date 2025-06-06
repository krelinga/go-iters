package iters

import "iter"

// Stop halts the processing of the sequence.
// This is a convenient way to discard all elements in the sequence without processing them.
func Stop[T any](seq iter.Seq[T]) {
	for range seq {
		return
	}
}

// Stop2 halts the processing of the sequence.
// This is a convenient way to discard all elements in the sequence without processing them.
func Stop2[T1, T2 any](seq iter.Seq2[T1, T2]) {
	for range seq {
		return
	}
}

// StopIf stops processing the sequence when the predicate returns true.
// It passes-through all elements until the predicate is true, then stops yielding further values.
func StopIf[T any](seq iter.Seq[T], pred PredFn[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val := range seq {
			if pred(val) {
				return // Stop processing if the predicate is true
			}
			if !yield(val) {
				return // Stop yielding values if the yield function returns false
			}
		}
	}
}

// StopIf2 stops processing the sequence when the predicate returns true.
// It passes-through all elements until the predicate is true, then stops yielding further values.
func StopIf2[T1, T2 any](seq iter.Seq2[T1, T2], pred PredFn2[T1, T2]) iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
		for one, two := range seq {
			if pred(one, two) {
				return // Stop processing if the predicate is true
			}
			if !yield(one, two) {
				return // Stop yielding values if the yield function returns false
			}
		}
	}
}
