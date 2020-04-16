package rng

import (
	cryptorand "crypto/rand"
	big "math/big"
	mathrand "math/rand"
)

type RNG interface {
	RandInt(max int) int
}

type SecureRNG struct{}

func (s *SecureRNG) RandInt(max int) int {
	res, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err.Error())
	}
	return int(res.Uint64())
}

var _ RNG = (*SecureRNG)(nil)

type SeededRNG struct {
	Rand *mathrand.Rand
}

func (i *SeededRNG) RandInt(max int) int {
	return i.Rand.Intn(max)
}

var _ RNG = (*SeededRNG)(nil)
