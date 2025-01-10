package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
)

func mains() error {
	logWriter, err := os.Create("output.log")
	if err != nil {
		return err
	}
	defer logWriter.Close()

	if _, ok := os.LookupEnv("ZWJS"); ok {
		// for WindowsTerminal 1.22
		readline.SetZWJSWidthGetter(func(w1, w2 int) int {
			if w2 > w1 {
				return w2
			}
			return w1
		})
	}

	editor := readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) {
			return io.WriteString(w, "  0123456789ABCDEF\n$ ")
		},
		// Coloring: &coloring.VimBatch{},
		Writer: io.MultiWriter(
			colorable.NewColorableStdout(),
			logWriter),
		Highlight: []readline.Highlight{
			{Pattern: readline.ResetColor, EscSeq: "\x1B[0m"},
			{Pattern: readline.DefaultColor, EscSeq: "\x1B[33;49;1m"},
			{Pattern: regexp.MustCompile("&"), EscSeq: "\x1B[33;49;22m"},
			{Pattern: regexp.MustCompile(`"[^"]*"`), EscSeq: "\x1B[31;49;22m"},
			{Pattern: regexp.MustCompile(`%[^%]*%`), EscSeq: "\x1B[36;49;1m"},
			{Pattern: regexp.MustCompile("\u3000"), EscSeq: "\x1B[37;41;22m"},
		},
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
