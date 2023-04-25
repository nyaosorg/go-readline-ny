package tty8

import (
	"github.com/mattn/go-tty"
)

type Tty struct {
	*tty.TTY
}

// New returns the instance for KeyGetter, which is the customized
// version of go-tty.TTY
func (t *Tty) Open() (err error) {
	t.TTY,err = tty.Open()
	return err
}

// GetResizeNotifier is the wrapper for the channel for resize-event.
// It returns the function to get the new screen-size on resized.
// When the channel is closed, it returns false as the third value.
// The reason to need the wrapper is to remove the dependency
// on "mattn/go-tty".WINSIZE .
func (t *Tty) GetResizeNotifier() func() (int, int, bool) {
	ws := t.TTY.SIGWINCH()
	return func() (int, int, bool) {
		ws1, ok := <-ws
		return ws1.W, ws1.H, ok
	}
}
