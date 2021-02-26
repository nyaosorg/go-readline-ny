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

func (c RawCodePoint) IsSpace() bool {
	return unicode.IsSpace(rune(c))
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

func (c EscCodePoint) IsSpace() bool {
	return false
}

// EscCodePoint is for the charactor to print as ^X
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

func (c CtrlCodePoint) IsSpace() bool {
	return unicode.IsSpace(rune(c))
}

func rune2moji(ch rune) Moji {
	if ch < ' ' {
		return CtrlCodePoint(ch)
	} else if isToBeEscaped(ch) {
		return EscCodePoint(ch)
	} else {
		return RawCodePoint(ch)
	}
}
