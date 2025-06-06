package iters

// Sink abstracts away the details of of values are written.
type Sink[T any] interface {
	// Close the sink, signaling that no more values will be sent.
	// Should be called exactly once.
	// It is an error to call Write() after Close() has been called.
	Close()

	// Write a value to the sink.
	// Returns true if the value was successfully written, or false if the sink will not accept more values.
	// It is not an error to call Write() after it has returned false; it will simply return false again.
	// It is an error to call Write() after Close() has been called.
	// It is up to the caller to call Close() when they are done writing values, even after Write() returns false.
	Write(T) bool
}

// ToSlice returns a Sink that appends values to the given slice (allocating it if necessary).
func ToSlice[T any](slice *[]T) Sink[T] {
	return toSlice[T]{slice: slice}
}

type toSlice[T any] struct {
	slice *[]T
}

func (s toSlice[T]) Close() {}

func (s toSlice[T]) Write(val T) bool {
	if s.slice == nil {
		return false // Slice is nil, cannot write
	}
	*s.slice = append(*s.slice, val)
	return true
}

// ToChan returns a Sink that sends values to the given data channel.
// If the given done channel is closed then the Sink's Write function will return false without sending a value.
func ToChan[T any](data chan<- T, done <-chan struct{}) Sink[T] {
	return toChan[T]{data: data, done: done}
}

type toChan[T any] struct {
	data chan<- T
	done <-chan struct{}
}

func (s toChan[T]) Close() {
	close(s.data) // Close the channel to signal no more values will be sent
}

func (s toChan[T]) Write(val T) bool {
	select {
	case s.data <- val:
		return true
	case <-s.done:
		return false
	}
}