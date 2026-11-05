// Package derive turns a stream of random bytes into the values that
// random.Service returns. Every source and the mock use it.
package derive

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	mathrand "math/rand/v2"
	"sync"
)

const base32Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

// Bytes reads n bytes from r.
func Bytes(r io.Reader, n int) []byte {
	if n < 0 {
		panic(fmt.Sprintf("random: Bytes(%d) needs n of at least 0", n))
	}

	b := make([]byte, n)

	_, err := io.ReadFull(r, b)
	if err != nil {
		panic(fmt.Sprintf("random: reading random bytes: %v", err))
	}

	return b
}

// Text returns 26 characters of the base32 alphabet, which hold 130 bits of
// randomness, the same way crypto/rand.Text does.
func Text(r io.Reader) string {
	b := Bytes(r, 26)
	for i := range b {
		b[i] = base32Alphabet[b[i]%32]
	}

	return string(b)
}

// IntN returns a uniformly random number from 0 to n-1. It panics if n is
// below 1.
func IntN(r io.Reader, n int) int {
	if n < 1 {
		panic(fmt.Sprintf("random: IntN(%d) needs n of at least 1", n))
	}

	// crypto/rand.Int rejects values that would favor small numbers, which a
	// plain modulo doesn't.
	v, err := rand.Int(r, big.NewInt(int64(n)))
	if err != nil {
		panic(fmt.Sprintf("random: reading random bytes: %v", err))
	}

	return int(v.Int64())
}

// Float64 returns a uniformly random number from 0 to below 1. 53 random bits
// fill a float64's mantissa, which makes every value equally likely.
func Float64(r io.Reader) float64 {
	return float64(binary.LittleEndian.Uint64(Bytes(r, 8))>>11) / (1 << 53)
}

// Seed turns a number into the bytes of a ChaCha8 seed.
func Seed(n uint64) [32]byte {
	var seed [32]byte
	binary.LittleEndian.PutUint64(seed[:], n)

	return seed
}

// Stream is a ChaCha8 stream that's safe to read from several goroutines.
type Stream struct {
	mu     sync.Mutex
	chacha *mathrand.ChaCha8
}

// NewStream returns a Stream whose bytes depend only on seed.
func NewStream(seed [32]byte) *Stream {
	return &Stream{chacha: mathrand.NewChaCha8(seed)}
}

// Read fills p and never fails.
func (s *Stream) Read(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.chacha.Read(p)
}
