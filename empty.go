package iters

import "iter"

// Empty returns an iterator that yields no elements.
func Empty[T any]() iter.Seq[T] {
	return func(yield func(T) bool) {}
}

// Empty2 returns an iterator that yields no elements.
func Empty2[T1, T2 any]() iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {}
}
