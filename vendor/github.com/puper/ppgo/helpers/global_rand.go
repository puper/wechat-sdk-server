package helpers

import (
	"math/rand"
	"time"
)

func GlobalRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomString(l int) string {
	bs := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lenBs := len(bs)
	result := make([]byte, 0, l)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bs[r.Intn(lenBs)])
	}
	return string(result)
}
