package iters

import "iter"

type AnySeq[T any] interface {
	iter.Seq[T] | ParSeq[T]
}

type AnySeq2[T1, T2 any] interface {
	iter.Seq2[T1, T2] | ParSeq2[T1, T2]
}