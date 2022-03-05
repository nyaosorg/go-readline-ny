package readline_test

import (
	"context"
	"strings"
	"testing"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/internal/dummyin"
)

const f = "\U0001F468\u200D\U0001F33E"

func tryAll(t *testing.T, texts ...string) (string, string) {
	var buffer strings.Builder
	editor := readline.Editor{
		OpenKeyGetter: dummyin.New(texts...),
		Writer:        &buffer,
		Prompt:        func() (int, error) { return 0, nil },
	}
	result, err := editor.ReadLine(context.Background())
	if err != nil {
		t.Fatalf("ERR=%s", err.Error())
		return "", buffer.String()
	}
	return result, buffer.String()
}

func TestKeyFuncBackSpace(t *testing.T) {
	result, _ := tryAll(t, f, f, f, f, f, "\b", "\b", "x")
	if result != f+f+f+"x" {
		t.Fatalf("TEXT=%s", result)
		return
	}
}
