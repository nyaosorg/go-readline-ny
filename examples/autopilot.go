//go:build run

package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/nyaosorg/go-ttyadapter/auto"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

func main() {
	editor := &readline.Editor{
		Default: "12345",
		Tty: &auto.Pilot{
			Text: []string{keys.CtrlE, keys.Left, "\b", "\r"},
		},
		Writer:       io.Discard,
		PromptWriter: func(w io.Writer) (int, error) { return 0, nil },
	}

	result, err := editor.ReadLine(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if result == "1235" {
		fmt.Println("PASS")
	} else {
		fmt.Println("FAIL")
	}
}
