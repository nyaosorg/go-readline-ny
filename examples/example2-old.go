//go:build run

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/completion"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-readline-ny/simplehistory"
)

type OSClipboard struct{}

func (OSClipboard) Read() (string, error) {
	return clipboard.ReadAll()
}

func (OSClipboard) Write(s string) error {
	return clipboard.WriteAll(s)
}

func main() {
	history := simplehistory.New()

	editor := &readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) {
			return io.WriteString(w, "\x1B[36;22m$ ") // print `$ ` with cyan
		},
		Writer:  colorable.NewColorableStdout(),
		History: history,
		Highlight: []readline.Highlight{
			{Pattern: regexp.MustCompile("&"), Sequence: "\x1B[33;49;22m"},
			{Pattern: regexp.MustCompile(`"[^"]*"`), Sequence: "\x1B[35;49;22m"},
			{Pattern: regexp.MustCompile(`%[^%]*%`), Sequence: "\x1B[36;49;1m"},
			{Pattern: regexp.MustCompile("\u3000"), Sequence: "\x1B[37;41;22m"},
		},
		HistoryCycling: true,
		PredictColor:   [...]string{"\x1B[3;22;34m", "\x1B[23;39m"},
		ResetColor:     "\x1B[0m",
		DefaultColor:   "\x1B[33;49;1m",

		Clipboard: OSClipboard{},
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
