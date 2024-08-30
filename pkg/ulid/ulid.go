package ulid

import (
	"crypto/rand"
	"log"
	"math"
	"math/big"
	mrand "math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

// ULIDGenerator defines the interface for generating ULIDs
type ULIDGenerator interface {
	Generate() string
}

// ulidGenerator is the concrete implementation of ULIDGenerator
type ulidGenerator struct {
	entropyPool sync.Pool
}

// NewULIDGenerator initializes a new ULID generator
func NewGenerator() ULIDGenerator {
	return &ulidGenerator{
		entropyPool: sync.Pool{
			New: func() interface{} {
				seed, err := cryptoRandSeed()
				if err != nil {
					log.Fatalf("Failed to seed random number generator: %v", err)
				}
				source := mrand.New(mrand.NewSource(seed))
				// Use ulid.Monotonic to ensure monotonicity within the same millisecond
				return ulid.Monotonic(source, 0)
			},
		},
	}
}

// cryptoRandSeed generates a cryptographically secure random seed
func cryptoRandSeed() (int64, error) {
	max := big.NewInt(math.MaxInt64)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return n.Int64(), nil
}

// Generate creates a new ULID string
func (g *ulidGenerator) Generate() string {
	t := time.Now().UTC()

	// Retrieve a monotonic entropy source from the pool
	entropy := g.entropyPool.Get().(*ulid.MonotonicEntropy)

	// Generate a new ULID
	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	// Put the entropy source back in the pool
	g.entropyPool.Put(entropy)

	return id.String()
}
