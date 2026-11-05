package random_test

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/The127/go-random"
)

func TestEverySourceReturnsTextOf26CharactersOfTheBase32Alphabet(t *testing.T) {
	t.Parallel()

	for name, service := range map[string]random.Service{
		"crypto": random.NewCryptoRandom(),
		"pseudo": random.NewPseudoRandom(),
		"seeded": random.NewPseudoRandomFromSeed(1),
	} {
		// act
		text := service.Text()

		// assert
		if len(text) != 26 || strings.Trim(text, "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567") != "" {
			t.Errorf("%s: Text() = %q, want 26 characters of A to Z and 2 to 7", name, text)
		}
	}
}

func TestPseudoRandomGetsAnotherSeedEachTime(t *testing.T) {
	t.Parallel()

	// arrange
	first, second := random.NewPseudoRandom(), random.NewPseudoRandom()

	// act
	a, b := first.Bytes(16), second.Bytes(16)

	// assert
	if bytes.Equal(a, b) {
		t.Errorf("both pseudo sources gave %x", a)
	}
}

func TestPseudoRandomFromSeedGivesTheSameBytesForTheSameSeed(t *testing.T) {
	t.Parallel()

	// arrange
	first, second := random.NewPseudoRandomFromSeed(7), random.NewPseudoRandomFromSeed(7)

	// act
	a, b := first.Bytes(16), second.Bytes(16)

	// assert
	if !bytes.Equal(a, b) {
		t.Errorf("seed 7 gave %x and %x", a, b)
	}
}

func TestCryptoIntNPanicsForNBelowOne(t *testing.T) {
	t.Parallel()

	// arrange
	service := random.NewCryptoRandom()

	defer func() {
		// assert
		if recover() == nil {
			t.Error("IntN(0) didn't panic")
		}
	}()

	// act
	service.IntN(0)
}

func TestPseudoRandomCanBeUsedFromSeveralGoroutines(t *testing.T) {
	t.Parallel()

	// arrange
	service := random.NewPseudoRandom()

	var wg sync.WaitGroup

	// act
	wg.Go(func() { _ = service.Bytes(8) })
	wg.Go(func() { _ = service.Text() })
	wg.Wait()

	// assert
	// The race detector fails the test if the calls race.
}
