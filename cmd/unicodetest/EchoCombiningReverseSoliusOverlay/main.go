package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("[%c", 0x20E5)
	os.Stdout.Sync()
	fmt.Printf(",%c]\n", 0x20E5)
}
