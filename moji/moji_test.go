package moji

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

const emojiMan = '\U0001F468'
const emojiRice = '\U0001F33E'
const emojiManFarmer = string(emojiMan) + string(zeroWidthJoinRune) + string(emojiRice)

func TestZeroWidthJoinSequenceWidth(t *testing.T) {
	manFarmer := StringToMoji(emojiManFarmer)

	expect := WidthT(5)
	if runtime.GOOS != "windows" || os.Getenv("WT_SESSION") == "" {
		// <1F468> as EscCodePoint because zeroWidthJoinRune does not work.
		expect = 7
	}
	if result := manFarmer[0].Width(); result != expect {
		t.Fatalf("`%#v`: EmojiManFarmer's width is invalid (%d). It should be %d", manFarmer[0], result, expect)
		return
	}
}

func TestCodePointPut(t *testing.T) {
	SurrogatePairOk = false
	ZeroWidthJoinSequenceOk = false
	VariationSequenceOk = false

	source := "1\b\t\x7F\u908A\U000E0104"
	expect := "1^H <7F>\u908A<E0104>"
	mojis := StringToMoji(source)

	var buffer strings.Builder
	for _, c := range mojis {
		c.PrintTo(&buffer)
	}
	result := buffer.String()
	if result != expect {
		t.Fatalf("StringToMoji(%v)='%s' (expect '%s')", source, result, expect)
	}

	buffer.Reset()
	for _, c := range mojis {
		c.WriteTo(&buffer)
	}
	result = buffer.String()
	if result != source {
		t.Fatalf("StringToMoji(%v).WriteTo()='%s' (expect '%s')", source, result, source)
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
		Count2 int
		Width  WidthT
		Width2 WidthT // not Windows Terminal, e.g. xterm
	}{
		{
			Source: "\U0001F926\u200D\u2640\uFE0F",
			Title:  "WOMAN FACEPALMING",
			Count:  1,
			Count2: 1,
			Width:  5,
			Width2: 6, // equals three code points
		},
		{
			Source: "#\uFE0F\u20E3",
			Title:  "EnclosedNumberSign",
			Count:  1,
			Count2: 1,
			Width:  3,
			Width2: 3,
		},
		{
			Source: "\U0001F3F3\uFE0F",
			Title:  "WhiteFlag",
			Count:  1,
			Count2: 1,
			Width:  2,
			Width2: 2,
		},
		/*{
			Source: "\U0001F647\U0001F3FF",
			Title:  "PersonBowing: dark skin tone",
			Count:  1,
			Count2: 2, // two code points
			Width:  4,
			Width2: 4,
		},
		*/
	}

	for _, p := range table {
		// t.Logf("try %s", p.Title)
		mojis := StringToMoji(p.Source)

		expectCount := p.Count
		expectWidth := p.Width
		if runtime.GOOS != "windows" || os.Getenv("WT_SESSION") == "" {
			expectWidth = p.Width2
			expectCount = p.Count2
		}
		if result := len(mojis); result != expectCount {
			t.Fatalf("StringToMoji: Count of %s == %d (expect %d)",
				p.Title, result, expectCount)
		}
		if result := mojis[0].Width(); result != expectWidth {
			t.Fatalf("StringToMoji: Width of %#v(%s) == %d (expect %d)",
				mojis[0], p.Title, result, expectWidth)
		}
		if w, c := MojiWidthAndCountInString(p.Source); w != expectWidth || c != expectCount {
			t.Fatalf("MojiWidthAndCountInString: Width and Count of %s == %d,%d (expect %d,%d)",
				p.Title, w, c, expectWidth, expectCount)
		}
	}
}
