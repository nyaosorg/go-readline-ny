package readline

import (
	"strings"
	"testing"
)

const emojiMan = '\U0001F468'
const emojiRice = '\U0001F33E'
const emojiManFarmer = string(emojiMan) + string(zeroWidthJoinRune) + string(emojiRice)

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
		c.PrintTo(&buffer)
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

func TestString2Moji(t *testing.T) {
	SurrogatePairOk = true
	ZeroWidthJoinSequenceOk = true
	VariationSequenceOk = true

	var table = []struct {
		Source string
		Title  string
		Count  int
		Width  WidthT
	}{
		{
			Source: "\U0001F926\u200D\u2640\uFE0F",
			Title:  "WOMAN FACEPALMING",
			Count:  1,
			Width:  5,
		},
		{
			Source: "#\uFE0F\u20E3",
			Title:  "EnclosedNumberSign",
			Count:  1,
			Width:  3,
		},
		{
			Source: "\U0001F3F3\uFE0F",
			Title:  "WhiteFlag",
			Count:  1,
			Width:  2,
		},
	}

	for _, p := range table {
		t.Logf("try %s", p.Title)
		mojis := string2moji(p.Source)

		if result := len(mojis); result != p.Count {
			t.Fatalf("Count of %s == %d (expect %d)",
				p.Title, result, p.Count)
		}
		if result := mojis[0].Width(); result != p.Width {
			t.Fatalf("Width of %s == %d (expect %d)",
				p.Title, result, p.Width)
		}
	}
}
