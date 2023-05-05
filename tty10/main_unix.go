//go:build !windows
// +build !windows

package tty10

func enable(handle int) (func(), error) {
	return func() {}, nil
}
