package ulid_test

import (
	"sync"
	"testing"

	"github.com/ZyoGo/default-ddd-http/pkg/ulid"
)

func TestULIDGenerator_ConcurrentGeneration(t *testing.T) {
	const numGoroutines = 1000
	generator := ulid.NewGenerator()

	// Use a map to store generated ULIDs and a mutex for concurrent writes
	ulidMap := make(map[string]struct{})
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()

			id := generator.Generate()

			// Lock map access to prevent race conditions
			mu.Lock()
			defer mu.Unlock()

			if _, exists := ulidMap[id]; exists {
				t.Errorf("Duplicate ULID generated: %s", id)
			}

			ulidMap[id] = struct{}{}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()
}
