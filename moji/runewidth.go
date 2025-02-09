//go:build !orgxwidth

package moji

import (
	"github.com/mattn/go-runewidth"
)

func newRuneWidth(ambiguousIsWide bool) func(rune) int {
	c := runewidth.NewCondition()
	c.EastAsianWidth = ambiguousIsWide
	return c.RuneWidth
}

var runeWidth = newRuneWidth(AmbiguousIsWide)
