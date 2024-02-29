package readline

import (
	"bytes"
	"testing"
)

func TestSliceInsert(t *testing.T) {
	slice := []byte{1, 2, 3, 4, 5}

	slice = sliceInsert(slice, 3, 6, 7, 8)
	expect := []byte{1, 2, 3, 6, 7, 8, 4, 5}

	if !bytes.Equal(slice, expect) {
		t.Fatalf("expect %v, but %v", expect, slice)
	}
}

func TestSliceDelete(t *testing.T) {
	slice := []byte{1, 2, 3, 4, 5}
	slice = sliceDelete(slice, 2, 2)
	expect := []byte{1, 2, 5}

	if !bytes.Equal(slice, expect) {
		t.Fatalf("expect %v, but %v", expect, slice)
	}
}
