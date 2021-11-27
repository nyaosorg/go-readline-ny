package coloring

import (
	"github.com/nyaosorg/go-readline-ny"
)

type VimBatch struct {
	bits int
}

func (s *VimBatch) Init() {
	s.bits = 0
}

const (
	envArea    = 1
	quotedArea = 2
)

func (s *VimBatch) Get(codepoint rune) int {
	newbits := s.bits
	if codepoint == '%' {
		newbits ^= envArea
	} else if codepoint == '"' {
		newbits ^= quotedArea
	}
	color := readline.White
	if ((s.bits | newbits) & envArea) != 0 {
		color = readline.Cyan
	} else if ((s.bits | newbits) & quotedArea) != 0 {
		color = readline.Magenta
	} else {
		color = readline.White
	}
	s.bits = newbits
	return color
}
