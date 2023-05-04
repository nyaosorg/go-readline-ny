//go:build run
// +build run

package main

import (
	"context"
	"fmt"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

func main() {
	editor := readline.Editor{
		Default: "AHAHA",
		Cursor:  3,
		Writer:  colorable.NewColorableStdout(),
	}

	enter_status := 0
	editor.BindKey(keys.CtrlP, readline.AnonymousCommand(
		func(ctx context.Context, r *readline.Buffer) readline.Result {
			enter_status = -1
			return readline.ENTER
		}))

	editor.BindKey(keys.CtrlN, readline.AnonymousCommand(
		func(ctx context.Context, r *readline.Buffer) readline.Result {
			enter_status = +1
			return readline.ENTER
		}))

	text, err := editor.ReadLine(context.Background())
	fmt.Printf("ENTER_STATUS=%d\n", enter_status)
	if err != nil {
		fmt.Printf("ERR=%s\n", err.Error())
	} else {
		fmt.Printf("TEXT=%s\n", text)
	}
}
