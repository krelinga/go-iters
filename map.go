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