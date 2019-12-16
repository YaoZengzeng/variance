package variance

import (
	"testing"
	"reflect"

	"github.com/YaoZengzeng/variance/types"
)

func TestMinVariance(t *testing.T) {
	tests := []struct {
		f int
		pool types.Pools
		expected map[string]int
	}{
		// All pi cut to the same value.
		{
			40,
			types.Pools{
				{
					"f1",
					120,
				},
				{
					"f2",
					100,
				},
			},
			map[string]int{
				"f1": 30,
				"f2": 10,
			},
		},
		// All pi cut to 0.
		{
			40,
			types.Pools{
				{
					"f1",
					30,
				},
				{
					"f2",
					10,
				},
			},
			map[string]int{
				"f1": 30,
				"f2": 10,
			},
		},
		// All pi cut to negative number.
		{
			40,
			types.Pools{
				{
					"f1",
					20,
				},
				{
					"f2",
					10,
				},
			},
			map[string]int{
				"f1": 25,
				"f2": 15,
			},
		},
		// remainder exists.
		{
			41,
			types.Pools{
				{
					"f1",
					20,
				},
				{
					"f2",
					10,
				},
			},
			map[string]int{
				"f1": 26,
				"f2": 15,
			},
		},
		// some pi are equal.
		{
			40,
			types.Pools{
				{
					"f1",
					15,
				},
				{
					"f2",
					10,
				},
				{
					"f3",
					10,
				},
				{
					"f4",
					5,
				},
			},
			map[string]int{
				"f1": 15,
				"f2": 10,
				"f3": 10,
				"f4": 5,
			},
		},
		// Only cut some of pi.
		{
			30,
			types.Pools{
				{
					"f1",
					100,
				},
				{
					"f2",
					100,
				},
				{
					"f3",
					10,
				},
			},
			map[string]int{
				"f1": 15,
				"f2": 15,
			},
		},
	}

	for _, tt := range tests {
		get, err := MinVariance(tt.pool, tt.f)
		if err != nil {
			t.Fatalf("failed to call MinVariance: %v", err)
		}

		if !reflect.DeepEqual(tt.expected, get) {
			t.Fatalf("expected to get %v, but actually get %v", tt.expected, get)
		}
	}
}
