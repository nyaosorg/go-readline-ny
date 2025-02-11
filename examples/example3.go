//go:build run

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/completion"
	"github.com/nyaosorg/go-readline-ny/keys"
)

func mains() error {
	var editor readline.Editor

	editor.PromptWriter = func(w io.Writer) (int, error) {
		return io.WriteString(w, "menu> ")
	}
	editor.BindKey(keys.CtrlI, &completion.CmdCompletionOrList2{
		Delimiter: "&|><;",
		Enclosure: `"'`,
		Postfix:   " ",
		Candidates: func(field []string) (forComp []string, forList []string) {
			if len(field) <= 1 {
				c := []string{
					"list", "say", "pewpew", "help", "exit", "Space Command",
				}
				return c, c
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
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		cmd := parts[0]
		switch cmd {
		case "list":
			fmt.Println("try to list")
		case "say":
			fmt.Println("try to say")
		case "pewpew":
			fmt.Println("try to pewpew")
		case "help":
			fmt.Println("try to help")
		case "exit":
			fmt.Println("try to exit")
		default:
			fmt.Println("idk")
		}
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
