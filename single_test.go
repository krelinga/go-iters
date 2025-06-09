package iters_test

import (
	"maps"
	"slices"
	"testing"

	"github.com/krelinga/go-iters"
)

func TestSingle(t *testing.T) {
	got := slices.Collect(iters.Single(int(42)))
	want := []int{42}
	if !slices.Equal(got, want) {
		t.Errorf("Single(42) = %v; want %v", got, want)
	}
}

func TestSingle2(t *testing.T) {
	got := maps.Collect(iters.Single2("hello", "world"))
	want := map[string]string{"hello": "world"}
	if !maps.Equal(got, want) {
		t.Errorf("Single2(\"hello\", \"world\") = %v; want %v", got, want)
	}
}
