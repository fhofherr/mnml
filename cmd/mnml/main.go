package main

import (
	"fmt"
	"os"

	"github.com/fhofherr/mnml/internal/cmd/mnml"
)

func main() {
	mnml := mnml.New()

	if err := mnml.Execute(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
