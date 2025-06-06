package iters_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/krelinga/go-iters"
)

func TestConcat(t *testing.T) {
	seq1 := slices.Values([]int{1, 2, 3})
	seq2 := slices.Values([]int{4, 5, 6})

	want := []int{1, 2, 3, 4, 5, 6}
	got := slices.Collect(iters.Concat(seq1, seq2))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Concat(%v, %v) = %v; want %v", seq1, seq2, got, want)
	}
}
