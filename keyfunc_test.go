package readline

import (
	"testing"
)

func TestInsertSelfZeroWithJoiner(t *testing.T) {
	pending := []Moji{RawCodePoint(emojiMan), RawCodePoint(zeroWidthJoinRune)}
	keys := moji2string(pending) + string(emojiRice)
	mojis := string2moji(keys)
	if result := len(mojis); result != 1 {
		t.Fatalf("len(MANFARMER) == %d (expect: 1)", result)
	}
	if result := mojis[0].Width(); result != 5 {
		t.Fatalf("MANFARMER.Width() == %d (expect: 5)", result)
	}
	if result := GetStringWidth(keys); result != 5 {
		t.Fatalf("GetStringWidth(MANFARMER) == %d (expect: 5)", result)
	}
}
