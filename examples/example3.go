package main

import (
	"context"
	"fmt"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-readline-ny/simplehistory"
)

func main() {
	var editor readline.Editor

	history := simplehistory.New()
	editor.History = history
	editor.BindKey(keys.CtrlS, readline.CmdBackwardChar)
	editor.BindKey(keys.CtrlD, readline.CmdForwardChar)
	editor.BindKey(keys.CtrlX, readline.CmdNextHistory)
	editor.BindKey(keys.CtrlE, readline.CmdPreviousHistory)
	editor.BindKey(keys.CtrlA, readline.CmdBackwardWord)
	editor.BindKey(keys.CtrlF, readline.CmdForwardWord)
	editor.BindKey(keys.CtrlG, readline.CmdDeleteChar)

	disableKeys := []keys.Code{
		keys.CtrlB,
		keys.CtrlP,
		keys.CtrlN,
	}
	for _, key := range disableKeys {
		editor.BindKey(key, readline.SelfInserter(key))
	}

	for {
		text, err := editor.ReadLine(context.Background())
		if err != nil {
			fmt.Printf("ERR=%s\n", err.Error())
			return
		}
		fmt.Printf("TEXT=%s\n", text)

		history.Add(text)
	}
}
