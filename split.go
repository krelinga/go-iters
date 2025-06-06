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

func Split[T1, T2 any](in iter.Seq2[T1, T2], oneSink Sink[T1], twoSink Sink[T2]) {
	var oneDone, twoDone bool
	for one, two := range in {
		if oneDone && twoDone {
			break // Both sinks are closed, stop processing
		}
		if !oneDone && !oneSink.Write(one) {
			oneDone = true  // Mark the first sink as closed if it cannot accept more values
			oneSink.Close() // Close the sink to signal no more values will be sent
		}
		if !twoDone && !twoSink.Write(two) {
			twoDone = true  // Mark the second sink as closed if it cannot accept more values
			twoSink.Close() // Close the sink to signal no more values will be sent
		}
	}
	if !oneDone {
		oneSink.Close() // Close the first sink if it was not closed during processing
	}
	if !twoDone {
		twoSink.Close() // Close the second sink if it was not closed during processing
	}
}
