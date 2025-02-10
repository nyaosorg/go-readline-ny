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

type CustomComp struct {
	Candidates []string
}

func (c *CustomComp) Enclosures() string { return `"` }
func (c *CustomComp) Delimiters() string { return "&|><;" }

func (c *CustomComp) List(field []string) (completionList []string, listingList []string) {
	if len(field) <= 1 {
		return c.Candidates, c.Candidates
	}
	return []string{}, []string{}
}

func mains() error {
	var editor readline.Editor

	editor.PromptWriter = func(w io.Writer) (int, error) {
		return io.WriteString(w, "menu> ")
	}
	editor.BindKey(keys.CtrlI, completion.CmdCompletionOrList{
		Completion: &CustomComp{
			Candidates: []string{
				"list", "say", "pewpew", "help", "exit", "Space Command",
			},
		},
		Postfix: " ",
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
