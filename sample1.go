// +build run

// go run sample1.go

package main

import (
	"context"
	"fmt"

	"github.com/mattn/go-colorable"

	"github.com/zetamatta/go-readline-ny"
)

func main() {
	editor := readline.Editor{
		Default: "InitialValue",
		Cursor:  3,
		Writer:  colorable.NewColorableStdout(),
	}
	text, err := editor.ReadLine(context.Background())

	if err != nil {
		fmt.Printf("ERR=%s\n", err.Error())
	} else {
		fmt.Printf("TEXT=%s\n", text)
	}
}
