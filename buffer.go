package readline

import (
	"io"
	"strings"
	"unicode"

	"github.com/nyaosorg/go-readline-ny/internal/moji"
)

const forbiddenWidth WidthT = 3 // = lastcolumn(1) and FULLWIDTHCHAR-SIZE(2)

type _Undo struct {
	pos  int
	del  int
	text string
}

type Cell struct {
	Moji     Moji
	color    ColorSequence
	position int16
}

func (C Cell) String() string {
	var buffer strings.Builder
	C.Moji.WriteTo(&buffer)
	return buffer.String()
}

// Buffer is ReadLine's internal data structure
type Buffer struct {
	*Editor
	Buffer         []Cell
	ViewStart      int
	termWidth      int // == topColumn + termWidth + forbiddenWidth
	topColumn      int // == width of Prompt
	historyPointer int
	undoes         []*_Undo
	pending        string
}

// ViewWidth returns the cell-width screen can show in the one-line.
func (B *Buffer) ViewWidth() WidthT {
	return WidthT(B.termWidth) - WidthT(B.topColumn) - forbiddenWidth
}

func (B *Buffer) view() _Range {
	view := B.Buffer[B.ViewStart:]
	width := B.ViewWidth()
	w := WidthT(0)
	for i, c := range view {
		w += c.Moji.Width()
		if w >= width {
			return view[:i]
		}
	}
	return _Range(view)
}

func (B *Buffer) view2() (all _Range, before _Range) {
	v := B.view()
	x := B.Cursor - B.ViewStart
	return v, v[:x]
}

func (B *Buffer) insert(csrPos int, insStr []Cell) {
	// expand buffer
	B.Buffer = append(B.Buffer, insStr...)

	// shift original string to make area
	copy(B.Buffer[csrPos+len(insStr):], B.Buffer[csrPos:])

	// insert insStr
	copy(B.Buffer[csrPos:csrPos+len(insStr)], insStr)

	u := &_Undo{
		pos: csrPos,
		del: len(insStr),
	}
	B.undoes = append(B.undoes, u)
}

func mojis2cells(mojis []Moji) []Cell {
	cells := make([]Cell, 0, len(mojis))
	for _, m := range mojis {
		cells = append(cells, Cell{Moji: m})
	}
	return cells
}

func (B *Buffer) insertString(pos int, s string) _Range {
	list := mojis2cells(StringToMoji(s))
	B.insert(pos, list)
	return _Range(list)
}

// InsertString inserts string s at pos (Do not update screen)
// It returns the count of runes
func (B *Buffer) InsertString(pos int, s string) int {
	return len(B.insertString(pos, s))
}

// Delete remove Buffer[pos:pos+n].
// It returns the width to clear the end of line.
// It does not update screen.
func (B *Buffer) Delete(pos int, n int) WidthT {
	if n <= 0 || len(B.Buffer) < pos+n {
		return 0
	}
	u := &_Undo{
		pos:  pos,
		text: cell2string(B.Buffer[pos : pos+n]),
	}
	B.undoes = append(B.undoes, u)

	delw := B.GetWidthBetween(pos, pos+n)
	copy(B.Buffer[pos:], B.Buffer[pos+n:])
	B.Buffer = B.Buffer[:len(B.Buffer)-n]
	return delw
}

// ResetViewStart set ViewStart the new value which should be.
// It does not update screen.
func (B *Buffer) ResetViewStart() {
	B.ViewStart = 0
	w := WidthT(0)
	for i := 0; i <= B.Cursor && i < len(B.Buffer); i++ {
		w += B.Buffer[i].Moji.Width()
		for w >= B.ViewWidth() {
			if B.ViewStart >= len(B.Buffer) {
				// When standard output is redirected.
				return
			}
			w -= B.Buffer[B.ViewStart].Moji.Width()
			B.ViewStart++
		}
	}
}

// GetWidthBetween returns the width between start and end
func (B *Buffer) GetWidthBetween(from int, to int) WidthT {
	return _Range(B.Buffer[from:to]).Width()
}

// SubString returns the readline string between start and end
func (B *Buffer) SubString(start, end int) string {
	return cell2string(B.Buffer[start:end])
}

func (B Buffer) String() string {
	return cell2string(B.Buffer)
}

// Delimiters means the quationmarks. The whitespace enclosed by them
// are not treat as parameters separator.
var Delimiters = "\"'"

// CurrentWordTop returns the position of the current word the cursor exists
func (B *Buffer) CurrentWordTop() (wordTop int) {
	wordTop = -1
	quotedchar := '\000'
	for i, m := range B.Buffer[:B.Cursor] {
		if ch, ok := moji.MojiToRune(m.Moji); ok {
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
		return B.Cursor
	}
	return wordTop
}

// CurrentWord returns the current word the cursor exists and word's position
func (B *Buffer) CurrentWord() (string, int) {
	start := B.CurrentWordTop()
	return B.SubString(start, B.Cursor), start
}

func (B *Buffer) joinUndo() {
	if len(B.undoes) < 2 {
		return
	}
	u1 := B.undoes[len(B.undoes)-2]
	u2 := B.undoes[len(B.undoes)-1]
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
	B.undoes = B.undoes[:len(B.undoes)-1]
}

func (B *Buffer) startChangeWidthEventLoop(_lastw int, getResizeEvent func() (int, int, bool)) {
	go func(lastw int) {
		for {
			w, _, ok := getResizeEvent()
			if !ok {
				break
			}
			if lastw != w {
				mu.Lock()
				B.termWidth = w
				B.RepaintAfterPrompt()
				mu.Unlock()
				lastw = w
			}
		}
	}(_lastw)
}

// GetKey reads one-key from Tty.
func (B *Buffer) GetKey() (string, error) {
	B.Out.Flush()
	return getKey(B.Tty)
}

func (B *Buffer) eraseline() {
	io.WriteString(B.Out, "\x1B[0K")
}

type _Range []Cell

func (B *Buffer) puts(s []Cell) _Range {
	defaultColor := ColorSequence(B.refreshColor())
	color := ColorSequence(-1)
	for _, ch := range s {
		if ch.color != color {
			color = ch.color
			color.WriteTo(B.Out)
		}
		ch.Moji.PrintTo(B.Out)
	}
	if color != defaultColor {
		defaultColor.WriteTo(B.Out)
	}
	return _Range(s)
}

func (s _Range) Width() (w WidthT) {
	for _, ch := range s {
		w += ch.Moji.Width()
	}
	return
}

func cell2string(m []Cell) string {
	var buffer strings.Builder
	for _, m1 := range m {
		m1.Moji.WriteTo(&buffer)
	}
	return buffer.String()
}

type WidthT = moji.WidthT

type Moji = moji.Moji

func StringToMoji(s string) []Moji {
	return moji.StringToMoji(s)
}

func GetStringWidth(s string) WidthT {
	w, _ := moji.MojiWidthAndCountInString(s)
	return w
}

func MojiCountInString(s string) int {
	_, c := moji.MojiWidthAndCountInString(s)
	return c
}

func ResetCharWidth() {
	moji.ResetCharWidth()
}

func SetCharWidth(c rune, width int) {
	moji.SetCharWidth(c, width)
}

func EnableSurrogatePair(value bool) {
	moji.SurrogatePairOk = value
}

func IsSurrogatePairEnabled() bool {
	return moji.SurrogatePairOk
}
