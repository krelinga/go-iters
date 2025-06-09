package iters_test

import (
	"maps"
	"slices"
	"testing"

	"github.com/krelinga/go-iters"
)

func TestNewPair(t *testing.T) {
	pair := iters.NewPair[int, string](42, "hello")

	if pair.One != 42 {
		t.Errorf("Expected One to be 42, got %d", pair.One)
	}

	if pair.Two != "hello" {
		t.Errorf("Expected Two to be 'hello', got '%s'", pair.Two)
	}
}

func TestFromPairs(t *testing.T) {
	pairs := []iters.Pair[int, string]{
		{One: 1, Two: "one"},
		{One: 2, Two: "two"},
	}

	got := maps.Collect(iters.FromPairs(slices.Values(pairs)))
	want := map[int]string{
		1: "one",
		2: "two",
	}
	if !maps.Equal(got, want) {
		t.Errorf("FromPairs() = %v; want %v", got, want)
	}
}

func TestToPairs(t *testing.T) {
	pairs := map[int]string{
		1: "one",
		2: "two",
	}
	got := slices.Collect(iters.ToPairs(maps.All(pairs)))
	want := []iters.Pair[int, string]{
		{One: 1, Two: "one"},
		{One: 2, Two: "two"},
	}
	if !slices.Equal(got, want) {
		t.Errorf("ToPairs() = %v; want %v", got, want)
	}
}
