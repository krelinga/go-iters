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

// PullFromChan creates a pull-based iterator from a channel.
// This may be more-efficient than `FromChan` combined with `iter.Pull`.
// `done` is closed the first time the `stop` function is called.
// `next` pulls a value off the `data` channel and returns it as long as the `data` channel is open and `stop` has not been called.
// See the comments on `iter.Pull` for more details on the semantics of `next` and `stop`, PullFromChan has the same behavior.
func PullFromChan[T any](data <-chan T, done chan<- struct{}) (next func() (T, bool), stop func()) {
	var stopCalled bool
	next = func() (T, bool) {
		if stopCalled {
			var zero T
			return zero, false // If stop was called, return zero value and false
		}
		val, ok := <-data
		return val, ok
	}
	stop = func() {
		if !stopCalled {
			stopCalled = true
			close(done) // Close the done channel to signal that the iterator is done
		}
	}
	return
}
