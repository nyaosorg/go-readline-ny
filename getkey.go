package readline

import (
	"strings"
)

type KeyGetter interface {
	Raw() (func() error, error)
	ReadRune() (rune, error)
	Buffered() bool
}

// GetKey reads one-key from tty.
func GetKey(tty1 KeyGetter) (string, error) {
	clean, err := tty1.Raw()
	if err != nil {
		return "", err
	}
	defer clean()

	var buffer strings.Builder
	escape := false
	var surrogated rune = 0
	for {
		r, err := tty1.ReadRune()
		if err != nil {
			return "", err
		}
		if r == 0 {
			continue
		}
		if surrogated > 0 {
			r = 0x10000 + (surrogated-0xD800)*0x400 + (r - 0xDC00)
			surrogated = 0
		} else if 0xD800 <= r && r < 0xE000 { // surrogate pair first word.
			surrogated = r
			continue
		}
		buffer.WriteRune(r)
		if r == '\x1B' {
			escape = true
		}
		if !(escape && tty1.Buffered()) && buffer.Len() > 0 {
			return buffer.String(), nil
		}
	}
}
