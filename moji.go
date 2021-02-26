package readline

import (
	"fmt"
	"io"
	"os"
	"strings"

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
	Put(io.Writer)
	IsSpace() bool
}

type ZeroWidthJoinSequence string

func (s ZeroWidthJoinSequence) Width() WidthT {
	// runewidth.StringWidth should be used because the width that it gives
	// has no compatible with WindowsTerminal's.
	w := WidthT(1)
	for _, c := range string(s) {
		w += WidthT(wtRuneWidth.RuneWidth(c))
	}
	return w
}

func (s ZeroWidthJoinSequence) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(s))
	return int64(n), err
}

func (s ZeroWidthJoinSequence) Put(w io.Writer) {
	io.WriteString(w, string(s))
}

func (s ZeroWidthJoinSequence) IsSpace() bool {
	return false
}

type VariationSequence string

func (s VariationSequence) Width() WidthT {
	return 4
}

func (s VariationSequence) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(s))
	return int64(n), err
}

func (s VariationSequence) Put(w io.Writer) {
	fmt.Fprintf(w, "    \x1B7\b\b\b\b%s\x1B8", string(s))
}

func (s VariationSequence) IsSpace() bool {
	return false
}

const (
	zeroWidthJoinRune = '\u200D'
	zeroWidthJoinStr  = "\u200D"
)

func string2moji(s string) []Moji {
	runes := []rune(s)
	mojis := make([]Moji, 0, len(runes))
	for i := 0; i < len(runes); i++ {
		if ZeroWidthJoinSequenceOk && runes[i] == zeroWidthJoinRune && i > 0 && i+1 < len(runes) {
			mojis[len(mojis)-1] = ZeroWidthJoinSequence(string(runes[i-1 : i+2]))
			i++
		} else if VariationSequenceOk && isVariationSelector(runes[i]) && i > 0 {
			mojis[len(mojis)-1] = VariationSequence(string(runes[i-1 : i+1]))

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
