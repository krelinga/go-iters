package iters

import "iter"
import "github.com/krelinga/go-views"

// Flatten consumes a sequence of `views.List`s and yields all elements from all lists in order.
func Flatten[T any](in iter.Seq[views.List[T]]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for list := range in {
			for val := range list.Values() {
				if !yield(val) {
					return
				}
			}
		}
	}
}

// FlattenOne consumes a sequence of `views.List`s paired with a second value and yields all elements from all lists in order, along with the second value.
// The second value is repeated for each element in the corresponding list.
func FlattenOne[T1, T2 any](in iter.Seq2[views.List[T1], T2]) iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
		for list, val := range in {
			for item := range list.Values() {
				if !yield(item, val) {
					return
				}
			}
		}
	}
}

// FlattenTwo consumes a sequence of `views.List`s paired with a first value and yields all elements from all lists in order, along with the first value.
// The first value is repeated for each element in the corresponding list.
func FlattenTwo[T1, T2 any](in iter.Seq2[T1, views.List[T2]]) iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
		for val, list := range in {
			for item := range list.Values() {
				if !yield(val, item) {
					return
				}
			}
		}
	}
}