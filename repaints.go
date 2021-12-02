package readline

import (
	"unicode/utf8"
)

type _MonoChrome struct{}

func (_MonoChrome) Init() {}

func (_MonoChrome) Next(rune) int {
	return White
}

func (buf *Buffer) RefreshColor() {
	if buf.Coloring == nil {
		buf.Coloring = _MonoChrome{}
	}
	buf.Coloring.Init()
	position := int16(0)
	for i, cell := range buf.Buffer {
		buf.Buffer[i].position = position
		if codepoint, ok := moji2rune(cell.Moji); ok {
			buf.Buffer[i].color = _PackedColorCode(buf.Coloring.Next(codepoint))
		} else {
			buf.Buffer[i].color = _PackedColorCode(buf.Coloring.Next(utf8.RuneError))
		}
		position += int16(cell.Moji.Width())
	}
}

// InsertAndRepaint inserts str and repaint the editline.
func (buf *Buffer) InsertAndRepaint(str string) {
	buf.ReplaceAndRepaint(buf.Cursor, str)
}

// GotoHead move screen-cursor to the top of the viewarea.
// It should be called before text is changed.
func (buf *Buffer) GotoHead() {
	buf.backspace(_Range(buf.Buffer[buf.ViewStart:buf.Cursor]).Width())
}

// DrawFromHead draw all text in viewarea and
// move screen-cursor to the position where it should be.
func (buf *Buffer) DrawFromHead() {
	// Repaint
	view, _, right := buf.view3()
	buf.puts(view)

	// Move to cursor position
	buf.eraseline()
	buf.backspace(right.Width())
}

// ReplaceAndRepaint replaces the string between `pos` and cursor's position to `str`
func (buf *Buffer) ReplaceAndRepaint(pos int, str string) {
	buf.GotoHead()

	// Replace Buffer
	buf.Delete(pos, buf.Cursor-pos)

	// Define ViewStart , Cursor
	buf.Cursor = pos + buf.InsertString(pos, str)

	buf.joinUndo() // merge delete and insert

	buf.ResetViewStart()

	buf.DrawFromHead()
}

// Repaint buffer[pos:] + " \b"*del but do not rewind cursor position
func (buf *Buffer) Repaint(pos int, del WidthT) {
	buf.GetWidthBetween(buf.ViewStart, pos)

	view := buf.view()
	bs := buf.puts(view[pos-buf.ViewStart:]).Width()

	buf.eraseline()
	if del > 0 {
		buf.backspace(bs)
	} else {
		// for readline_keyfunc.go: KeyFuncInsertSelf()
		buf.backspace(bs + del)
	}
}

// RepaintAfterPrompt repaints the all characters in the editline except for prompt.
func (buf *Buffer) RepaintAfterPrompt() {
	buf.ResetViewStart()
	buf.puts(buf.Buffer[buf.ViewStart:buf.Cursor])
	buf.Repaint(buf.Cursor, 0)
}

// RepaintAll repaints the all characters in the editline including prompt.
func (buf *Buffer) RepaintAll() {
	buf.Out.Flush()
	buf.topColumn, _ = buf.Prompt()
	buf.RepaintAfterPrompt()
}
