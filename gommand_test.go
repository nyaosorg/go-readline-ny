package readline_test

import (
	"context"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/auto"
)

const f = "\U0001F468\u200D\U0001F33E"

func tryAll(t *testing.T, texts ...string) (string, []string) {
	var buffer strings.Builder
	editor := readline.Editor{
		Tty:          &auto.Pilot{Text: texts},
		Writer:       &buffer,
		PromptWriter: func(w io.Writer) (int, error) { return 0, nil },
	}
	result, err := editor.ReadLine(context.Background())
	if err != nil {
		t.Fatalf("ERR=%s", err.Error())
		return "", nil
	}

	outputPieces := strings.Split(buffer.String(), "\x1B[?25h\x1B[?25l")
	for i, s := range outputPieces {
		s = strings.ReplaceAll(s, "\x1B[?25h", "")
		s = strings.ReplaceAll(s, "\x1B[?25l", "")
		s = strings.ReplaceAll(s, "\u200D", "<ZWJ>")
		s = strings.ReplaceAll(s, "\x1B", "<ESC>")
		outputPieces[i] = s
	}
	return result, outputPieces
}

func TestKeyFuncBackSpace(t *testing.T) {
	result, outputs := tryAll(t, f, "\b", "x")
	expect := "x"
	if runtime.GOOS != "windows" || os.Getenv("WT_SESSION") == "" {
		expect = "\U0001F468\u200Dx"
	}
	if result != expect {
		t.Fatalf("TEXT=%s", result)
		return
	}
	if false {
		for _, o := range outputs {
			println(o)
		}
	}
}
