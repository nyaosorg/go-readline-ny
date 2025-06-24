package moji

import (
	"io"
	"unicode"
	"unicode/utf8"
)

// Moji is the interface for minimum unit to edit in readline
//
// When we make a new implement type of Moji,
// we have to append the code in the function:
// StringToMoji() and KeyFuncInsertSelf().
type Moji interface {
	Width() WidthT
	WriteTo(io.Writer) (int64, error)
	PrintTo(io.Writer)
	Len() int
}

var GetZWJSWidth = func(w1, w2 int) int {
	return w1 + 1 + w2 // WindowsTerminal until 1.21
}

type ZeroWidthJoinSequence [2]Moji

func (s ZeroWidthJoinSequence) Len() int {
	return s[0].Len() + len(string(zeroWidthJoinRune)) + s[1].Len()
}

func (s ZeroWidthJoinSequence) Width() WidthT {
	// runewidth.StringWidth should not be used because the width that it gives
	// has no compatible with WindowsTerminal's.
	return WidthT(GetZWJSWidth(int(s[0].Width()), int(s[1].Width())))
}

func (s ZeroWidthJoinSequence) WriteTo(w io.Writer) (int64, error) {
	n1, err := s[0].WriteTo(w)
	if err != nil {
		return n1, err
	}
	n2, err := writeRune(w, zeroWidthJoinRune)
	if err != nil {
		return n1 + n2, err
	}
	n3, err := s[1].WriteTo(w)
	return n1 + n2 + n3, err
}

func (s ZeroWidthJoinSequence) PrintTo(w io.Writer) {
	switch s0 := s[0].(type) {
	case WavingWhiteFlagCodePoint:
		saveCursorAfterN(w, s.Width())
		s0.WriteTo(w)
		writeRune(w, zeroWidthJoinRune)
		s[1].WriteTo(w)
		restoreCursor(w)
	default:
		s0.PrintTo(w)
		writeRune(w, zeroWidthJoinRune)
		s[1].PrintTo(w)
	}
}

type ModifierSequence [2]Moji

func isEmojiModifier(ch rune) bool {
	if !ModifierSequenceOk {
		return false
	}
	return '\U0001F3FB' <= ch && ch <= '\U0001F3FF'
}

func AreEmojiModifier(s string) bool {
	if !ModifierSequenceOk {
		return false
	}
	u, _ := utf8.DecodeRuneInString(s)
	return isEmojiModifier(u)
}

func (s ModifierSequence) Len() int {
	return s[0].Len() + s[1].Len()
}

func (s ModifierSequence) Width() WidthT {
	return s[0].Width() + s[1].Width()
}

func (s ModifierSequence) WriteTo(w io.Writer) (int64, error) {
	n1, err := s[0].WriteTo(w)
	if err != nil {
		return n1, err
	}
	n2, err := s[1].WriteTo(w)
	return n1 + n2, err
}

func (s ModifierSequence) PrintTo(w io.Writer) {
	s.WriteTo(w)
}

type VariationSequence [2]Moji

func (s VariationSequence) Len() int {
	return s[0].Len() + s[1].Len()
}

func (s VariationSequence) Width() WidthT {
	switch s0 := s[0].(type) {
	case WavingWhiteFlagCodePoint:
		return s0.Width()
	default:
		return s0.Width() + 1
	}
}

func (s VariationSequence) WriteTo(w io.Writer) (int64, error) {
	n1, err := s[0].WriteTo(w)
	if err != nil {
		return n1, err
	}
	n2, err := s[1].WriteTo(w)
	return n1 + n2, err
}

func saveCursorAfterN(w io.Writer, n WidthT) {
	for i := WidthT(0); i < n; i++ {
		w.Write([]byte{' '})
	}
	w.Write([]byte{'\x1B', '7'})
	for i := WidthT(0); i < n; i++ {
		w.Write([]byte{'\b'})
	}
}

func restoreCursor(w io.Writer) {
	w.Write([]byte{'\x1B', '8'})
}

func (s VariationSequence) PrintTo(w io.Writer) {
	saveCursorAfterN(w, s.Width())
	// The sequence 'ESC 7' can not remember the cursor position more than one.
	// When _VariationSequence contains another _VariationSequence
	// s[0].PrintTo(w) does not work as we expect.
	s[0].WriteTo(w)
	s[1].WriteTo(w)
	restoreCursor(w)
}

const zeroWidthJoinRune = '\u200D'

func isZeroWidthJoin(r rune) bool {
	return ZeroWidthJoinSequenceOk && unicode.Is(unicode.Join_Control, r)
}

func AreZeroWidthJoin(s string) bool {
	if !ZeroWidthJoinSequenceOk {
		return false
	}
	r, _ := utf8.DecodeRuneInString(s)
	return isZeroWidthJoin(r)
}

func StringToMoji(s string) []Moji {
	mojis := make([]Moji, 0, len(s))
	var last Moji
	pos := 0
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		s = s[size:]

		if last != nil {
			if isZeroWidthJoin(r) {
				next, nextsize := utf8.DecodeRuneInString(s)
				s = s[nextsize:]

				last = ZeroWidthJoinSequence([...]Moji{last, RawCodePoint(next)})
				mojis[len(mojis)-1] = last
				continue
			}
			if isVariationSelectorLike(r) {
				last = VariationSequence([...]Moji{last, RawCodePoint(r)})
				mojis[len(mojis)-1] = last
				continue
			}
			if isEmojiModifier(r) {
				last = ModifierSequence([...]Moji{last, RawCodePoint(r)})
				mojis[len(mojis)-1] = last
				continue
			}
		}
		last = rune2moji(r, pos)
		pos += int(last.Width())
		mojis = append(mojis, last)
	}
	return mojis
}

func MojiWidthAndCountInString(s string) (width WidthT, count int) {
	var last Moji
	var lastWidth WidthT
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		s = s[size:]

		if last != nil {
			if isZeroWidthJoin(r) {
				next, nextsize := utf8.DecodeRuneInString(s)
				s = s[nextsize:]

				width -= lastWidth
				last = ZeroWidthJoinSequence([...]Moji{last, RawCodePoint(next)})
				lastWidth = last.Width()
				width += lastWidth
				continue
			}
			if isVariationSelectorLike(r) {
				width -= lastWidth
				last = VariationSequence([...]Moji{last, RawCodePoint(r)})
				lastWidth = last.Width()
				width += lastWidth
				continue
			}
			if isEmojiModifier(r) {
				width -= lastWidth
				last = ModifierSequence([...]Moji{last, RawCodePoint(r)})
				lastWidth = last.Width()
				width += lastWidth
				continue
			}
		}
		last = rune2moji(r, int(width))
		lastWidth = last.Width()
		width += lastWidth
		count++
	}
	return
}
