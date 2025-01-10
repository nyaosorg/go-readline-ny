package readline

import (
	"io"
	"unicode/utf8"
)

type Highlight struct {
	Pattern  interface{ FindAllStringIndex(string, int) [][]int }
	Sequence string
}

type escapeSequenceType string

func (e escapeSequenceType) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(e))
	return int64(n), err
}

func (e escapeSequenceType) Equals(other colorInterface) bool {
	o, ok := other.(escapeSequenceType)
	return ok && o == e
}

type highlightColorSequence struct {
	colorMap []escapeSequenceType
	index    int
	resetSeq escapeSequenceType
}

func highlightToColoring(input string, resetColor, defaultColor string, H []Highlight) *highlightColorSequence {
	colorMap := make([]escapeSequenceType, len(input))
	for i := 0; i < len(input); i++ {
		colorMap[i] = escapeSequenceType(defaultColor)
	}
	for _, h := range H {
		positions := h.Pattern.FindAllStringIndex(input, -1)
		if positions == nil {
			continue
		}
		for _, p := range positions {
			for i := p[0]; i < p[1]; i++ {
				colorMap[i] = escapeSequenceType(h.Sequence)
			}
		}
	}
	return &highlightColorSequence{
		colorMap: colorMap,
		resetSeq: escapeSequenceType(resetColor),
	}
}

func (H *highlightColorSequence) Init() colorInterface {
	H.index = 0
	return escapeSequenceType(H.resetSeq)
}

func (H *highlightColorSequence) Next(r rune) colorInterface {
	if r == CursorPositionDummyRune {
		return escapeSequenceType("")
	}
	rv := escapeSequenceType(H.colorMap[H.index])
	H.index += utf8.RuneLen(r)
	return rv
}

type colorBridge struct {
	base Coloring
}

func (c *colorBridge) Init() colorInterface {
	return c.base.Init()
}

func (c *colorBridge) Next(r rune) colorInterface {
	return c.base.Next(r)
}
