package ulid_test

import (
	"sync"
	"testing"

	"github.com/ZyoGo/default-ddd-http/pkg/ulid"
)

func BenchmarkULIDGenerator_ConcurrentGeneration(b *testing.B) {
	const numGoroutines = 10000
	generator := ulid.NewGenerator()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for j := 0; j < numGoroutines; j++ {
			go func() {
				defer wg.Done()
				generator.Generate()
			}()
		}

		wg.Wait()
	}
}
