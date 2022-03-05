package dummyin

import (
	"unicode/utf8"

	"github.com/nyaosorg/go-readline-ny"
)

type _Tty struct {
	Text []string
}

func (*_Tty) Raw() (func() error, error) {
	return func() error { return nil }, nil
}

func (M *_Tty) ReadRune() (rune, error) {
	if M.Text == nil || len(M.Text) <= 0 {
		return '\r', nil
	}
	if len(M.Text[0]) <= 0 {
		M.Text = M.Text[1:]
	}
	if M.Text == nil || len(M.Text) <= 0 {
		return '\r', nil
	}
	r, size := utf8.DecodeRuneInString(M.Text[0])
	M.Text[0] = M.Text[0][size:]

	return r, nil
}

func (M *_Tty) Buffered() bool {
	return len(M.Text[0]) > 0
}

func (M *_Tty) Close() error {
	return nil
}

func (M *_Tty) Size() (int, int, error) {
	return 80, 25, nil
}

func (M *_Tty) GetResizeNotifier() func() (int, int, bool) {
	return func() (int, int, bool) {
		return 80, 25, true
	}
}

func New(texts ...string) func() (readline.KeyGetter, error) {
	return func() (readline.KeyGetter, error) {
		return &_Tty{Text: texts}, nil
	}
}
