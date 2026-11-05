# go-random

go-random gives Go code cryptographically secure randomness through an interface. A test can then set the values instead of getting new ones each run.

Code that calls `crypto/rand` directly is hard to test. A test can't check the token a function creates, because the token is different every time. With go-random, the code asks a `random.Service` for its randomness. In production that's `crypto/rand`. In a test it's a mock whose values you set.

## Install

```sh
go get github.com/The127/go-random
```

## Use it

Pass the service to the code that needs randomness:

```go
type Tokens struct {
	random random.Service
}

func NewTokens(r random.Service) *Tokens {
	return &Tokens{random: r}
}

func (t *Tokens) New() string {
	return base64.RawURLEncoding.EncodeToString(t.random.Bytes(32))
}
```

In production code, give it `crypto/rand`:

```go
tokens := NewTokens(random.NewCryptoRandom())
```

For values that aren't secret, for example to spread out retries, `random.NewPseudoRandom()` is faster. It uses a ChaCha8 stream from `math/rand/v2` and takes its seed from `crypto/rand`, which makes it differ on every start. `random.NewPseudoRandomFromSeed(seed)` gives the same values in every run, for a simulation you want to repeat. Never use either one for tokens, salts or keys.

In a test, give it a mock from `randomtest`:

```go
func TestNewToken(t *testing.T) {
	mock := randomtest.NewMock()
	mock.QueueBytes(make([]byte, 32))
	tokens := NewTokens(mock)

	if got := tokens.New(); got != "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" {
		t.Errorf("New() = %q", got)
	}
}
```

A `random.Service` has five methods. `Bytes(n)` returns n random bytes, `Text()` a text of 26 characters, `IntN(n)` a number below n and `Float64()` a number from 0 to below 1. `Reader()` returns an `io.Reader` for libraries that take one, such as `ecdsa.GenerateKey`.

## The mock

`randomtest.NewMock()` returns values from a stream with a fixed seed. A test then gets the same values in every run. `randomtest.NewMockFromSeed(seed)` takes another seed.

To set a value, queue it. `QueueBytes`, `QueueText`, `QueueIntN` and `QueueFloat64` decide what the next call returns, and after that the mock returns values from its stream again. If a queued value doesn't fit the call, for example 3 queued bytes for `Bytes(4)`, the mock panics.

The mock also panics outside a test binary. That keeps predictable randomness out of production code. You can use the mock from several goroutines at once.

go-random has no dependencies outside Go's standard library.

## Contribute

[CONTRIBUTING.md](CONTRIBUTING.md) explains how to set up the tools, how to sign off your commits and how releases work.

## License

go-random is licensed under the [MIT License](LICENSE.md).
