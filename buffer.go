package readline

import (
	"fmt"
	"strings"
	"unicode"
)

const forbiddenWidth WidthT = 3 // = lastcolumn(1) and FULLWIDTHCHAR-SIZE(2)

type undoT struct {
	pos  int
	del  int
	text string
}

type Buffer struct {
	*Editor
	Buffer         []Moji
	TTY            KeyGetter
	ViewStart      int
	termWidth      int // == topColumn + termWidth + forbiddenWidth
	topColumn      int // == width of Prompt
	HistoryPointer int
	undoes         []*undoT
	pending        string
}

func (this *Buffer) ViewWidth() WidthT {
	return WidthT(this.termWidth) - WidthT(this.topColumn) - forbiddenWidth
}

func (this *Buffer) view() _Range {
	view := this.Buffer[this.ViewStart:]
	width := this.ViewWidth()
	w := WidthT(0)
	for i, c := range view {
		w += c.Width()
		if w >= width {
			return view[:i]
		}
	}
	return _Range(view)
}

func (this *Buffer) view3() (_Range, _Range, _Range) {
	v := this.view()
	x := this.Cursor - this.ViewStart
	return v, v[:x], v[x:]
}

func (this *Buffer) insert(csrPos int, insStr []Moji) {
	// expand buffer
	this.Buffer = append(this.Buffer, insStr...)

	// shift original string to make area
	copy(this.Buffer[csrPos+len(insStr):], this.Buffer[csrPos:])

	// insert insStr
	copy(this.Buffer[csrPos:csrPos+len(insStr)], insStr)

	u := &undoT{
		pos: csrPos,
		del: len(insStr),
	}
	this.undoes = append(this.undoes, u)
}

func (this *Buffer) insertString(pos int, s string) _Range {
	list := string2moji(s)
	this.insert(pos, list)
	return _Range(list)
}

// Insert String :s at :pos (Do not update screen)
// returns
//    count of rune
func (this *Buffer) InsertString(pos int, s string) int {
	return len(this.insertString(pos, s))
}

// Delete remove Buffer[pos:pos+n].
// It returns the width to clear the end of line.
// It does not update screen.
func (this *Buffer) Delete(pos int, n int) WidthT {
	if n <= 0 || len(this.Buffer) < pos+n {
		return 0
	}
	u := &undoT{
		pos:  pos,
		text: moji2string(this.Buffer[pos : pos+n]),
	}
	this.undoes = append(this.undoes, u)

	delw := this.GetWidthBetween(pos, pos+n)
	copy(this.Buffer[pos:], this.Buffer[pos+n:])
	this.Buffer = this.Buffer[:len(this.Buffer)-n]
	return delw
}

// ResetViewStart set ViewStart the new value which should be.
// It does not update screen.
func (this *Buffer) ResetViewStart() {
	this.ViewStart = 0
	w := WidthT(0)
	for i := 0; i <= this.Cursor && i < len(this.Buffer); i++ {
		w += this.Buffer[i].Width()
		for w >= this.ViewWidth() {
			if this.ViewStart >= len(this.Buffer) {
				// When standard output is redirected.
				return
			}
			w -= this.Buffer[this.ViewStart].Width()
			this.ViewStart++
		}
	}
}

func (this *Buffer) GetWidthBetween(from int, to int) WidthT {
	return _Range(this.Buffer[from:to]).Width()
}

func (this *Buffer) SubString(start, end int) string {
	return moji2string(this.Buffer[start:end])
}

func (this Buffer) String() string {
	return moji2string(this.Buffer)
}

var Delimiters = "\"'"

func (this *Buffer) CurrentWordTop() (wordTop int) {
	wordTop = -1
	quotedchar := '\000'
	for i, moji := range this.Buffer[:this.Cursor] {
		if ch, ok := moji2rune(moji); ok {
			if quotedchar == '\000' {
				if strings.ContainsRune(Delimiters, ch) {
					quotedchar = ch
				}
			} else if ch == quotedchar {
				quotedchar = '\000'
			}
			if unicode.IsSpace(ch) && quotedchar == '\000' {
				wordTop = -1
			} else if wordTop < 0 {
				wordTop = i
			}
		}
	}
	if wordTop < 0 {
		return this.Cursor
	} else {
		return wordTop
	}
}

func (this *Buffer) CurrentWord() (string, int) {
	start := this.CurrentWordTop()
	return this.SubString(start, this.Cursor), start
}

func (this *Buffer) joinUndo() {
	if len(this.undoes) < 2 {
		return
	}
	u1 := this.undoes[len(this.undoes)-2]
	u2 := this.undoes[len(this.undoes)-1]
	if u1.pos != u2.pos {
		return
	}
	if u1.text != "" && u2.text != "" {
		return
	}
	if u1.del != 0 && u2.del != 0 {
		return
	}
	if u1.text == "" {
		u1.text = u2.text
	}
	if u1.del == 0 {
		u1.del = u2.del
	}
	this.undoes = this.undoes[:len(this.undoes)-1]
}

func (b *Buffer) startChangeWidthEventLoop(lastw_ int, getResizeEvent func() (int, int, bool)) {
	go func(lastw int) {
		for {
			w, _, ok := getResizeEvent()
			if !ok {
				break
			}
			if lastw != w {
				mu.Lock()
				b.termWidth = w
				fmt.Fprintf(b.Out, "\x1B[%dG", b.topColumn+1)
				b.RepaintAfterPrompt()
				mu.Unlock()
				lastw = w
			}
		}
	}(lastw_)
}

// GetKey reads one-key from tty.
func (this *Buffer) GetKey() (string, error) {
	return GetKey(this.TTY)
}
