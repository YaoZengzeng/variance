package variance

import (
	"sort"

	"github.com/YaoZengzeng/variance/types"
)

func MinVariance(P types.Pools, f int64) (map[string]int64, error) {
	sort.Sort(P)

	var (
		aggr, remainder, pos, delta int64
	)

	for pos = 1; pos <= int64(len(P)); pos++ {
		if pos == int64(len(P)) {
			aggr += f / pos
			remainder = f % pos
			break
		}
		if P[pos] != P[pos-1] {
			delta = P[pos-1].Value - P[pos].Value
			if f <= delta*pos {
				aggr += (f / pos)
				remainder = f % pos
				break
			} else {
				aggr += delta
				f -= (delta * pos)
			}
		}
	}

	delta = 0
	result := map[string]int64{}
	for i := int64(0); i < pos; i++ {
		result[P[i].Name] = aggr
		if i != 0 {
			delta += (P[i-1].Value - P[i].Value)
			result[P[i].Name] -= delta
		}
		if remainder != 0 {
			result[P[i].Name] += 1
			remainder--
		}
	}

	return result, nil
}
