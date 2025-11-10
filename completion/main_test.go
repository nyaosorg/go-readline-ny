package completion_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/completion"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-ttyadapter/auto"
)

func TestComplete(t *testing.T) {
	editor := &readline.Editor{
		Tty:          &auto.Pilot{Text: []string{keys.CtrlE, "\t", "\r"}},
		Writer:       io.Discard,
		PromptWriter: func(w io.Writer) (int, error) { return 0, nil },
		Default:      "rollback ",
	}

	editor.BindKey(keys.CtrlI, &completion.CmdCompletionOrList2{
		Enclosure: `"'`,
		Delimiter: ",",
		Postfix:   " ",
		Candidates: func(field []string) ([]string, []string) {
			L := len(field)
			var r []string
			if L >= 2 && strings.EqualFold(field[L-2], "ROLLBACK") {
				r = []string{
					"transaction",
					"to",
				}
			} else {
				r = []string{}
			}
			return r, r
		},
	})

	result, err := editor.ReadLine(context.Background())
	if err != nil {
		t.Fatal(err.Error())
	}
	if expect := "rollback t"; result != expect {
		t.Fatalf("expect %#v, but %#v", expect, result)
	}
}
