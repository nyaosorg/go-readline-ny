package dummyin

import (
	"unicode/utf8"
)

type Tty struct {
	Text []string
}

func (*Tty) Open() error {
	return nil
}

func (*Tty) Raw() (func() error, error) {
	return func() error { return nil }, nil
}

func (M *Tty) ReadRune() (rune, error) {
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

func (M *Tty) Buffered() bool {
	return len(M.Text[0]) > 0
}

func (M *Tty) Close() error {
	return nil
}

func (M *Tty) Size() (int, int, error) {
	return 80, 25, nil
}

func (M *Tty) GetResizeNotifier() func() (int, int, bool) {
	return func() (int, int, bool) {
		return 80, 25, true
	}
}
