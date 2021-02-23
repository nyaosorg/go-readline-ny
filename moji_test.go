package readline

import (
	"testing"
)

const emojiManFarmer = "\U0001F468\u200D\U0001F33E"

func TestZeroWidthJoinSequenceWidth(t *testing.T) {
	manFarmer := ZeroWidthJoinSequence(emojiManFarmer)

	if w := manFarmer.Width(); w != 5 {
		t.Fatalf("EmojiManFarmer's width is invalid (%d). It should be 5", w)
		return
	}
}
