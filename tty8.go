//go:build !tty10

package readline

import "github.com/nyaosorg/go-readline-ny/tty8"

type defaultTty = tty8.Tty
