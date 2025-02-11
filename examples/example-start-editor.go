//go:build run
// +build run

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

func startEditor(source string) (string, error) {
	textEditor, ok := os.LookupEnv("EDITOR")
	if !ok {
		return source, errors.New("$EDITOR is not defined")
	}
	fd, err := os.CreateTemp("", "example3")
	if err != nil {
		return source, err
	}
	io.WriteString(fd, source)
	fd.Close()
	fname := fd.Name()
	defer os.Remove(fname)

	cmd := exec.Command(textEditor, fname)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return source, err
	}
	update, err := os.ReadFile(fname)
	if err != nil {
		return source, err
	}
	update = bytes.TrimSuffix(update, []byte{'\n'})
	update = bytes.TrimSuffix(update, []byte{'\r'})
	return string(update), nil
}

func cmdStartEditor(ctx context.Context, B *readline.Buffer) readline.Result {
	result, err := startEditor(B.String())
	if err != nil {
		return readline.CONTINUE
	}
	B.Buffer = B.Buffer[:0]
	B.InsertString(0, result)
	B.Cursor = len(B.Buffer)
	B.RepaintAll()
	return readline.CONTINUE
}

func main() {
	var editor readline.Editor
	editor.Writer = colorable.NewColorableStdout()
	editor.PromptWriter = func(w io.Writer) (int, error) {
		return io.WriteString(w, "\r> ")
	}

	enter_status := 0
	editor.BindKey(keys.Escape, &readline.GoCommand{
		Name: "START_EDITOR",
		Func: cmdStartEditor,
	})

	text, err := editor.ReadLine(context.Background())
	fmt.Printf("ENTER_STATUS=%d\n", enter_status)
	if err != nil {
		fmt.Printf("ERR=%s\n", err.Error())
	} else {
		fmt.Printf("TEXT=%s\n", text)
	}
}
