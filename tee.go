package iters

import "iter"

// Tee copies values from the input sequence to the given output sinks before yielding them to the returned iterator.
func Tee[T any](seq iter.Seq[T], sinks ...Sink[T]) iter.Seq[T] {
	if len(sinks) == 0 {
		return seq
	}
	return func(yield func(T) bool) {
		closedCount := 0
		closed := make([]bool, len(sinks))
		for i := range closed {
			if sinks[i] == nil {
				closed[i] = true // If the sink is nil, mark it as closed
				closedCount++
			}
		}
		var retClosed bool

		tryWrite := func(sink Sink[T], val T, closed *bool) {
			if *closed {
				return
			}
			if !sink.Write(val) {
				sink.Close()
				*closed = true
				closedCount++
			}
		}

		for val := range seq {
			if closedCount == len(sinks)+1 {
				break // All sinks and default output are closed, exit the loop
			}
			for i, sink := range sinks {
				tryWrite(sink, val, &closed[i])
			}
			if !retClosed && !yield(val) {
				retClosed = true
				closedCount++
				break
			}
		}

		for i, sink := range sinks {
			if !closed[i] {
				sink.Close()
			}
		}
	}
}
