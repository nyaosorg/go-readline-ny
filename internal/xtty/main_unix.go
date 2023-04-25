//go:build !windows
// +build !windows

package xtty

func nop() {}

func enable(handle int) (func(), error) {
	return nop, nil
}
