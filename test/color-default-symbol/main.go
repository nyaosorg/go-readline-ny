package main

import (
	"context"
	"fmt"

	"github.com/nyaosorg/go-readline-ny"
)

type Coloring struct {
	quoted bool
}

var (
	defaultColor = readline.White   // SGR4(0, 37, 40, 1)
	quotedColor  = readline.Magenta //SGR4(0, 35, 40, 1)
)

func (c *Coloring) Init() readline.ColorSequence {
	c.quoted = false
	return defaultColor
}

func (c *Coloring) Next(r rune) readline.ColorSequence {
	next := c.quoted

	if r == '"' {
		next = !next
	}

	defer func() {
		c.quoted = next
	}()

	if next || c.quoted {
		return quotedColor
	}
	return defaultColor
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
