package iters

import "iter"

type PartDest[T any] struct {
	Pred PredFn[T]
	Sink Sink[T]
}

func NewPartDest[T any](pred PredFn[T], out Sink[T]) PartDest[T] {
	return PartDest[T]{Pred: pred, Sink: out}
}

// Partition distributes elements to multiple sinks based on predicates.
// Predicates are evaluated in-order, and the first matching predicate's sink will receive the value.
// If a sink does not accept the value then the value will be discarded.
// If no predicate matches, the value is sent to iterator returned by Partition.
func Partition[T any](seq iter.Seq[T], dests ...PartDest[T]) iter.Seq[T] {
	if len(dests) == 0 {
		return seq // No destinations, return the original sequence
	}

	return func(yield func(T) bool) {
		destDone := make([]bool, len(dests))
		var defDone bool
		var doneCount int

		tryWrite := func(sink Sink[T], val T, done *bool) {
			if *done {
				return // Sink is closed, do not write
			}
			if !sink.Write(val) {
				*done = true // Sink is full or closed, mark it as closed
				sink.Close() // Close the sink to signal no more values will be sent
				doneCount++
			}
		}

	vals:
		for val := range seq {
			if doneCount == len(dests)+1 {
				// All sinks + return iterator are closed, stop processing
				break
			}

			for i, predSink := range dests {
				if predSink.Pred(val) {
					tryWrite(predSink.Sink, val, &destDone[i])
					continue vals // Move to the next value after writing to a matching sink
				}
			}
			// If no predicate matched, write to the default sink
			if !defDone && !yield(val) {
				defDone = true // Mark the default sink as closed
				doneCount++
			}
		}
		for i, done := range destDone {
			if !done {
				dests[i].Sink.Close() // Close any sinks that were not closed during processing
			}
		}
	}
}
