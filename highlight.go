package readline

import (
	"io"
	"unicode/utf8"
)

type Highlight struct {
	Pattern  interface{ FindAllStringIndex(string, int) [][]int }
	Sequence string
}

type EscapeSequenceId uint

var (
	escapeSequences          = []string{}
	escapeSequenceStringToId = map[string]EscapeSequenceId{}
)

func NewEscapeSequenceId(s string) EscapeSequenceId {
	if code, ok := escapeSequenceStringToId[s]; ok {
		return code
	}
	code := EscapeSequenceId(len(escapeSequences))
	escapeSequences = append(escapeSequences, s)
	escapeSequenceStringToId[s] = code
	return code
}

func (e EscapeSequenceId) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, escapeSequences[e])
	return int64(n), err
}

func (e EscapeSequenceId) Equals(other ColorInterface) bool {
	o, ok := other.(EscapeSequenceId)
	return ok && o == e
}

type HighlightColorSequence struct {
	colorMap []EscapeSequenceId
	index    int
	resetSeq EscapeSequenceId
}

func HighlightToColoring(input string, resetColor, defaultColor string, H []Highlight) *HighlightColorSequence {
	colorMap := make([]EscapeSequenceId, len(input))
	defaultSeq := NewEscapeSequenceId(defaultColor)
	for i := 0; i < len(input); i++ {
		colorMap[i] = defaultSeq
	}
	for _, h := range H {
		positions := h.Pattern.FindAllStringIndex(input, -1)
		if positions == nil {
			continue
		}
		seq := NewEscapeSequenceId(h.Sequence)
		for _, p := range positions {
			for i := p[0]; i < p[1]; i++ {
				colorMap[i] = seq
			}
		}
	}
	return &HighlightColorSequence{
		colorMap: colorMap,
		resetSeq: NewEscapeSequenceId(resetColor),
	}
}

func (H *HighlightColorSequence) Init() ColorInterface {
	H.index = 0
	return H.resetSeq
}

func (H *HighlightColorSequence) Next(r rune) ColorInterface {
	if r == CursorPositionDummyRune {
		return NewEscapeSequenceId("")
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

func (c *colorBridge) Init() ColorInterface {
	return c.base.Init()
}

func (c *colorBridge) Next(r rune) ColorInterface {
	return c.base.Next(r)
}
