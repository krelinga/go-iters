package iters

import "iter"

// Merge combines multiple `ParSeq`s into a single `iter.Seq`.
// It consumes elements from each sequence in parallel, yielding them as they become available.
func Merge[T any](seqs ...ParSeq[T]) iter.Seq[T] {
	done := make(chan struct{})
	data := make(chan T)

	for _, seq := range seqs {
		ParConsume(seq, func(s iter.Seq[T]) {
			for val := range s {
				select {
				case data <- val:
				case <-done:
					return
				}
			}
		})
	}

	return func(yield func(T) bool) {
		defer close(done)
		for val := range data {
			if !yield(val) {
				return
			}
		}
	}
}

// Merge2 combines multiple `ParSeq2`s into a single `iter.Seq2`.
// It consumes pairs of elements from each sequence in parallel, yielding them as they become available.
func Merge2[T1, T2 any](seqs ...ParSeq2[T1, T2]) iter.Seq2[T1, T2] {
	done := make(chan struct{})
	data := make(chan Pair[T1, T2])

	for _, seq := range seqs {
		ParConsume2(seq, func(s iter.Seq2[T1, T2]) {
			for one, two := range s {
				select {
				case data <- Pair[T1, T2]{one, two}:
				case <-done:
					return
				}
			}
		})
	}

	return func(yield func(T1, T2) bool) {
		defer close(done)
		for pair := range data {
			if !yield(pair.One, pair.Two) {
				return
			}
		}
	}
}