package readline

type VimBatchColor struct {
	bits int
}

func (s *VimBatchColor) Init() {
	s.bits = 0
}

func (s *VimBatchColor) Get(codepoint rune) ColorCode {
	newbits := s.bits
	if codepoint == '%' {
		newbits ^= 1
	} else if codepoint == '"' {
		newbits ^= 2
	}
	var color ColorCode = 37
	if ((s.bits | newbits) & 1) != 0 {
		color = 36
	} else if ((s.bits | newbits) & 2) != 0 {
		color = 35
	} else {
		color = 37
	}
	s.bits = newbits
	return color
}
