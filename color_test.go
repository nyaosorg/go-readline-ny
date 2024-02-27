package readline

import (
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

	c = c.Chain(0)

}
