package main

import (
	"context"
	"fmt"

	"github.com/nyaosorg/go-readline-ny"
)

type Coloring struct {
	quoted bool
}

func (c *Coloring) Init() int {
	c.quoted = false
	return 0
}

func (c *Coloring) Next(r rune) int {
	next := c.quoted

	if r == '"' {
		next = !next
	}

	defer func() {
		c.quoted = next
	}()

	if next || c.quoted {
		return readline.Magenta
	}
	return 0
}

func main() {
	var editor readline.Editor

	editor.Coloring = &Coloring{}
	text, err := editor.ReadLine(context.Background())
	if err != nil {
		fmt.Printf("ERR=%s\n", err.Error())
	} else {
		fmt.Printf("TEXT=%s\n", text)
	}
}
