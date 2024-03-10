package readline

import (
	"strings"
	"unicode/utf16"
)

// XTty is the interface of tty to use GetKey function.
type XTty interface {
	Raw() (func() error, error)
	ReadRune() (rune, error)
	Buffered() bool
}

// GetKey reads one-key from *tty*.
// The *tty* object must have Raw(),ReadRune(), and Buffered() method.
func GetKey(tty XTty) (string, error) {
	clean, err := tty.Raw()
	if err != nil {
		return "", err
	}
	defer clean()

	var buffer strings.Builder
	escape := false
	var surrogated rune = 0
	for {
		r, err := tty.ReadRune()
		if err != nil {
			return "", err
		}
		if r == 0 {
			continue
		}
		if surrogated > 0 {
			r = utf16.DecodeRune(surrogated, r)
			surrogated = 0
		} else if utf16.IsSurrogate(r) { // surrogate pair first word.
			surrogated = r
			continue
		}
		buffer.WriteRune(r)
		if r == '\x1B' {
			escape = true
		}
		if !(escape && tty.Buffered()) && buffer.Len() > 0 {
			return buffer.String(), nil
		}
	}
}
