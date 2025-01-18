package readline

import (
	"io"
	"unicode/utf8"
)

type Highlight struct {
	Pattern  interface{ FindAllStringIndex(string, int) [][]int }
	Sequence string
}

type escapeSequenceId uint

var (
	escapeSequences          = []string{}
	escapeSequenceStringToId = map[string]escapeSequenceId{}
)

func newEscapeSequenceId(s string) escapeSequenceId {
	if code, ok := escapeSequenceStringToId[s]; ok {
		return code
	}
	code := escapeSequenceId(len(escapeSequences))
	escapeSequences = append(escapeSequences, s)
	escapeSequenceStringToId[s] = code
	return code
}

func (e escapeSequenceId) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, escapeSequences[e])
	return int64(n), err
}

func (e escapeSequenceId) Equals(other colorInterface) bool {
	o, ok := other.(escapeSequenceId)
	return ok && o == e
}

type highlightColorSequence struct {
	colorMap []escapeSequenceId
	index    int
	resetSeq escapeSequenceId
}

func highlightToColoring(input string, resetColor, defaultColor string, H []Highlight) *highlightColorSequence {
	colorMap := make([]escapeSequenceId, len(input))
	defaultSeq := newEscapeSequenceId(defaultColor)
	for i := 0; i < len(input); i++ {
		colorMap[i] = defaultSeq
	}
	for _, h := range H {
		positions := h.Pattern.FindAllStringIndex(input, -1)
		if positions == nil {
			continue
		}
		seq := newEscapeSequenceId(h.Sequence)
		for _, p := range positions {
			for i := p[0]; i < p[1]; i++ {
				colorMap[i] = seq
			}
		}
	}
	return &highlightColorSequence{
		colorMap: colorMap,
		resetSeq: newEscapeSequenceId(resetColor),
	}
}

func (H *highlightColorSequence) Init() colorInterface {
	H.index = 0
	return H.resetSeq
}

func (H *highlightColorSequence) Next(r rune) colorInterface {
	if r == CursorPositionDummyRune {
		return newEscapeSequenceId("")
	}
	if H.index >= len(H.colorMap) {
		return H.resetSeq
	}
	rv := H.colorMap[H.index]
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
