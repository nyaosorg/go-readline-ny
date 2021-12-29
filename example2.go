//go:build run
// +build run

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/coloring"
	"github.com/nyaosorg/go-readline-ny/simplehistory"
)

func main() {
	history := simplehistory.New()

	editor := readline.Editor{
		Prompt:         func() (int, error) { return fmt.Print("$ ") },
		Writer:         colorable.NewColorableStdout(),
		History:        history,
		Coloring:       &coloring.VimBatch{},
		HistoryCycling: true,
	}
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
