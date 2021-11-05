package main

import (
	"log"

	"github.com/speedata/xts/core"
)

func main() {

	err := core.Dothings()
	if err != nil {
		log.Fatal(err)
	}
}
