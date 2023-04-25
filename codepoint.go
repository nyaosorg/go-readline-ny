package readline

import (
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

// _RawCodePoint is for the character to print as is.
type _RawCodePoint rune

func (c _RawCodePoint) Width() WidthT {
	return getWidth(rune(c))
}

func writeRune(w io.Writer, r rune) (int64, error) {
	var b [utf8.UTFMax]byte
	n := utf8.EncodeRune(b[:], r)
	n, err := w.Write(b[:n])
	return int64(n), err
}

func (c _RawCodePoint) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(c))
}

func (c _RawCodePoint) PrintTo(w io.Writer) {
	writeRune(w, rune(c))
}

// _EscCodePoint is for the character to print as <XXXXX>
type _EscCodePoint rune

func (c _EscCodePoint) Width() WidthT {
	return lenEscaped(rune(c))
}

func (c _EscCodePoint) PrintTo(w io.Writer) {
	fmt.Fprintf(w, "<%X>", c)
}

func (c _EscCodePoint) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(c))
}

type _RegionalIndicator rune

func (r _RegionalIndicator) Width() WidthT {
	return 2
}

func (r _RegionalIndicator) PrintTo(w io.Writer) {
	writeRune(w, rune(r))
}

func (r _RegionalIndicator) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(r))
}

// _CtrlCodePoint is for the character to print as ^X
type _CtrlCodePoint rune

func (c _CtrlCodePoint) Width() WidthT {
	return 2
}

func (c _CtrlCodePoint) PrintTo(w io.Writer) {
	w.Write([]byte{'^', byte('A' + (byte(c) - 1))})
}

func (c _CtrlCodePoint) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(c))
}

// _WavingWhiteFlagCodePoint is for U+1F3F3 (WAVING WHITE FLAG)
// In WindowsTerminal:
// - "\U0001F3F3"       needs 2cells-width (It should needs 1cell-width)
// - "\U0001F3F3\uFE0F" needs 2cells-width,too.
// (\uFE0F is the variation selector-15)
type _WavingWhiteFlagCodePoint rune

func (s _WavingWhiteFlagCodePoint) Width() WidthT {
	return 2
}

func (s _WavingWhiteFlagCodePoint) PrintTo(w io.Writer) {
	saveCursorAfterN(w, s.Width())
	writeRune(w, rune(s))
	restoreCursor(w)
}

func (s _WavingWhiteFlagCodePoint) WriteTo(w io.Writer) (int64, error) {
	return writeRune(w, rune(s))
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

func rune2moji(ch rune) Moji {
	if ch < ' ' {
		return _CtrlCodePoint(ch)
	} else if boxDrawingBegin <= ch && ch <= boxDrawingEnd && AmbiguousIsWide {
		return _WavingWhiteFlagCodePoint(ch)
	} else if isToBeEscaped(ch) {
		return _EscCodePoint(ch)
	} else if regionalIndicatorBegin <= ch && ch <= regionalIndicatorEnd {
		return _RegionalIndicator(ch)
	} else if mathematicalBoldCapitalBegin <= ch && ch <= mathematicalBoldCapitalEnd {
		return _WavingWhiteFlagCodePoint(ch)
	} else if ch == wavingWhiteFlagCodePoint {
		return _WavingWhiteFlagCodePoint(ch)
	} else {
		return _RawCodePoint(ch)
	}
}

func moji2rune(m Moji) (rune, bool) {
	switch r := m.(type) {
	case _RawCodePoint:
		return rune(r), true
	case _CtrlCodePoint:
		return rune(r), true
	case _EscCodePoint:
		return rune(r), true
	case _RegionalIndicator:
		return rune(r), true
	case _WavingWhiteFlagCodePoint:
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
