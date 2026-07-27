package cat

import (
	"reflect"
	"testing"
)

func TestEffectiveMaxFlowsPerMessage(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in, want int
	}{
		{0, 1},
		{-5, 1},
		{1, 1},
		{100, 100},
	}
	for _, tc := range cases {
		if got := effectiveMaxFlowsPerMessage(tc.in); got != tc.want {
			t.Fatalf("effectiveMaxFlowsPerMessage(%d) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestJchfExportSliceBounds(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		batchSize int
		keep      int
		want      [][2]int
	}{
		{"single", 10000, 3, [][2]int{{0, 3}}},
		{"exact", 4, 8, [][2]int{{0, 4}, {4, 8}}},
		{"partial-last", 4, 10, [][2]int{{0, 4}, {4, 8}, {8, 10}}},
		{"zero-batch-clamped", 0, 5, [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := jchfExportSliceBounds(tc.batchSize, tc.keep)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("jchfExportSliceBounds(%d, %d) = %v, want %v", tc.batchSize, tc.keep, got, tc.want)
			}
		})
	}
}
