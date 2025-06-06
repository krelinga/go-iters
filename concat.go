package iters

// Concat combines multiple `iter.Seq`s into a single `iter.Seq`.
// It yields elements from each sequence in the order they are provided.
func Concat[Seq AnySeq[T], T any](seqs ...Seq) Seq {
	return func(yield func(T) bool) {
		for _, seq := range seqs {
			for val := range seq {
				if !yield(val) {
					return
				}
			}
		}
	}
}

// Concat2 combines multiple `iter.Seq2`s into a single `iter.Seq2`.
// It yields pairs of elements from each sequence in the order they are provided.
func Concat2[Seq AnySeq2[T1, T2], T1, T2 any](seqs ...Seq) Seq {
	return func(yield func(T1, T2) bool) {
		for _, seq := range seqs {
			for one, two := range seq {
				if !yield(one, two) {
					return
				}
			}
		}
	}
}