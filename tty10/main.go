package tty10

import (
	"os"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
)

type Tty struct {
	buffer [128]byte
	key   []byte
	done   chan struct{}
	ticker *time.Ticker
}

func (M *Tty) Open() error {
	M.done = make(chan struct{})
	M.ticker = time.NewTicker(time.Second)
	return nil
}

func (*Tty) Raw() (func() error, error) {
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

func (M *Tty) ReadRune() (rune, error) {
	if M.key == nil || len(M.key) <= 0 {
		n, err := os.Stdin.Read(M.buffer[:])
		if err != nil {
			return 0, err
		}
		M.key = M.buffer[:n]
	}
	r, size := utf8.DecodeRune(M.key)
	M.key = M.key[size:]

	return r, nil
}

func (M *Tty) Buffered() bool {
	return len(M.key) > 0
}

func (M *Tty) Close() error {
	M.done <- struct{}{}
	return nil
}

func (M *Tty) Size() (int, int, error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

func (M *Tty) GetResizeNotifier() func() (int, int, bool) {
	return func() (int, int, bool) {
		select {
		case <-M.done:
			M.ticker.Stop()
			return 0,0, false
		case <-M.ticker.C:
			w, h, err := M.Size()
			return w, h, err == nil
		}
	}
}
