// Package random gives Go code cryptographically secure randomness through an
// interface that tests can replace. Package randomtest has the mock.
package random

import (
	"crypto/rand"
	"io"

	"github.com/The127/go-random/internal/derive"
)

// Service returns random values. If crypto/rand fails, the program crashes,
// as crypto/rand.Read does since Go 1.24, and recover can't catch it.
type Service interface {
	// Bytes returns n random bytes. It panics if n is negative.
	Bytes(n int) []byte
	// Text returns 26 characters of the base32 alphabet, A to Z and 2 to 7,
	// which hold more than 128 bits of randomness.
	Text() string
	// IntN returns a uniformly random number from 0 to n-1. It panics if n is
	// below 1.
	IntN(n int) int
	// Float64 returns a uniformly random number from 0 to below 1.
	Float64() float64
	// Reader returns the source of the random bytes, for libraries that read
	// from an io.Reader. Since Go 1.26, Go's own crypto packages, like
	// crypto/ecdsa, ignore a reader they're given and use crypto/rand.
	Reader() io.Reader
}

type source struct{ r io.Reader }

// NewCryptoRandom returns a Service that reads crypto/rand. Use it for
// everything secret, like tokens, salts and keys.
func NewCryptoRandom() Service {
	return source{rand.Reader}
}

// NewPseudoRandom returns a fast Service for values that aren't secret, for
// example to spread out retries. Its ChaCha8 stream takes its seed from
// crypto/rand and differs on every start.
func NewPseudoRandom() Service {
	var seed [32]byte

	_, _ = rand.Read(seed[:])

	return source{derive.NewStream(seed)}
}

// NewPseudoRandomFromSeed returns a Service like NewPseudoRandom whose values
// depend only on seed and are the same in every run. Never use it for
// anything secret.
func NewPseudoRandomFromSeed(seed uint64) Service {
	return source{derive.NewStream(derive.Seed(seed))}
}

func (s source) Bytes(n int) []byte { return derive.Bytes(s.r, n) }

func (s source) Text() string { return derive.Text(s.r) }

func (s source) IntN(n int) int { return derive.IntN(s.r, n) }

func (s source) Float64() float64 { return derive.Float64(s.r) }

func (s source) Reader() io.Reader { return s.r }
