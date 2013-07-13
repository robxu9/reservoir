package util

import (
	"math/rand"
	"time"
)

type RandomGenerator struct {
	pool string
	rg   *rand.Rand
	used map[string]int
}

func Random_NewAlphaNumericGenerator() *RandomGenerator {
	return NewGenerator("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
}

func Random_NewGenerator(pool string) *RandomGenerator {
	return &RandomGenerator{
		pool,
		rand.New(rand.NewSource(time.Now().UnixNano())),
		make(map[string]int),
	}
}

func (g *RandomGenerator) NewRandomString(length int) (r string) {
	if length < 1 {
		return
	}
	b := make([]byte, length)
	for retries := 0; ; retries++ {
		for i, _ := range b {
			b[i] = g.pool[g.rg.Intn(len(g.pool))]
		}
		r = string(b)
		_, used := g.used[r]
		if !used {
			break
		}
		if retries == 3 {
			return ""
		}
	}
	g.used[r] = 0
	return
}
