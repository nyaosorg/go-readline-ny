//go:build run

package main

import (
	"context"
	"fmt"

	"github.com/nyaosorg/go-readline-ny"
)

func main() {
	var editor readline.Editor
	text, err := editor.ReadLine(context.Background())
	if err != nil {
		fmt.Printf("ERR=%s\n", err.Error())
	} else {
		fmt.Printf("TEXT=%s\n", text)
	}
}
