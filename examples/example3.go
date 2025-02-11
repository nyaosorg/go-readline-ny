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

	// If you do not want to list files with double-tab-key,
	// use `CmdCompletion2` instead of `CmdCompletionOrList2`

	editor.BindKey(keys.CtrlI, &completion.CmdCompletionOrList2{
		// Characters listed here are excluded from completion.
		Delimiter: "&|><;",
		// Enclose candidates with these characters when they contain spaces
		Enclosure: `"'`,
		// String to append when only one candidate remains
		Postfix: " ",
		// Function for listing candidates
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
