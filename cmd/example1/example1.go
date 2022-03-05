package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/nyaosorg/go-readline-ny"
)

func main() {
	signal.Ignore(os.Interrupt)

	editor := readline.Editor{}
	text, err := editor.ReadLine(context.Background())
	if err != nil {
		fmt.Printf("ERR=%s\n", err.Error())
	} else {
		fmt.Printf("TEXT=%s\n", text)
	}
}
