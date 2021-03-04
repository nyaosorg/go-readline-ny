package readline

import (
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

var isWindowsTerminal = os.Getenv("WT_SESSION") != "" && os.Getenv("WT_PROFILE_ID") != ""

var (
	SurrogatePairOk         = isWindowsTerminal
	ZeroWidthJoinSequenceOk = isWindowsTerminal
	VariationSequenceOk     = isWindowsTerminal
)

var wtRuneWidth *runewidth.Condition

func init() {
	wtRuneWidth = runewidth.NewCondition()
	if isWindowsTerminal {
		wtRuneWidth.EastAsianWidth = false
	}
}

type Moji interface {
	Width() WidthT
	WriteTo(io.Writer) (int64, error)
	PrintTo(io.Writer)
}

type ZeroWidthJoinSequence [2]Moji

func (s ZeroWidthJoinSequence) Width() WidthT {
	// runewidth.StringWidth should not be used because the width that it gives
	// has no compatible with WindowsTerminal's.
	return s[0].Width() + 1 + s[1].Width()
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

type VariationSequence [2]Moji

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
	// When VariationSequence contains another VariationSequence
	// s[0].PrintTo(w) does not work as we expect.
	s[0].WriteTo(w)
	s[1].WriteTo(w)
	restoreCursor(w)
}

const (
	zeroWidthJoinRune = '\u200D'
	zeroWidthJoinStr  = "\u200D"
)

func isZeroWidthJoinRune(r rune) bool {
	return unicode.Is(unicode.Join_Control, r)
}

func isZeroWidthJoinStr(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return isZeroWidthJoinRune(r)
}

func string2moji(s string) []Moji {
	runes := []rune(s)
	mojis := make([]Moji, 0, len(runes))
	for i := 0; i < len(runes); i++ {
		if ZeroWidthJoinSequenceOk && isZeroWidthJoinRune(runes[i]) && i > 0 && i+1 < len(runes) {
			mojis[len(mojis)-1] =
				ZeroWidthJoinSequence(
					[...]Moji{mojis[len(mojis)-1], RawCodePoint(runes[i+1])})
			i++
		} else if VariationSequenceOk && isVariationSelectorLikeRune(runes[i]) && i > 0 {
			mojis[len(mojis)-1] =
				VariationSequence(
					[...]Moji{mojis[len(mojis)-1], RawCodePoint(runes[i])})

		} else {
			mojis = append(mojis, rune2moji(runes[i]))
		}
	}
	return mojis
}

func moji2string(m []Moji) string {
	var buffer strings.Builder
	for _, m1 := range m {
		m1.WriteTo(&buffer)
	}
	return buffer.String()
}
