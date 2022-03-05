package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	bytes, _ := io.ReadAll(os.Stdin)
	for _, ch := range string(bytes) {
		fmt.Printf("U+%X\n", ch)
	}
}
