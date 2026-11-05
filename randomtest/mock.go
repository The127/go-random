// Package randomtest has a mock random.Service for tests. The mock panics
// outside a test binary, which keeps predictable randomness out of production
// code.
package randomtest

import (
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/The127/go-random/internal/derive"
)

// Mock is a random.Service that returns queued values first and then values
// from a stream that its seed decides. It's safe to use from several
// goroutines.
type Mock struct {
	mu     sync.Mutex
	stream *derive.Stream
	bytes  []byte
	sizes  []int
	texts  []string
	ints   []int
	floats []float64
}

// NewMock returns a Mock with a fixed seed. It returns the same values in
// every run.
func NewMock() *Mock {
	return NewMockFromSeed(0)
}

// NewMockFromSeed returns a Mock whose stream depends only on seed.
func NewMockFromSeed(seed uint64) *Mock {
	if !testing.Testing() {
		panic("randomtest: a mock outside a test binary would make randomness predictable")
	}

	return &Mock{stream: derive.NewStream(derive.Seed(seed))}
}

// QueueBytes makes the next call to Bytes return b. That call must ask for
// len(b) bytes. Reader reads queued bytes too.
func (m *Mock) QueueBytes(b []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.bytes = append(m.bytes, b...)
	m.sizes = append(m.sizes, len(b))
}

// QueueText makes the next call to Text return s.
func (m *Mock) QueueText(s string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.texts = append(m.texts, s)
}

// QueueIntN makes the next call to IntN return v. That call's n must be
// above v. It panics if v is negative.
func (m *Mock) QueueIntN(v int) {
	if v < 0 {
		panic(fmt.Sprintf("randomtest: QueueIntN(%d) needs a number of at least 0", v))
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.ints = append(m.ints, v)
}

// QueueFloat64 makes the next call to Float64 return v. It panics if v isn't
// from 0 to below 1.
func (m *Mock) QueueFloat64(v float64) {
	if v < 0 || v >= 1 {
		panic(fmt.Sprintf("randomtest: QueueFloat64(%v) needs a number from 0 to below 1", v))
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.floats = append(m.floats, v)
}

// Bytes returns the queued bytes, or n bytes from the stream. It panics if
// the queued bytes aren't n bytes.
func (m *Mock) Bytes(n int) []byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.sizes) == 0 {
		return derive.Bytes(m.stream, n)
	}

	if m.sizes[0] != n {
		panic(fmt.Sprintf("randomtest: Bytes(%d) got %d queued bytes", n, m.sizes[0]))
	}

	b := m.bytes[:n:n]
	m.bytes, m.sizes = m.bytes[n:], m.sizes[1:]

	return b
}

// Text returns the queued text, or text from the stream.
func (m *Mock) Text() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.texts) == 0 {
		return derive.Text(m.stream)
	}

	s := m.texts[0]
	m.texts = m.texts[1:]

	return s
}

// IntN returns the queued number, or a number from the stream. It panics if
// n is below 1 or the queued number isn't below n.
func (m *Mock) IntN(n int) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.ints) == 0 {
		return derive.IntN(m.stream, n)
	}

	v := m.ints[0]
	if v >= n {
		panic(fmt.Sprintf("randomtest: IntN(%d) got the queued number %d, which isn't below it", n, v))
	}

	m.ints = m.ints[1:]

	return v
}

// Float64 returns the queued number, or a number from the stream.
func (m *Mock) Float64() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.floats) == 0 {
		return derive.Float64(m.stream)
	}

	v := m.floats[0]
	m.floats = m.floats[1:]

	return v
}

// Reader returns a reader that reads the queued bytes first and then the
// stream.
func (m *Mock) Reader() io.Reader {
	return reader{m}
}

type reader struct{ m *Mock }

func (r reader) Read(p []byte) (int, error) {
	r.m.mu.Lock()
	defer r.m.mu.Unlock()

	n := copy(p, r.m.bytes)
	r.m.bytes = r.m.bytes[n:]
	r.m.dropReadSizes(n)

	if n < len(p) {
		_, _ = r.m.stream.Read(p[n:])
	}

	return len(p), nil
}

// dropReadSizes removes the sizes of the queued slices that Reader has read.
// Bytes then no longer expects them.
func (m *Mock) dropReadSizes(n int) {
	for n > 0 && len(m.sizes) > 0 {
		if m.sizes[0] > n {
			m.sizes[0] -= n

			return
		}

		n -= m.sizes[0]
		m.sizes = m.sizes[1:]
	}
}
