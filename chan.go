package iters

import "iter"

// FromChan creates an iterator from a channel.
// It consumes values from the channel and yields them until the channel is closed, or until the iterator is no-longer needed.
// The `done` channel is closed when the iterator completes (either because the data channel is closed or the iterator is no-longer needed).
func FromChan[T any](data <-chan T, done chan<- struct{}) iter.Seq[T] {
	return func(yield func(T) bool) {
		defer close(done) // Ensure the done channel is closed when the iterator completes
		for val := range data {
			if !yield(val) {
				return
			}
		}
	}
}