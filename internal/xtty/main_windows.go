//go:build windows
// +build windows

package xtty

import (
	"golang.org/x/sys/windows"
)

const _ENABLE_VIRTUAL_TERMINAL_INPUT = 0x0200

func nop() {}

func enable(handle int) (func(), error) {
	h := windows.Handle(handle)
	var orig uint32
	if err := windows.GetConsoleMode(h, &orig); err != nil {
		return nop, err
	}
	bits := orig | _ENABLE_VIRTUAL_TERMINAL_INPUT
	if err := windows.SetConsoleMode(h, bits); err != nil {
		return nop, err
	}
	return func() { windows.SetConsoleMode(h, orig) }, nil
}
