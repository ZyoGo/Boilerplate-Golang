package main

import (
	"log"

	"github.com/ZyoGo/default-ddd-http/cmd/api-http/modules"
)

func main() {
	c, err := modules.New()
	if err != nil {
		log.Fatal("failed to instantiate server: ", err)
	}

	if err = c.Run(); err != nil {
		log.Fatal("error when running server: ", err)
	}
}
