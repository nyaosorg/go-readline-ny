package runewidth

import (
	"github.com/mattn/go-runewidth"
)

type Condition struct {
	*runewidth.Condition
}

func New(ambiguousIsWide bool) Condition {
	var c Condition
	c.Condition = runewidth.NewCondition()
	c.Condition.EastAsianWidth = ambiguousIsWide
	return c
}
