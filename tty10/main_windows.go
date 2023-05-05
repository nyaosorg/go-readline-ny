//go:build windows
// +build windows

package tty10

import (
	"golang.org/x/sys/windows"
)

func enable(handle int) (func(), error) {
	h := windows.Handle(handle)
	var orig uint32
	if err := windows.GetConsoleMode(h, &orig); err != nil {
		return func() {}, err
	}
	bits := orig | windows.ENABLE_VIRTUAL_TERMINAL_INPUT
	if err := windows.SetConsoleMode(h, bits); err != nil {
		return func() {}, err
	}
	return func() { windows.SetConsoleMode(h, orig) }, nil
}
