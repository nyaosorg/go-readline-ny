package readline

import (
	"fmt"
	"unicode/utf8"

	"github.com/nyaosorg/go-readline-ny/internal/moji"
)

type _MonoChrome struct{}

func (_MonoChrome) Init() ColorSequence {
	return DefaultForeGroundColor
}

func (_MonoChrome) Next(rune) ColorSequence {
	return DefaultForeGroundColor
}

func (buf *Buffer) RefreshColor() ColorSequence {
	if buf.Coloring == nil {
		buf.Coloring = _MonoChrome{}
	}
	defaultColor := buf.Coloring.Init()
	position := int16(0)
	for i, cell := range buf.Buffer {
		buf.Buffer[i].position = position
		if codepoint, ok := moji.MojiToRune(cell.Moji); ok {
			buf.Buffer[i].color = ColorSequence(buf.Coloring.Next(codepoint))
		} else {
			buf.Buffer[i].color = ColorSequence(buf.Coloring.Next(utf8.RuneError))
		}
		position += int16(cell.Moji.Width())
	}
	return defaultColor
}

// InsertAndRepaint inserts str and repaint the editline.
func (buf *Buffer) InsertAndRepaint(str string) {
	buf.ReplaceAndRepaint(buf.Cursor, str)
}

// GotoHead move screen-cursor to the top of the viewarea.
// It should be called before text is changed.
func (B *Buffer) GotoHead() {
	fmt.Fprintf(B, "\x1B[%dG", B.topColumn+1)
}

func (B *Buffer) repaint() {
	all, left := B.view2()
	B.GotoHead()
	B.puts(all)
	B.eraseline()
	B.GotoHead()
	B.puts(left)
}

// DrawFromHead draw all text in viewarea and
// move screen-cursor to the position where it should be.
func (B *Buffer) DrawFromHead() {
	B.repaint()
}

// ReplaceAndRepaint replaces the string between `pos` and cursor's position to `str`
func (buf *Buffer) ReplaceAndRepaint(pos int, str string) {
	// Replace Buffer
	buf.Delete(pos, buf.Cursor-pos)

	// Define ViewStart , Cursor
	buf.Cursor = pos + buf.InsertString(pos, str)
	buf.joinUndo() // merge delete and insert
	buf.ResetViewStart()
	buf.repaint()
}

// RepaintAfterPrompt repaints the all characters in the editline except for prompt.
func (B *Buffer) RepaintAfterPrompt() {
	B.ResetViewStart()
	B.repaint()
}

// RepaintAll repaints the all characters in the editline including prompt.
func (B *Buffer) RepaintAll() {
	B.Out.Flush()
	B.topColumn, _ = B.Prompt()
	B.repaint()
}
