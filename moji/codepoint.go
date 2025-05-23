package moji

import (
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

type Tab struct {
	pos int16
}

func (t *Tab) Len() int {
	return 1
}

func (t *Tab) Width() WidthT {
	return WidthT(4 - t.pos%4)
}

func (t *Tab) PrintTo(w io.Writer) {
	io.WriteString(w, "    "[t.pos%4:])
}

func (t *Tab) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte{'\t'})
	return int64(n), err
}

func (t *Tab) GoString() string {
	return "moji.Tab{}"
}

func (t *Tab) SetPosition(pos int16) {
	t.pos = pos
}

// RawCodePoint is for the character to print as is.
type RawCodePoint rune

func (c RawCodePoint) Len() int {
	return utf8.RuneLen(rune(c))
}

func (c RawCodePoint) Width() WidthT {
	return getWidth(rune(c))
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

func (c RawCodePoint) PrintTo(w io.Writer) {
	writeRune(w, rune(c))
}

func (c RawCodePoint) GoString() string {
	return fmt.Sprintf("moji.RawCodePoint{%X}", rune(c))
}

// EscCodePoint is for the character to print as <XXXXX>
type EscCodePoint rune

func (c EscCodePoint) Len() int {
	return utf8.RuneLen(rune(c))
}

func (c EscCodePoint) Width() WidthT {
	return lenEscaped(rune(c))
}

func (c EscCodePoint) PrintTo(w io.Writer) {
	fmt.Fprintf(w, "<%X>", c)
}

func (c EscCodePoint) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(c))
}

func (c EscCodePoint) GoString() string {
	return fmt.Sprintf("moji.EscCodePoint{%X}", rune(c))
}

type RegionalIndicator rune

func (r RegionalIndicator) Len() int {
	return utf8.RuneLen(rune(r))
}

func (r RegionalIndicator) Width() WidthT {
	return 2
}

func (r RegionalIndicator) PrintTo(w io.Writer) {
	writeRune(w, rune(r))
}

func (r RegionalIndicator) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(r))
}

func (r RegionalIndicator) GoString() string {
	return fmt.Sprintf("moji.RegionalIndicator{%X}", rune(r))
}

// CtrlCodePoint is for the character to print as ^X
type CtrlCodePoint rune

func (c CtrlCodePoint) Len() int {
	return 1
}

func (c CtrlCodePoint) Width() WidthT {
	return 2
}

func (c CtrlCodePoint) PrintTo(w io.Writer) {
	w.Write([]byte{'^', byte('A' + (byte(c) - 1))})
}

func (c CtrlCodePoint) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(c))
}

func (c CtrlCodePoint) GoString() string {
	return fmt.Sprintf("moji.CtrlCodePoint{%X}", rune(c))
}

// WavingWhiteFlagCodePoint is for U+1F3F3 (WAVING WHITE FLAG)
// In WindowsTerminal:
// - "\U0001F3F3"       needs 2cells-width (It should needs 1cell-width)
// - "\U0001F3F3\uFE0F" needs 2cells-width,too.
// (\uFE0F is the variation selector-15)
type WavingWhiteFlagCodePoint rune

func (s WavingWhiteFlagCodePoint) Len() int {
	return utf8.RuneLen(rune(s))
}

func (s WavingWhiteFlagCodePoint) Width() WidthT {
	return 2
}

func (s WavingWhiteFlagCodePoint) PrintTo(w io.Writer) {
	saveCursorAfterN(w, s.Width())
	writeRune(w, rune(s))
	restoreCursor(w)
}

func (s WavingWhiteFlagCodePoint) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(s))
}

func (s WavingWhiteFlagCodePoint) GoString() string {
	return fmt.Sprintf("moji.WavingWhiteFlagCodePoint{%X}", rune(s))
}

const (
	regionalIndicatorBegin       = 0x1F1E6
	regionalIndicatorEnd         = 0x1F1FF
	mathematicalBoldCapitalBegin = 0x1D400
	mathematicalBoldCapitalEnd   = 0x1D7FF
	wavingWhiteFlagCodePoint     = 0x1F3F3
	boxDrawingBegin              = 0x2500
	boxDrawingEnd                = 0x257F
)

func rune2moji(ch rune, pos int) Moji {
	if ch == '\t' {
		return &Tab{pos: int16(pos)}
	} else if ch < ' ' {
		return CtrlCodePoint(ch)
	} else if boxDrawingBegin <= ch && ch <= boxDrawingEnd && AmbiguousIsWide {
		return WavingWhiteFlagCodePoint(ch)
	} else if isToBeEscaped(ch) {
		return EscCodePoint(ch)
	} else if regionalIndicatorBegin <= ch && ch <= regionalIndicatorEnd {
		return RegionalIndicator(ch)
	} else if mathematicalBoldCapitalBegin <= ch && ch <= mathematicalBoldCapitalEnd {
		return WavingWhiteFlagCodePoint(ch)
	} else if ch == wavingWhiteFlagCodePoint {
		return WavingWhiteFlagCodePoint(ch)
	} else {
		return RawCodePoint(ch)
	}
}

func MojiToRune(m Moji) (rune, bool) {
	switch r := m.(type) {
	case RawCodePoint:
		return rune(r), true
	case CtrlCodePoint:
		return rune(r), true
	case EscCodePoint:
		return rune(r), true
	case RegionalIndicator:
		return rune(r), true
	case WavingWhiteFlagCodePoint:
		return rune(r), true
	default:
		return 0, false
	}
}

func IsSpaceMoji(m Moji) bool {
	if r, ok := MojiToRune(m); ok {
		return unicode.IsSpace(r)
	}
	return false
}
