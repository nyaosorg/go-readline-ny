package readline

import (
	"fmt"
	"io"
	"strings"

	"github.com/nyaosorg/go-readline-ny/moji"
)

const (
	CursorPositionDummyRune = '\uE000'
)

func (B *Buffer) refreshColor() ColorInterface {
	var ci interface {
		Init() ColorInterface
		Next(rune) ColorInterface
	}
	if B.Coloring != nil {
		ci = &colorBridge{base: B.Coloring}
	} else {
		str := B.String()
		cursorPos := 0
		for i := 0; i < B.Cursor; i++ {
			cursorPos += B.Buffer[i].Moji.Len()
		}
		if B.memoHighlightSource == str && B.memoHighlightResult != nil {
			ci = B.memoHighlightResult
		} else {
			result := HighlightToColorSequence(str, B.ResetColor, B.DefaultColor, B.Highlight, -1-cursorPos)
			B.memoHighlightSource = str
			B.memoHighlightResult = result
			ci = result
		}
	}
	var defaultColor ColorInterface = ci.Init()
	position := int16(0)
	var tmpbuf strings.Builder
	for i, cell := range B.Buffer {
		if i == B.Cursor {
			ci.Next(CursorPositionDummyRune)
		}
		B.Buffer[i].position = position
		if tab, ok := cell.Moji.(*moji.Tab); ok {
			tab.SetPosition(position)
			B.Buffer[i].color = ci.Next('\t')
		} else if codepoint, ok := moji.MojiToRune(cell.Moji); ok {
			B.Buffer[i].color = ci.Next(codepoint)
		} else {
			cell.Moji.PrintTo(&tmpbuf)
			var cs ColorInterface
			for _, c := range tmpbuf.String() {
				cs = ci.Next(c)
			}
			B.Buffer[i].color = cs
			tmpbuf.Reset()
		}
		position += int16(cell.Moji.Width())
	}
	if len(B.Buffer) == B.Cursor {
		ci.Next(CursorPositionDummyRune)
	}
	return defaultColor
}

// InsertAndRepaint inserts str and repaint the editline.
func (B *Buffer) InsertAndRepaint(str string) {
	B.ReplaceAndRepaint(B.Cursor, str)
}

// GotoHead move screen-cursor to the top of the viewarea.
// It should be called before text is changed.
func (B *Buffer) GotoHead() {
	fmt.Fprintf(B.Out, "\x1B[%dG", B.topColumn+1)
}

func (B *Buffer) repaint() {
	B.updateSuffix()
	B.repaintWithoutUpdateSuffix()
}

func (B *Buffer) repaintWithoutUpdateSuffix() {
	all, left, w := B.getView2()
	B.GotoHead()
	puts := B.newPrinter()
	puts(all)
	if B.PredictColor[0] != "" && len(B.suffix) > 0 {
		io.WriteString(B.Out, B.PredictColor[0]) // "\x1B[3;22;39m"
		for _, c := range B.getSuffix() {
			w += c.Width()
			if w >= B.ViewWidth() {
				break
			}
			c.PrintTo(B.Out)
		}
		io.WriteString(B.Out, B.PredictColor[1]) // "\x1B[23m"
	}
	B.eraseline()
	B.GotoHead()
	puts(left)
}

// DrawFromHead draw all text in viewarea and
// move screen-cursor to the position where it should be.
func (B *Buffer) DrawFromHead() {
	B.repaint()
}

// ReplaceAndRepaint replaces the string between `pos` and cursor's position to `str`
func (B *Buffer) ReplaceAndRepaint(pos int, str string) {
	// Replace Buffer
	B.Delete(pos, B.Cursor-pos)

	// Define ViewStart , Cursor
	B.Cursor = pos + B.InsertString(pos, str)
	B.joinUndo() // merge delete and insert
	B.ResetViewStart()
	B.repaint()
}

// RepaintAfterPrompt repaints the all characters in the editline except for prompt.
func (B *Buffer) RepaintAfterPrompt() {
	B.ResetViewStart()
	B.repaint()
}

// RepaintAll repaints the all characters in the editline including prompt.
func (B *Buffer) RepaintAll() {
	B.Out.Flush()
	B.topColumn, _ = B.callPromptWriter()
	B.repaint()
}

// RepaintLastLine repaints the last line of the prompt and input-line.
// IMPORTANT: This method requires setting valid Editor.PromptWriter
func (B *Buffer) RepaintLastLine() {
	B.Out.Flush()
	var prompt string
	if B.PromptWriter == nil {
		prompt = "\r> "
	} else {
		var buffer strings.Builder
		buffer.WriteByte('\r')
		B.PromptWriter(&buffer)
		prompt = buffer.String()
		prompt = strings.ReplaceAll(prompt, "\n", "\r")
	}
	B.Out.WriteString(prompt)
	B.repaint()
}
