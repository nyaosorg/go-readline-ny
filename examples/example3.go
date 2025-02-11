//go:build run

package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/completion"
	"github.com/nyaosorg/go-readline-ny/keys"
)

func mains() error {
	var editor readline.Editor

	editor.PromptWriter = func(w io.Writer) (int, error) {
		return io.WriteString(w, "menu> ")
	}
	candidates := []string{"list", "say", "pewpew", "help", "exit", "Space Command"}

	editor.BindKey(keys.CtrlI, &completion.CmdCompletionOrList2{
		Delimiter: "&|><;",
		Enclosure: `"'`,
		Postfix:   " ",
		Candidates: func(field []string) (forComp []string, forList []string) {
			if len(field) <= 1 {
				return candidates, candidates
			}
			return nil, nil
		},
	})
	ctx := context.Background()
	for {
		line, err := editor.ReadLine(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("TEXT=%#v\n", line)
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
