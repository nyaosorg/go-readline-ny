//go:build !windows
// +build !windows

package tty10

func nop() {}

func enable(handle int) (func(), error) {
	return nop, nil
}
