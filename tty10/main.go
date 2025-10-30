package tty10

import (
	"os"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
)

type Tty struct {
	done    chan struct{}
	disable func()
	buf     []byte
}

func (M *Tty) Open(onSize func(int)) error {
	var err error

	stdin := int(os.Stdin.Fd())
	M.disable, err = enable(stdin)
	if err != nil {
		return err
	}

	if onSize != nil {
		w, _, err := M.Size()
		if err != nil {
			return err
		}
		M.done = make(chan struct{})
		go func(lastw int) {
			ticker := time.NewTicker(time.Second)
			for {
				select {
				case <-M.done:
					ticker.Stop()
					return
				case <-ticker.C:
					w, _, err := M.Size()
					if err == nil && w != lastw {
						onSize(w)
						lastw = w
					}
				}
			}
		}(w)
	}
	return nil
}

func (M *Tty) getKey() ([]byte, error) {
	stdin := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(stdin)
	if err != nil {
		return nil, err
	}
	defer term.Restore(stdin, oldState)

	var buffer [1024]byte
	n, err := os.Stdin.Read(buffer[:])
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}

func (M *Tty) GetKey() (string, error) {
	if len(M.buf) <= 0 {
		var err error
		M.buf, err = M.getKey()
		if err != nil || len(M.buf) <= 0 {
			return "", err
		}
	}
	var result string
	r, size := utf8.DecodeRune(M.buf)
	if r == '\x1B' {
		result = string(M.buf)
		M.buf = nil
	} else {
		result = string(M.buf[:size])
		M.buf = M.buf[size:]
	}
	return result, nil
}

func (M *Tty) Close() error {
	if M.done != nil {
		M.done <- struct{}{}
		close(M.done)
		M.done = nil
	}
	if M.disable != nil {
		M.disable()
		M.disable = nil
	}
	return nil
}

func (M *Tty) Size() (int, int, error) {
	return term.GetSize(int(os.Stderr.Fd()))
}
