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

// Deprecated: it has a problem where only the first line is taken when multiple lines are pasted from the clipboard, etc
func GetKey(tty XTty) (string, error) {
	clean, err := tty.Raw()
	if err != nil {
		return "", err
	}
	defer clean()

	return getOneKey(tty)
}

func getOneKey(tty XTty) (string, error) {
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

func getKeys(tty XTty) ([]string, error) {
	clean, err := tty.Raw()
	if err != nil {
		return nil, err
	}
	defer clean()

	keys := []string{}

	for {
		key1, err := getOneKey(tty)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key1)
		if !tty.Buffered() {
			return keys, nil
		}
	}
}
