package iters

import "iter"

// SplitOne returns an iterator that yields only the first element of each pair in the input sequence.
func SplitOne[T1, T2 any](in iter.Seq2[T1, T2]) iter.Seq[T1] {
	return func(yield func(T1) bool) {
		for one, _ := range in {
			if !yield(one) {
				return
			}
		}
	}
}

// SplitTwo returns an iterator that yields only the second element of each pair in the input sequence.
func SplitTwo[T1, T2 any](in iter.Seq2[T1, T2]) iter.Seq[T2] {
	return func(yield func(T2) bool) {
		for _, two := range in {
			if !yield(two) {
				return
			}
		}
	}
}

// Split splits an `iter.Seq2` into two separate `iter.Seq`s, one for each element of the pair.
// Callers *MUST* ensure that the two returned sequences are consumed in parallel to avoid deadlocks.
func Split[T1, T2 any](in iter.Seq2[T1, T2]) (iter.Seq[T1], iter.Seq[T2]) {
	oneDone := make(chan struct{})
	twoDone := make(chan struct{})
	oneData := make(chan T1)
	twoData := make(chan T2)

	go func() {
		var oneIsDone, twoIsDone bool
		for one, two := range in {
			if !oneIsDone {
				select {
				case oneData <- one:
				case <-oneDone:
					oneIsDone = true
					close(oneData)
				}
			}
			if !twoIsDone {
				select {
				case twoData <- two:
				case <-twoDone:
					twoIsDone = true
					close(twoData)
				}
			}
			if oneIsDone && twoIsDone {
				break
			}
		}
		if !oneIsDone {
			close(oneData)
		}
		if !twoIsDone {
			close(twoData)
		}
	}()

	out1 := func(yield func(T1) bool) {
		defer close(oneDone)
		for one := range oneData {
			if !yield(one) {
				return
			}
		}
	}

	out2 := func(yield func(T2) bool) {
		defer close(twoDone)
		for two := range twoData {
			if !yield(two) {
				return
			}
		}
	}

	return out1, out2
}