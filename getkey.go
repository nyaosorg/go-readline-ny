package readline

import (
	"strings"

	"github.com/mattn/go-tty"
)

type KeyGetter interface {
	Raw() (func() error, error)
	ReadRune() (rune, error)
	Buffered() bool
	Close() error
	Size() (int, int, error)

	GetChangeWidthEvent() func() int
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

type DefaultTty struct {
	*tty.TTY
}

func NewDefaultTty() (KeyGetter, error) {
	tty1, err := tty.Open()
	if err != nil {
		return nil, err
	}
	return &DefaultTty{TTY: tty1}, nil
}

// GetChangeWidthEvent is the wrapper for the channel for resize-event.
// It returns the function to get the new width on resized.
// When the channel is closed, it returns -1.
// The reason to need the wrapper is to remove the dependency
// on "mattn/go-tty".WINSIZE.
func (t *DefaultTty) GetChangeWidthEvent() func() int {
	ws := t.TTY.SIGWINCH()
	return func() int {
		ws1, ok := <-ws
		if !ok {
			return -1
		}
		return ws1.W
	}
}
