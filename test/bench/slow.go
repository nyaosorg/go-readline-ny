//go:build run

package main

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/completion"
	"github.com/nyaosorg/go-readline-ny/keys"
)

type SlowPattern struct{}

func (SlowPattern) FindAllStringIndex(string, int) [][]int {
	time.Sleep(10 * time.Millisecond)
	return nil
}

func main() {
	editor := &readline.Editor{
		Writer: colorable.NewColorableStdout(),
		Highlight: []readline.Highlight{
			{Pattern: regexp.MustCompile("&"), Sequence: "\x1B[33;49;22m"},
			{Pattern: regexp.MustCompile(`"[^"]*"`), Sequence: "\x1B[35;49;22m"},
			{Pattern: regexp.MustCompile(`%[^%]*%`), Sequence: "\x1B[36;49;1m"},
			{Pattern: regexp.MustCompile("\u3000"), Sequence: "\x1B[37;41;22m"},
			{Pattern: SlowPattern{}, Sequence: ""},
		},
		PredictColor: [...]string{"\x1B[3;22;34m", "\x1B[23;39m"},
		ResetColor:   "\x1B[0m",
		DefaultColor: "\x1B[33;49;1m",
	}
	editor.BindKey(keys.CtrlI, completion.CmdCompletionOrList{
		Completion: completion.File{},
		Postfix:    " ",
	})
	text, err := editor.ReadLine(context.Background())

	if err != nil {
		fmt.Printf("ERR=%s\n", err.Error())
		return
	}
	fmt.Printf("TEXT=%s\n", text)
}
