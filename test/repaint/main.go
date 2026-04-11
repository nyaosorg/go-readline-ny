package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

func mains() error {
	colorable.EnableColorsStdout(nil)

	var w io.Writer = colorable.NewColorableStdout()

	editor := readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) {
			return io.WriteString(w, "String: \"")
		},
		OnAfterRender: func(B *readline.Buffer, availWidth int) {
			if availWidth >= 1 {
				B.Out.Write([]byte{'"', '\b'})
			}
		},
		// Coloring: &coloring.VimBatch{},
		Writer: w,
	}
	editor.BindKey(keys.CtrlL, readline.CmdRepaintLine)
	editor.BindKey(keys.CtrlBackslash, &readline.GoCommand{
		Name: "noise command",
		Func: func(_ context.Context, b *readline.Buffer) readline.Result {
			io.WriteString(b.Out, "noise noise noise noise")
			return readline.CONTINUE
		},
	})
	text, err := editor.ReadLine(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("TEXT=%s\n", text)
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
