package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
)

var flagLog = flag.String("l", "", "output log filename")

func mains() error {
	colorable.EnableColorsStdout(nil)

	var w io.Writer = colorable.NewColorableStdout()
	if *flagLog != "" {
		logWriter, err := os.Create(*flagLog)
		if err != nil {
			return err
		}
		defer logWriter.Close()

		w = io.MultiWriter(w, logWriter)
	}

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
			return io.WriteString(w, "\x1B[30;49;1m  0123456789ABCDEF\n$ ")
		},
		// Coloring: &coloring.VimBatch{},
		Writer: w,
		Highlight: []readline.Highlight{
			{Pattern: regexp.MustCompile("&"), Sequence: "\x1B[33;49;22m"},
			{Pattern: regexp.MustCompile(`"[^"]*"`), Sequence: "\x1B[31;49;22m"},
			{Pattern: regexp.MustCompile(`%[^%]*%`), Sequence: "\x1B[36;49;1m"},
			{Pattern: regexp.MustCompile("\u3000"), Sequence: "\x1B[37;41;22m"},
		},
		ResetColor:   "\x1B[0m",
		DefaultColor: "\x1B[33;49;1m",
	}
	text, err := editor.ReadLine(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("TEXT=%s\n", text)
	return nil
}

func main() {
	flag.Parse()
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
