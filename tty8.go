package readline

import (
	"fmt"

	"github.com/mattn/go-tty"
)

type _Tty struct {
	*tty.TTY
	buf []string
}

func (m *_Tty) Open(onResize func(int)) error {
	var err error
	m.TTY, err = tty.Open()
	if err != nil {
		return fmt.Errorf("go-tty.Open: %w", err)
	}
	_lastw, _, err := m.TTY.Size()
	if err != nil {
		return fmt.Errorf("go-tty.Size: %w", err)
	}
	ws := m.TTY.SIGWINCH()
	go func(lastw int) {
		for {
			wh, ok := <-ws
			if !ok {
				break
			}
			if lastw != wh.W {
				onResize(wh.W)
				lastw = wh.W
			}
		}
	}(_lastw)
	return nil
}

func (m *_Tty) GetKey() (string, error) {
	if len(m.buf) <= 0 {
		var err error
		m.buf, err = getKeys(m.TTY)
		if err != nil || len(m.buf) <= 0 {
			return "", err
		}
	}
	var top string
	top, m.buf = m.buf[0], m.buf[1:]
	return top, nil
}

func (m *_Tty) Close() error {
	if m.TTY != nil {
		m.TTY.Close()
		m.TTY = nil
	}
	return nil
}
