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

var FunISearchBackward = &KeyGoFuncT{
	Name: F_ISEARCH_BACKWARD,
	Func: func(ctx context.Context, this *Buffer) Result {
		var searchBuf strings.Builder
		foundStr := ""
		searchStr := ""
		lastFoundPos := this.History.Len() - 1
		this.GotoHead()

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
			for _, ch := range StringToMoji(drawStr) {
				w1 := ch.Width()
				if drawWidth+w1 >= this.ViewWidth() {
					break
				}
				ch.PrintTo(this.Out)
				drawWidth += w1
			}
			this.eraseline()
			io.WriteString(this.Out, ansiCursorOn)
			this.Out.Flush()
			key, err := this.GetKey()
			if err != nil {
				println(err.Error())
				return CONTINUE
			}
			io.WriteString(this.Out, ansiCursorOff)
			this.GotoHead()
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
				u := &_Undo{
					pos:  0,
					text: cell2string(this.Buffer),
				}
				this.undoes = append(this.undoes, u)
				this.Buffer = this.Buffer[:0]
				this.Cursor = 0
				this.ReplaceAndRepaint(0, foundStr)
				return CONTINUE
			case "\x03", "\x07", "\x1B":
				all, left := this.view2()
				this.puts(all)
				this.eraseline()
				this.GotoHead()
				this.puts(left)
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
	},
}
