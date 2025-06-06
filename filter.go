package iters

import "iter"

// PredFn is a predicate function type that takes an element of type T and returns true if the element satisfies the condition, false otherwise.
type PredFn[T any] func(T) bool

// Filter returns an iterator that yields elements from the input sequence that satisfy the given predicate function.
func Filter[T any](seq iter.Seq[T], predicate PredFn[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for val := range seq {
			if predicate(val) && !yield(val) {
				return
			}
		}
	}
}

// PredFn2 is a predicate function type that takes two elements of types T1 and T2 and returns true if the pair satisfies the condition, false otherwise.
type PredFn2[T1, T2 any] func(T1, T2) bool

// Filter2 returns an iterator that yields pairs from the input sequence that satisfy the given predicate function.
func Filter2[T1, T2 any](seq iter.Seq2[T1, T2], predicate PredFn2[T1, T2]) iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
		for one, two := range seq {
			if predicate(one, two) && !yield(one, two) {
				return
			}
		}
	}
}
