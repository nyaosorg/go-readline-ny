package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/coloring"
	"github.com/nyaosorg/go-readline-ny/tty10"
)

func mains() error {
	logWriter, err := os.Create("output.log")
	if err != nil {
		return err
	}
	defer logWriter.Close()

	editor := readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) {
			return io.WriteString(w, "  0123456789ABCDEF\n$ ")
		},
		Coloring: &coloring.VimBatch{},
		Writer: io.MultiWriter(
			colorable.NewColorableStdout(),
			logWriter),
		Tty: &tty10.Tty{},
	}
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
