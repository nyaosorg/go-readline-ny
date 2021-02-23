package readline

import "testing"

func TestGetCharWidth(t *testing.T) {
	if w := GetCharWidth('\b'); w != 2 { // ^H
		t.Fatalf("GetCharWidth('\\b') == %d : failed", w)
		return
	}
	if w := GetCharWidth(0x7F); w != 4 { // <7F>
		t.Fatalf("GetCharWidth(DEL) == %d : failed", w)
		return
	}
	if w := GetCharWidth('\uFEFF'); w != 6 { // <FEFF>
		t.Fatalf("GetCharWidth('\\uFEFF') == %d : failed", w)
		return
	}
	if w := GetCharWidth('a'); w != 1 { // Alfabet
		t.Fatalf("GetCharWidth('a') == %d : failed", w)
	}
	if w := GetCharWidth('\u85ab'); w != 2 { // Kanji
		t.Fatalf("GetCharWidth('a') == %d : failed", w)
	}
}
