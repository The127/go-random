package random_test

import (
	"encoding/base64"
	"testing"

	"github.com/The127/go-random"
	"github.com/The127/go-random/randomtest"
)

type tokens struct {
	random random.Service
}

func (t *tokens) New() string {
	return base64.RawURLEncoding.EncodeToString(t.random.Bytes(32))
}

func TestTheREADMEExampleWorks(t *testing.T) {
	t.Parallel()

	// arrange
	mock := randomtest.NewMock()
	mock.QueueBytes(make([]byte, 32))
	tokens := &tokens{random: mock}

	// act
	got := tokens.New()

	// assert
	if got != "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" {
		t.Errorf("New() = %q", got)
	}
}
