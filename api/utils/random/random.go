package random

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"

	"github.com/seehuhn/mt19937"
)

func GenerateRandomString(length int) (string, error) {
	const charas = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	clen := len(charas)
	r := make([]byte, length)
	for i := range r {
		rand_n, err := GenSecRand(clen)
		if err != nil {
			return "", err
		}
		r[i] = charas[rand_n]
	}
	return string(r), nil
}

func GenSecRand(k int) (int64, error) {
	seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return 0, err
	}
	generator := rand.New(mt19937.New())
	generator.Seed(seed.Int64())
	return generator.Int63n(int64(k)), nil
}