package readline

import (
	"strings"
	"testing"
)

func TestColorSequenceAdd(t *testing.T) {
	c := ColorSequence(0)
	c = c.Add(1)
	if c != 0x0101 {
		t.Fatalf("Add(1) failed(%x)", int64(c))
	}
	c = c.Add(2)
	if c != 0x020102 {
		t.Fatalf("Add(2) failed(%x)", int64(c))
	}
	c = c.Chain(c)
	if c != 0x0201020104 {
		t.Fatalf("Chain() failed(%x)", int64(c))
	}
}

func TestDarkDefaultForeGroundColor(t *testing.T) {
	var buffer strings.Builder
	DarkDefaultForeGroundColor.WriteTo(&buffer)
	result := buffer.String()
	expect := "\x1B[39;22;49m"

	if result != expect {
		t.Fatalf("Expect ESC %s, but ESC %s", expect[1:], result[1:])
	}
}
