//go:build run

package main

import (
	"context"
	"fmt"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/simplehistory"
)

func main() {
	var editor readline.Editor
	h := simplehistory.New()
	editor.History = h
	for {
		text, err := editor.ReadLine(context.Background())
		if err != nil {
			fmt.Printf("ERR=%s\n", err.Error())
			return
		}
		fmt.Printf("TEXT=%s\n", text)
		h.Add(text)

		// editor.History.Add(text)
		// -> compile error: History interface does not have Add method
		//    But the value returned by `simplehistory.New()` does.
	}
}
