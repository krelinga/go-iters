package iters

import "iter"

// ParSeq is used as a type alias for iter.Seq[T] to indicate that it is a sequence that the caller *MUST* consume in parallel to avoid deadlocks.
type ParSeq[T any] iter.Seq[T]

// ParSeq2 is used as a type alias for iter.Seq2[T1, T2] to indicate that it is a sequence that the caller *MUST* consume in parallel to avoid deadlocks.
type ParSeq2[T1, T2 any] iter.Seq2[T1, T2]

// WaitFn is a function type that is used to wait for a ParSeq or ParSeq2 to be consumed by a parallel goroutine.
// Calls to a WaitFn block until the corresponding sequence is fully consumed.
// It is safe to call a WaitFn multiple times; calls after the first will return immediately.
type WaitFn func()

// ParConsume starts a new goroutine to consume the provided ParSeq[T] using the provided function f.
// It returns a WaitFn that can be called to block until the sequence is fully consumed.
func ParConsume[T any](seq ParSeq[T], f func(iter.Seq[T])) WaitFn {
	done := make(chan struct{})
	go func() {
		defer close(done)
		f(iter.Seq[T](seq))
	}()
	return func() {
		<-done
	}
}

// ParConsume2 starts a new goroutine to consume the provided ParSeq2[T1, T2] using the provided function f.
// It returns a WaitFn that can be called to block until the sequence is fully consumed.
func ParConsume2[T1, T2 any](seq ParSeq2[T1, T2], f func(iter.Seq2[T1, T2])) WaitFn {
	done := make(chan struct{})
	go func() {
		defer close(done)
		f(iter.Seq2[T1, T2](seq))
	}()
	return func() {
		<-done
	}
}

// Wait invokes each of the provided WaitFns and blocks until all of them are completed.
func Wait(waitFns ...WaitFn) {
	for _, waitFn := range waitFns {
		waitFn()
	}
}