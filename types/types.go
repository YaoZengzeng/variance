package types

type Pool struct {
	Name  string
	Value int64
}

type Pools []Pool

func (p Pools) Len() int { return len(p) }

func (p Pools) Less(i, j int) bool { return p[i].Value > p[j].Value }

func (p Pools) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
