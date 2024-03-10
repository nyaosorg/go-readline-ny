package tty10

import (
	"os"
	"time"

	"golang.org/x/term"
)

type Tty struct {
	done    chan struct{}
	disable func()
}

func (M *Tty) Open(onSize func(int)) error {
	var err error

	stdin := int(os.Stdin.Fd())
	M.disable, err = enable(stdin)
	if err != nil {
		return err
	}

	M.done = make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-M.done:
				ticker.Stop()
				return
			case <-ticker.C:
				w, _, err := M.Size()
				if err == nil {
					onSize(w)
				}
			}
		}
	}()

	return nil
}

func (M *Tty) GetKey() (string, error) {
	stdin := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(stdin)
	if err != nil {
		return "", err
	}
	defer term.Restore(stdin, oldState)

	var buffer [100]byte
	n, err := os.Stdin.Read(buffer[:])
	if err != nil {
		return "", err
	}
	return string(buffer[:n]), nil
}

func (M *Tty) Close() error {
	if M.done != nil {
		M.done <- struct{}{}
		close(M.done)
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
