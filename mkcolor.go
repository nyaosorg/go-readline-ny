package readline

import (
	"regexp"
	"unicode/utf8"
)

type Highlight struct {
	Pattern *regexp.Regexp

	// Sgr is the array of the Select Graphic Redition parameters.
	// []int{0,34,1} represents "ESC[0;34;1m"
	Sgr []int
}

func newColorSequence(sgr []int) ColorSequence {
	var colSeq ColorSequence
	for i := len(sgr) - 1; i >= 0; i-- {
		colSeq |= ColorSequence(sgr[i])
		colSeq <<= 8
	}
	colSeq |= ColorSequence(len(sgr))
	return colSeq
}

type highlightColorSequence struct {
	sgrs  []ColorSequence
	index int
}

func highlightToColoring(input string, H []Highlight) *highlightColorSequence {
	colorMap := make([]ColorSequence, len(input))
	resetSeq := ColorSequence(1)
	for i := range colorMap {
		colorMap[i] = resetSeq
	}
	for _, h := range H {
		positions := h.Pattern.FindAllStringIndex(input, -1)
		if positions == nil {
			continue
		}
		colSeq := newColorSequence(h.Sgr)
		for _, p := range positions {
			for i := p[0]; i < p[1]; i++ {
				colorMap[i] = colSeq
			}
		}
	}
	return &highlightColorSequence{
		sgrs:  colorMap,
		index: 0,
	}
}

func (H *highlightColorSequence) Init() ColorSequence {
	H.index = 0
	return 1 // == "ESC[0m"
}

func (H *highlightColorSequence) Next(r rune) ColorSequence {
	if r == CursorPositionDummyRune {
		return 0
	}
	rv := H.sgrs[H.index]
	H.index += utf8.RuneLen(r)
	return rv
}
