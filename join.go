package iters

import "iter"

// JoinPad combines two `iter.Seq`s into a single `iter.Seq2`.
// It stops when the longer sequence is exhausted, padding the shorter one with zero values.
// It consumes elements from both sequences in parallel.
func JoinPad[T1, T2 any](one iter.Seq[T1], two iter.Seq[T2]) iter.Seq2[T1, T2] {
	return joinImpl(one, two, func(oneOk, twoOk bool) bool {
		return !oneOk && !twoOk // Stop when both sequences are exhausted
	})
}

// JoinTrim combines two `iter.Seq`s into a single `iter.Seq2`.
// It stops when the shorter sequence is exhausted, ignoring any remaining elements in the longer sequence.
// It consumes elements from both sequences in parallel.
func JoinTrim[T1, T2 any](one iter.Seq[T1], two iter.Seq[T2]) iter.Seq2[T1, T2] {
	return joinImpl(one, two, func(oneOk, twoOk bool) bool {
		return !oneOk || !twoOk // Stop when either sequence is exhausted
	})
}

func joinImpl[T1, T2 any](one iter.Seq[T1], two iter.Seq[T2], stopWhen func(bool, bool) bool) iter.Seq2[T1, T2] {
	return func(yield func (T1, T2) bool) {
		oneGet, oneDone := iter.Pull(one)
		twoGet, twoDone := iter.Pull(two)
		defer oneDone()
		defer twoDone()

		for {
			oneVal, oneOk := oneGet()
			twoVal, twoOk := twoGet()

			if stopWhen(oneOk, twoOk) {
				break // Stop when the condition is met
			}

			if !yield(oneVal, twoVal) {
				return // Stop yielding if the consumer stops
			}
		}
	}
}
