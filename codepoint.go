package readline

import (
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

// RawCodePoint is for the charactor to print as is.
type RawCodePoint rune

func (c RawCodePoint) Width() WidthT {
	r := rune(c)
	width, ok := widthCache[r]
	if !ok {
		width = WidthT(wtRuneWidth.RuneWidth(r))
	}
	return width
}

func writeRune(w io.Writer, r rune) (int64, error) {
	var b [utf8.UTFMax]byte
	n := utf8.EncodeRune(b[:], r)
	n, err := w.Write(b[:n])
	return int64(n), err
}

func (c RawCodePoint) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(c))
}

func (c RawCodePoint) Put(w io.Writer) {
	writeRune(w, rune(c))
}

// EscCodePoint is for the charactor to print as <XXXXX>
type EscCodePoint rune

func (c EscCodePoint) Width() WidthT {
	return lenEscaped(rune(c))
}

func (c EscCodePoint) Put(w io.Writer) {
	fmt.Fprintf(w, "<%X>", c)
}

func (c EscCodePoint) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(c))
}

type RegionalIndicator rune

func (r RegionalIndicator) Width() WidthT {
	return 2
}

func (r RegionalIndicator) Put(w io.Writer) {
	writeRune(w, rune(r))
}

func (r RegionalIndicator) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(r))
}

// CtrlCodePoint is for the charactor to print as ^X
type CtrlCodePoint rune

func (c CtrlCodePoint) Width() WidthT {
	return 2
}

func (c CtrlCodePoint) Put(w io.Writer) {
	w.Write([]byte{'^', byte('A' + (byte(c) - 1))})
}

func (c CtrlCodePoint) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(c))
}

// WavingWhiteFlagCodePoint is for U+1F3F3 (WAVING WHITE FLAG)
// In WindowsTerminal:
// - "\U0001F3F3"       needs 2cells-width (It should needs 1cell-width)
// - "\U0001F3F3\uFE0F" needs 2cells-width,too.
// (\uFE0F is the variation selector-15)
type WavingWhiteFlagCodePoint rune

func (s WavingWhiteFlagCodePoint) Width() WidthT {
	return 2
}

func (s WavingWhiteFlagCodePoint) Put(w io.Writer) {
	saveCursorAfterN(w, s.Width())
	writeRune(w, rune(s))
	restoreCursor(w)
}

func (s WavingWhiteFlagCodePoint) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(s))
}

func rune2moji(ch rune) Moji {
	if ch < ' ' {
		return CtrlCodePoint(ch)
	} else if isToBeEscaped(ch) {
		return EscCodePoint(ch)
	} else if 0x1F1E6 <= ch && ch <= 0x1F1FF {
		return RegionalIndicator(ch)
	} else if ch == 0x1F3F3 {
		return WavingWhiteFlagCodePoint(ch)
	} else {
		return RawCodePoint(ch)
	}
}

func moji2rune(m Moji) (rune, bool) {
	switch r := m.(type) {
	case RawCodePoint:
		return rune(r), true
	case CtrlCodePoint:
		return rune(r), true
	case EscCodePoint:
		return rune(r), true
	case RegionalIndicator:
		return rune(r), true
	default:
		return 0, false
	}
}

func isSpaceMoji(m Moji) bool {
	if r, ok := moji2rune(m); ok {
		return unicode.IsSpace(r)
	}
	return false
}
