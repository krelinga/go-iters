package iters

import "iter"

type Sink[T any] struct {
	Pred PredFn[T]
	Out  *iter.Seq[T]
}

func NewSink[T any](pred PredFn[T], out *iter.Seq[T]) Sink[T] {
	return Sink[T]{Pred: pred, Out: out}
}

// Partition takes a sequence and multiple sinks, each with a pointer to an iterator and a corresponding predicate.
// It begins by creating & assigning an iterator for each sink's Out.
// It then yields elements from the sequence into the first iterator whose predicate matches the element.
// If no predicate matches, the element is yielded to the returned (default) sequence.
// If any consumer abandons an iterator before it is completed then any elements that would have been written to that sink will be lost.
// Callers *MUST* ensure that all output iterators (including the default) are consumed in parallel to avoid deadlocks.
func Partition[T any](seq iter.Seq[T], sinks ...Sink[T]) iter.Seq[T] {
	if len(sinks) == 0 {
		return seq // No sinks, return the original sequence
	}

	// One output channel & done channel for each sink, and the same for the default output.
	sinkData := make([]chan T, len(sinks))
	for i := range sinks {
		sinkData[i] = make(chan T)
	}
	defaultData := make(chan T)
	sinkDone := make([]chan struct{}, len(sinks))
	for i := range sinks {
		sinkDone[i] = make(chan struct{})
	}
	defaultDone := make(chan struct{})

	// Start a goroutine to consume the input sequence and distribute elements to the appropriate sinks.
	go func() {
		var defaultIsDone bool
		sinkIsDone := make([]bool, len(sinks))
	val:
		for val := range seq {
			allDone := true
			for _, done := range sinkIsDone {
				allDone = allDone && done
			}
			allDone = allDone && defaultIsDone
			if allDone {
				return // All sinks and the default iterator are done, stop processing
			}

			// Check each sink's predicate and send the value to the first matching sink.
			for i, sink := range sinks {
				if !sink.Pred(val) {
					continue // This sink's predicate does not match, skip to the next sink
				}
				if sinkIsDone[i] {
					continue val // This sink's predicate matches, but the iterator is already done.  Disacrd the value.
				}
				select {
				case sinkData[i] <- val:
				case <-sinkDone[i]:
					sinkIsDone[i] = true // The sink's iterator is done, stop sending to it
					close(sinkData[i])   // Close the channel to signal no more data will be sent
				}
				continue val // Value has been sent to a sink, skip to the next value
			}

			// If we got here, then no sink's predicate matched the value.  Output to the default iterator.
			if defaultIsDone {
				continue // Default iterator is already done, discard the value
			}
			select {
			case defaultData <- val:
			case <-defaultDone:
				defaultIsDone = true // The default iterator is done, stop sending to it
				close(defaultData)   // Close the channel to signal no more data will be sent
			}
		}

		// If we got here then it means that the input sequence has been fully consumed and at least one sink is still open.
		for i, done := range sinkIsDone {
			if !done {
				close(sinkData[i]) // Close each sink's channel to signal no more data will be sent
			}
		}
		if !defaultIsDone {
			close(defaultData) // Close the default channel to signal no more data will be sent
		}
	}()

	// Start goroutines for each sink to consume data.
	iteratorFn := func(data <-chan T, done chan<- struct{}) func(yield func(T) bool) {
		return func(yield func(T) bool) {
			defer close(done)
			for val := range data {
				if !yield(val) {
					return // Stop yielding if the consumer stops
				}
			}
		}
	}
	for i, sink := range sinks {
		*sink.Out = iteratorFn(sinkData[i], sinkDone[i])
	}
	return iteratorFn(defaultData, defaultDone)
}
