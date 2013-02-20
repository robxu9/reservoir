package reservoir

import (
	"fmt"
	"rand"
	"time"
)

type RS struct {
	pool string
	rg   *rand.Rand
	used map[string]int
}

func NewAlphaNumericRS() *RS {
	return NewRS("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
}

func NewRS(pool string) *RS {
	return &RS{
		pool,
		rand.New(rand.NewSource(time.Nanoseconds())),
		make(map[string]int),
	}
}

func (rs *RS) NewRandomString(length int) (r string) {
	if length < 1 {
		return
	}
	b := make([]byte, length)
	for retries := 0; ; retries++ {
		for i, _ := range b {
			b[i] = rs.pool[rs.rg.Intn(len(rs.pool))]
		}
		r = string(b)
		_, used := rs.used[r]
		if !used {
			break
		}
		if retries == 3 {
			return ""
		}
	}
	rs.used[r] = 0
	return
}
