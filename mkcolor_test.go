package readline

import (
	"regexp"
	"testing"
)

func TestColorSequence(t *testing.T) {
	testcase := []struct {
		expect ColorSequence
		source []int
	}{
		{expect: 0, source: []int{}},
	}

	for _, case1 := range testcase {
		result := newColorSequence(case1.source)
		if result != case1.expect {
			t.Fatalf("expect %#v, but %#v from %#v",
				case1.expect, result, case1.source)
		}
	}
}

func TestNewHighlight(t *testing.T) {
	hl := []Highlight{
		{regexp.MustCompile(`"[^"]*"`), []int{1, 31, 40}},
	}
	r1 := ColorSequence(3 | (1 << 8) | (31 << 16) | (40 << 24))

	result := highlightToColoring(`123"456"789"ABC"D`, hl)
	expect := []ColorSequence{
		0, 0, 0, r1, r1, r1, r1, r1, 0, 0, 0, r1, r1, r1, r1, r1, 0,
	}
	if len(result.sgrs) != len(expect) {
		t.Fatalf("expect len=%d, but %d", len(expect), len(result.sgrs))
	}
	for i, r := range result.sgrs {
		if e := expect[i]; e != r {
			t.Fatalf("[%d] expect %v, but %v", i, e, r)
		}
	}
}
