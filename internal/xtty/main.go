package xtty

import (
	"os"
	"unicode/utf8"

	"golang.org/x/term"
)

type TTY struct {
	buffer [128]byte
	text   []byte
}

func (*TTY) Raw() (func() error, error) {
	stdin := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(stdin)
	if err != nil {
		return func() error { return nil }, err
	}
	disable, err := enable(stdin)
	if err != nil {
		return func() error { return nil }, err
	}
	return func() error {
		disable()
		term.Restore(stdin, oldState)
		return nil
	}, nil
}

func (M *TTY) ReadRune() (rune, error) {
	if M.text == nil || len(M.text) <= 0 {
		n, err := os.Stdin.Read(M.buffer[:])
		if err != nil {
			return 0, err
		}
		M.text = M.buffer[:n]
	}
	r, size := utf8.DecodeRune(M.text)
	M.text = M.text[size:]

	return r, nil
}

func (M *TTY) Buffered() bool {
	return len(M.text) > 0
}

func (M *TTY) Close() error {
	return nil
}

func (M *TTY) Size() (int, int, error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

func (M *TTY) GetResizeNotifier() func() (int, int, bool) {
	return func() (int, int, bool) {
		w, h, err := M.Size()
		return w, h, err == nil
	}
}
