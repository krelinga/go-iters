package iters_test

import (
	"maps"
	"slices"
	"testing"

	"github.com/krelinga/go-iters"
)

func TestEmpty(t *testing.T) {
	out := slices.Collect(iters.Empty[int]())
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %v", out)
	}
}

func TestEmpty2(t *testing.T) {
	out := maps.Collect(iters.Empty2[int, string]())
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %v", out)
	}
}
