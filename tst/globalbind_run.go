// +build run

package main

import (
	"context"
	"fmt"

	"github.com/mattn/go-colorable"
	"github.com/zetamatta/go-readline-ny"
)

func main() {
	editor := readline.Editor{
		Writer: colorable.NewColorableStdout(),
	}
	readline.GlobalKeyMap.BindKeyClosure(readline.K_F1, func(_ context.Context, buffer *readline.Buffer) readline.Result {
		buffer.Writer.Write([]byte{'\x07'})
		return readline.CONTINUE
	})

	text, err := editor.ReadLine(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(text)
}
