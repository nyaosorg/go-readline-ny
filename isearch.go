package readline

import (
	"context"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

func caseInsensitiveStringContains(s, t string) bool {
	return strings.Contains(strings.ToUpper(s), strings.ToUpper(t))
}

func keyFuncIncSearch(ctx context.Context, this *Buffer) Result {
	var searchBuf strings.Builder
	foundStr := ""
	searchStr := ""
	lastFoundPos := this.History.Len() - 1
	this.backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))

	update := func() {
		for i := this.History.Len() - 1; ; i-- {
			if i < 0 {
				foundStr = ""
				break
			}
			line := this.History.At(i)
			if caseInsensitiveStringContains(line, searchStr) {
				foundStr = line
				lastFoundPos = i
				break
			}
		}
	}
	for {
		drawStr := fmt.Sprintf("(i-search)[%s]:%s", searchStr, foundStr)
		drawWidth := WidthT(0)
		for _, ch := range string2moji(drawStr) {
			w1 := ch.Width()
			if drawWidth+w1 >= this.ViewWidth() {
				break
			}
			ch.PrintTo(this.Out)
			drawWidth += w1
		}
		this.Eraseline()
		io.WriteString(this.Out, ansiCursorOn)
		this.Out.Flush()
		key, err := this.GetKey()
		if err != nil {
			println(err.Error())
			return CONTINUE
		}
		io.WriteString(this.Out, ansiCursorOff)
		this.backspace(drawWidth)
		switch key {
		case "\b":
			searchBuf.Reset()
			// chop last char
			var lastchar rune
			for i, c := range searchStr {
				if i > 0 {
					searchBuf.WriteRune(lastchar)
				}
				lastchar = c
			}
			searchStr = searchBuf.String()
			update()
		case "\r":
			this.ViewStart = 0
			u := &undoT{
				pos:  0,
				text: moji2string(this.Buffer),
			}
			this.undoes = append(this.undoes, u)
			this.Buffer = this.Buffer[:0]
			this.Cursor = 0
			this.ReplaceAndRepaint(0, foundStr)
			return CONTINUE
		case "\x03", "\x07", "\x1B":
			all, _, right := this.view3()
			this.puts(all)
			this.Eraseline()
			this.backspace(right.Width())
			return CONTINUE
		case "\x12":
			for i := lastFoundPos - 1; ; i-- {
				if i < 0 {
					i = this.History.Len() - 1
				}
				if i == lastFoundPos {
					break
				}
				line := this.History.At(i)
				if caseInsensitiveStringContains(line, searchStr) && foundStr != line {
					foundStr = line
					lastFoundPos = i
					break
				}
			}
		case "\x13":
			for i := lastFoundPos + 1; ; i++ {
				if i >= this.History.Len() {
					break
				}
				if i == lastFoundPos {
					break
				}
				line := this.History.At(i)
				if caseInsensitiveStringContains(line, searchStr) && foundStr != line {
					foundStr = line
					lastFoundPos = i
					break
				}
			}
		default:
			charcode, _ := utf8.DecodeRuneInString(key)
			if unicode.IsControl(charcode) {
				break
			}
			searchBuf.WriteRune(charcode)
			searchStr = searchBuf.String()
			update()
		}
	}
}
