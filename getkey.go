package readline

import (
	"strings"
	"unicode/utf16"

	"github.com/nyaosorg/go-readline-ny/internal/github.com/mattn/go-tty"
)

type MinimumTty interface {
	Raw() (func() error, error)
	ReadRune() (rune, error)
	Buffered() bool
}

// KeyGetter is the interface from which the ReadLine can read console input
type KeyGetter interface {
	MinimumTty
	Close() error
	Size() (int, int, error)

	GetResizeNotifier() func() (int, int, bool)
}

// GetKey reads one-key from tty.
func GetKey(tty1 MinimumTty) (string, error) {
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
		if !(escape && tty1.Buffered()) && buffer.Len() > 0 {
			return buffer.String(), nil
		}
	}
}

type _DefaultTty struct {
	*tty.TTY
}

// NewDefaultTty returns the instance for KeyGetter, which is the customized
// version of go-tty.TTY
func NewDefaultTty() (KeyGetter, error) {
	tty1, err := tty.Open()
	if err != nil {
		return nil, err
	}
	return &_DefaultTty{TTY: tty1}, nil
}

// GetResizeNotifier is the wrapper for the channel for resize-event.
// It returns the function to get the new screen-size on resized.
// When the channel is closed, it returns false as the third value.
// The reason to need the wrapper is to remove the dependency
// on "mattn/go-tty".WINSIZE .
func (t *_DefaultTty) GetResizeNotifier() func() (int, int, bool) {
	ws := t.TTY.SIGWINCH()
	return func() (int, int, bool) {
		ws1, ok := <-ws
		return ws1.W, ws1.H, ok
	}
}
