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
	return func(yield func(T1, T2) bool) {
		oneNext, oneDone := iter.Pull(iter.Seq[T1](one))
		defer oneDone()
		twoNext, twoDone := iter.Pull(iter.Seq[T2](two))
		defer twoDone()

		done := make(chan struct{})
		defer close(done)

		oneData := make(chan T1)
		go func() {
			defer close(oneData)
			for {
				val, ok := oneNext()
				if !ok {
					return // No more elements in the first sequence
				}
				select {
				case oneData <- val:
				case <-done:
					return // Stop if the main routine has finished
				}
			}
		}()

		twoData := make(chan T2)
		go func() {
			defer close(twoData)
			for {
				val, ok := twoNext()
				if !ok {
					return // No more elements in the second sequence
				}
				select {
				case twoData <- val:
				case <-done:
					return // Stop if the main routine has finished
				}
			}
		}()

		for {
			oneVal, oneOk := <-oneData
			twoVal, twoOk := <-twoData
			if stopWhen(oneOk, twoOk) {
				return // Stop condition met
			}
			if !yield(oneVal, twoVal) {
				return // Yield function returned false, stop iteration
			}
		}
	}
}
