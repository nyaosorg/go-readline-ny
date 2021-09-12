package main

import (
	"context"
	"fmt"

	"github.com/nyaosorg/go-readline-ny"
)

func main() {
	editor := readline.Editor{
		Prompt: func() (int, error) {
			print("  0123456789ABCDEF\n$ ")
			return 2, nil
		},
	}
	text, err := editor.ReadLine(context.Background())
	if err != nil {
		fmt.Printf("ERR=%s\n", err.Error())
	} else {
		fmt.Printf("TEXT=%s\n", text)
	}
}
