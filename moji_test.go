package readline

import (
	"strings"
	"testing"
)

const emojiMan = '\U0001F468'
const emojiRice = '\U0001F33E'
const emojiManFarmer = string(emojiMan) + zeroWidthJoinStr + string(emojiRice)

func TestZeroWidthJoinSequenceWidth(t *testing.T) {
	manFarmer := string2moji(emojiManFarmer)

	if w := manFarmer[0].Width(); w != 5 {
		t.Fatalf("EmojiManFarmer's width is invalid (%d). It should be 5", w)
		return
	}
}

func TestCodePointPut(t *testing.T) {
	SurrogatePairOk = false
	ZeroWidthJoinSequenceOk = false
	VariationSequenceOk = false

	source := "1\b\t\x7F\u908A\U000E0104"
	expect := "1^H^I<7F>\u908A<E0104>"
	mojis := string2moji(source)

	var buffer strings.Builder
	for _, c := range mojis {
		c.Put(&buffer)
	}
	result := buffer.String()
	if result != expect {
		t.Fatalf("string2moji(%v)='%s' (expect '%s')", source, result, expect)
	}

	buffer.Reset()
	for _, c := range mojis {
		c.WriteTo(&buffer)
	}
	result = buffer.String()
	if result != source {
		t.Fatalf("string2moji(%v).WriteTo()='%s' (expect '%s')", source, result, source)
	}
}

func TestWomanFacepalming(t *testing.T) {
	SurrogatePairOk = true
	ZeroWidthJoinSequenceOk = true
	VariationSequenceOk = true

	source := "\U0001F926\u200D\u2640\uFE0F"
	mojis := string2moji(source)

	if result := len(mojis); result != 1 {
		t.Fatalf("len(string2moji(WOMAN FACEPALMING)) == %d (expect %d)", result, 1)
	}
	if result := mojis[0].Width(); result != 5 {
		t.Fatalf("Width(WOMAN FACEPALMING)) == %d (expect %d)", result, 4)
	}
}
