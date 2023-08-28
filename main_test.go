package readline

import (
	"testing"
)

func TestCutEscapeSequenceAndOldLine(t *testing.T) {
	// backspace test for the issue: https://github.com/nyaosorg/nyagos/issues/422
	result := cutEscapeSequenceAndOldLine("2022-08-08 19:12:55.96\b\b\b>")
	expect := "2022-08-08 19:12:55>"
	if result != expect {
		t.Fatalf("expect '%s' but '%s'", expect, result)
	}

	// get the last line only
	result = cutEscapeSequenceAndOldLine("first line\nsecond line")
	expect = "second line"
	if result != expect {
		t.Fatalf("expect '%s' but '%s'", expect, result)
	}
}
