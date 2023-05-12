//go:build run
// +build run

package main

import (
	"context"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/completion"
	"github.com/nyaosorg/go-readline-ny/keys"
)

func mains() error {
	var editor readline.Editor

	editor.BindKey(keys.CtrlI, completion.CmdCompletionOrList{
		Completion: completion.File{},
	})
	// editor.BindKey("\t", completion.CmdCompletion{
	//	   Completion: completion.File{},
	// })

	text, err := editor.ReadLine(context.Background())
	if err != nil {
		return err
	}
	println("TEXT:", text)
	return nil
}

func main() {
	if err := mains(); err != nil {
		println("ERR:", err.Error())
	}
}
