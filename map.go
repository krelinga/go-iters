package iters

import "iter"

func Map[T, R any](seq iter.Seq[T], fn func(T) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		for val := range seq {
			if !yield(fn(val)) {
				return
			}
		}
	}
}

func Map2[T1, T2, R1, R2 any](seq iter.Seq2[T1, T2], fn func(T1, T2) (R1, R2)) iter.Seq2[R1, R2] {
	return func(yield func(R1, R2) bool) {
		for one, two := range seq {
			if !yield(fn(one, two)) {
				return
			}
		}
	}
}
