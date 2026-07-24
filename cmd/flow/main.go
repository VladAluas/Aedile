package main

import (
	"log"

	"github.com/VladAluas/Aedile/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
