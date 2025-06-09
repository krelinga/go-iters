package iters

import "iter"

// Single returns an iterator that yields t as its only element.
func Single[T any](t T) iter.Seq[T] {
	return func(yield func(T) bool) {
		yield(t)
	}
}

// Single2 returns an iterator that yields (1, 2) as its only element.
func Single2[T1, T2 any](one T1, two T2) iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
		yield(one, two)
	}
}
