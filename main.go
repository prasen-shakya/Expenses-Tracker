package main

import (
	"fmt"
	"os"
)

func main() {
	store, err := NewStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open store: %v\n", err)
		os.Exit(1)
	}
	defer store.CloseStore()
}
