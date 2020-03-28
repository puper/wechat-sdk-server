package helpers

import (
	"errors"
)

var (
	alphabet    = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	alphabetMap = func() map[byte]int {
		result := make(map[byte]int)
		for i, v := range alphabet {
			result[v] = i
		}
		return result
	}()
	InputNotValid = errors.New("input not valid")
)

func Dec2X(a int, x int) (string, error) {
	if x < 2 || x > 62 || a < 0 {
		return "", InputNotValid
	}
	b := make([]byte, 0)
	for {
		b = append(b, alphabet[a%x])
		a = a / x
		if a == 0 {
			break
		}
	}
	l := len(b)
	for i, j := 0, l-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b), nil
}

func X2Dec(a string, x int) (int, error) {
	if x < 2 || x > 62 {
		return 0, InputNotValid
	}
	b := []byte(a)
	l := len(b)
	if l < 1 {
		return 0, InputNotValid
	}
	r := 0
	for _, v := range b {
		t, ok := alphabetMap[v]
		if !ok {
			return 0, InputNotValid
		}
		r = t + x*r
	}
	return r, nil
}
