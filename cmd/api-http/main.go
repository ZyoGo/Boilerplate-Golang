package main

import (
	"os"
	"sync"

	server "github.com/ZyoGo/default-ddd-http/cmd/api-http/modules"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		os.Exit(server.Run())
		defer wg.Done()
	}()

	wg.Wait()
}
