package dummyin

import (
	"context"
	"testing"

	"github.com/nyaosorg/go-readline-ny"
)

func TestDummyIn(t *testing.T) {
	editor := &readline.Editor{
		OpenKeyGetter: New("a", "i", "u", "\b", "\x1B[D", "e"),
	}
	text, err := editor.ReadLine(context.Background())
	if err != nil {
		t.Fatalf("ERR=%s\n", err.Error())
		return
	}
	if text != "aei" {
		t.Fatalf("TEXT=%s\n", text)
	}
}
