package coloring

import (
	"github.com/nyaosorg/go-readline-ny"
)

type VimBatch struct {
	bits int
}

func (s *VimBatch) Init() readline.ColorSequence {
	s.bits = 0
	return readline.ColorReset
}

const (
	envArea    = 1
	quotedArea = 2
)

func (s *VimBatch) Next(codepoint rune) readline.ColorSequence {
	newbits := s.bits
	if codepoint == '%' {
		newbits ^= envArea
	} else if codepoint == '"' {
		newbits ^= quotedArea
	}
	color := readline.DefaultForeGroundColor
	if codepoint == '\u3000' {
		color = readline.SGR3(37, 22, 41)
	} else if ((s.bits | newbits) & envArea) != 0 {
		color = readline.Cyan
	} else if ((s.bits | newbits) & quotedArea) != 0 {
		color = readline.Magenta
	} else if codepoint == '&' {
		color = readline.DarkYellow
	} else {
		color = readline.DefaultForeGroundColor
	}
	s.bits = newbits
	return color
}
