package variance

import (
	"reflect"
	"testing"

	"github.com/YaoZengzeng/variance/types"
)

func TestMinVariance(t *testing.T) {
	tests := []struct {
		f        int64
		pools    types.Pools
		expected map[string]int64
	}{
		// All pi cut to the same value.
		{
			f: 40,
			pools: types.Pools{
				{
					Name:  "f1",
					Value: 120,
				},
				{
					Name:  "f2",
					Value: 100,
				},
			},
			expected: map[string]int64{
				"f1": 30,
				"f2": 10,
			},
		},
		// All pi cut to 0.
		{
			f: 40,
			pools: types.Pools{
				{
					Name:  "f1",
					Value: 30,
				},
				{
					Name:  "f2",
					Value: 10,
				},
			},
			expected: map[string]int64{
				"f1": 30,
				"f2": 10,
			},
		},
		// All pi cut to negative number.
		{
			f: 40,
			pools: types.Pools{
				{
					Name:  "f1",
					Value: 20,
				},
				{
					Name:  "f2",
					Value: 10,
				},
			},
			expected: map[string]int64{
				"f1": 25,
				"f2": 15,
			},
		},
		// remainder exists.
		{
			f: 41,
			pools: types.Pools{
				{
					Name:  "f1",
					Value: 20,
				},
				{
					Name:  "f2",
					Value: 10,
				},
			},
			expected: map[string]int64{
				"f1": 26,
				"f2": 15,
			},
		},
		// some pi are equal.
		{
			f: 40,
			pools: types.Pools{
				{
					Name:  "f1",
					Value: 15,
				},
				{
					Name:  "f2",
					Value: 10,
				},
				{
					Name:  "f3",
					Value: 10,
				},
				{
					Name:  "f4",
					Value: 5,
				},
			},
			expected: map[string]int64{
				"f1": 15,
				"f2": 10,
				"f3": 10,
				"f4": 5,
			},
		},
		// Only cut some of pi.
		{
			f: 30,
			pools: types.Pools{
				{
					Name:  "f1",
					Value: 100,
				},
				{
					Name:  "f2",
					Value: 100,
				},
				{
					Name:  "f3",
					Value: 10,
				},
			},
			expected: map[string]int64{
				"f1": 15,
				"f2": 15,
			},
		},
	}

	for _, tt := range tests {
		get, err := MinVariance(tt.pools, tt.f)
		if err != nil {
			t.Fatalf("failed to call MinVariance: %v", err)
		}

		if !reflect.DeepEqual(tt.expected, get) {
			t.Fatalf("expected to get %v, but actually get %v", tt.expected, get)
		}
	}
}
