//go:build run

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/coloring"
	"github.com/nyaosorg/go-readline-ny/completion"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-readline-ny/simplehistory"
)

func main() {
	history := simplehistory.New()

	editor := &readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) {
			return io.WriteString(w, "\x1B[1;36m$ \x1B[0m") // print `$ ` with cyan
		},
		Writer:         colorable.NewColorableStdout(),
		History:        history,
		Coloring:       &coloring.VimBatch{},
		HistoryCycling: true,
	}

	editor.BindKey(keys.CtrlI, completion.CmdCompletionOrList{
		Completion: completion.File{},
		Postfix:    " ",
	})
	// If you do not want to list files with double-tab-key,
	// use `CmdCompletion` instead of `CmdCompletionOrList`

	fmt.Println("Tiny Shell. Type Ctrl-D to quit.")
	for {
		text, err := editor.ReadLine(context.Background())

		if err != nil {
			fmt.Printf("ERR=%s\n", err.Error())
			return
		}

		fields := strings.Fields(text)
		if len(fields) <= 0 {
			continue
		}
		cmd := exec.Command(fields[0], fields[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		cmd.Run()

		history.Add(text)
	}
}
