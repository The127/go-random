package randomtest_test

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/The127/go-random"
	"github.com/The127/go-random/randomtest"
)

var _ random.Service = randomtest.NewMock()

func TestNewMockGivesTheSameValuesInEveryRun(t *testing.T) {
	t.Parallel()

	// arrange
	first, second := randomtest.NewMock(), randomtest.NewMock()

	// act
	got := []any{first.Bytes(16), first.Text(), first.IntN(1000), first.Float64()}
	want := []any{second.Bytes(16), second.Text(), second.IntN(1000), second.Float64()}

	// assert
	for i := range got {
		if !equal(got[i], want[i]) {
			t.Errorf("value %d = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestNewMockFromSeedGivesOtherBytesForAnotherSeed(t *testing.T) {
	t.Parallel()

	// arrange
	first, second := randomtest.NewMockFromSeed(1), randomtest.NewMockFromSeed(2)

	// act
	a, b := first.Bytes(16), second.Bytes(16)

	// assert
	if bytes.Equal(a, b) {
		t.Errorf("both seeds gave %x", a)
	}
}

func TestBytesReturnsTheQueuedBytesBeforeTheStream(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()
	mock.QueueBytes([]byte{1, 2, 3})

	// act
	queued, next := mock.Bytes(3), mock.Bytes(3)

	// assert
	if !bytes.Equal(queued, []byte{1, 2, 3}) {
		t.Errorf("first Bytes(3) = %v, want [1 2 3]", queued)
	}

	if want := randomtest.NewMock().Bytes(3); !bytes.Equal(next, want) {
		t.Errorf("second Bytes(3) = %v, want the stream's first bytes %v", next, want)
	}
}

func TestBytesPanicsWhenTheQueuedBytesHaveAnotherLength(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()
	mock.QueueBytes([]byte{1, 2, 3})

	// act and assert
	assertPanics(t, func() { mock.Bytes(4) })
}

func TestTextReturnsTheQueuedTextBeforeTheStream(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()
	mock.QueueText("abc")

	// act
	queued, next := mock.Text(), mock.Text()

	// assert
	if queued != "abc" {
		t.Errorf("first Text() = %q, want %q", queued, "abc")
	}

	if want := randomtest.NewMock().Text(); next != want {
		t.Errorf("second Text() = %q, want the stream's first text %q", next, want)
	}
}

func TestTextHas26CharactersOfTheBase32Alphabet(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()

	// act
	text := mock.Text()

	// assert
	if len(text) != 26 || strings.Trim(text, "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567") != "" {
		t.Errorf("Text() = %q, want 26 characters of A to Z and 2 to 7", text)
	}
}

func TestIntNReturnsTheQueuedNumberBeforeTheStream(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()
	mock.QueueIntN(2)

	// act
	queued, next := mock.IntN(3), mock.IntN(1000)

	// assert
	if queued != 2 {
		t.Errorf("first IntN(3) = %d, want 2", queued)
	}

	if want := randomtest.NewMock().IntN(1000); next != want {
		t.Errorf("second IntN(1000) = %d, want the stream's first number %d", next, want)
	}
}

func TestIntNPanicsWhenTheQueuedNumberIsntBelowN(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()
	mock.QueueIntN(3)

	// act and assert
	assertPanics(t, func() { mock.IntN(3) })
}

func TestQueueIntNPanicsForANegativeNumber(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()

	// act and assert
	assertPanics(t, func() { mock.QueueIntN(-1) })
}

func TestIntNPanicsForNBelowOne(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()

	// act and assert
	assertPanics(t, func() { mock.IntN(0) })
}

func TestIntNReturnsEveryNumberBelowNAndNoOther(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()
	seen := map[int]bool{}

	// act
	for range 1000 {
		seen[mock.IntN(3)] = true
	}

	// assert
	if len(seen) != 3 || !seen[0] || !seen[1] || !seen[2] {
		t.Errorf("IntN(3) returned %v, want exactly 0, 1 and 2", seen)
	}
}

func TestFloat64ReturnsTheQueuedNumberBeforeTheStream(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()
	mock.QueueFloat64(0.25)

	// act
	queued, next := mock.Float64(), mock.Float64()

	// assert
	if queued != 0.25 {
		t.Errorf("first Float64() = %v, want 0.25", queued)
	}

	if want := randomtest.NewMock().Float64(); next != want {
		t.Errorf("second Float64() = %v, want the stream's first number %v", next, want)
	}
}

func TestQueueFloat64PanicsOutsideZeroToOne(t *testing.T) {
	t.Parallel()

	for _, value := range []float64{-0.1, 1, 1.5} {
		// arrange
		mock := randomtest.NewMock()

		// act and assert
		assertPanics(t, func() { mock.QueueFloat64(value) })
	}
}

func TestFloat64StaysAtOrAboveZeroAndBelowOne(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()

	for range 1000 {
		// act
		value := mock.Float64()

		// assert
		if value < 0 || value >= 1 {
			t.Fatalf("Float64() = %v, want 0 <= value < 1", value)
		}
	}
}

func TestReaderReadsTheQueuedBytesBeforeTheStream(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()
	mock.QueueBytes([]byte{1, 2})

	want := append([]byte{1, 2}, randomtest.NewMock().Bytes(3)...)
	got := make([]byte, 5)

	// act
	_, err := io.ReadFull(mock.Reader(), got)
	if err != nil {
		t.Fatal(err)
	}

	// assert
	if !bytes.Equal(got, want) {
		t.Errorf("Reader() read %v, want %v", got, want)
	}
}

func TestMockCanBeUsedFromSeveralGoroutines(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()

	var wg sync.WaitGroup

	// act
	wg.Go(func() { mock.QueueText("abc") })
	wg.Go(func() { _ = mock.Text() })
	wg.Go(func() { _ = mock.Bytes(8) })
	wg.Wait()

	// assert
	// The race detector fails the test if the calls race.
}

func assertPanics(t *testing.T, f func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Error("the call didn't panic")
		}
	}()

	f()
}

func equal(a, b any) bool {
	if x, ok := a.([]byte); ok {
		return bytes.Equal(x, b.([]byte))
	}

	return a == b
}
