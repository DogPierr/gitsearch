//go:build !solution

package main

import (
	"os"

	"github.com/DogPierr/gitsearch/gitsearch/internal/gitsearch"
)

func main() {
	os.Exit(gitsearch.Init())
}
