package readline

import (
	"context"
	"io"
	"strings"

	"github.com/atotto/clipboard"

	"github.com/nyaosorg/go-readline-ny/internal/moji"
)

// Gommand is the implement of KeyFuncT which has a name and a function
type Gommand struct {
	Name string
	Func func(ctx context.Context, buffer *Buffer) Result
}

// Deprecate: use Gommand instead
type KeyGoFuncT = Gommand

// String returns Gommand's name
func (K Gommand) String() string {
	return K.Name
}

// Call calls the function the receiver contains
func (K *Gommand) Call(ctx context.Context, buffer *Buffer) Result {
	if K.Func == nil {
		return CONTINUE
	}
	return K.Func(ctx, buffer)
}

var CmdAcceptLine = &Gommand{
	Name: F_ACCEPT_LINE,
	Func: func(ctx context.Context, this *Buffer) Result { // Ctrl-M
		return ENTER
	},
}

var CmdInterrupt = &Gommand{
	Name: F_INTR,
	Func: func(ctx context.Context, this *Buffer) Result { // Ctrl-C
		this.Buffer = this.Buffer[:0]
		this.Cursor = 0
		this.ViewStart = 0
		this.undoes = nil
		return INTR
	},
}

var CmdBeginningOfLine = &Gommand{
	Name: F_BEGINNING_OF_LINE,
	Func: func(ctx context.Context, this *Buffer) Result { // Ctrl-A
		this.Cursor = 0
		this.ViewStart = 0
		this.repaint()
		return CONTINUE
	},
}

var CmdBackwardChar = &Gommand{
	Name: F_BACKWARD_CHAR,
	Func: func(ctx context.Context, this *Buffer) Result { // Ctrl-B
		if this.Cursor <= 0 {
			return CONTINUE
		}
		this.Cursor--
		if this.Cursor < this.ViewStart {
			this.ViewStart--
		}
		this.repaint()
		return CONTINUE
	},
}

var CmdEndOfLine = &Gommand{
	Name: F_END_OF_LINE,
	Func: func(ctx context.Context, this *Buffer) Result { // Ctrl-E
		allength := this.GetWidthBetween(this.ViewStart, len(this.Buffer))
		if allength < this.ViewWidth() {
			this.puts(this.Buffer[this.Cursor:])
			this.Cursor = len(this.Buffer)
		} else {
			this.GotoHead()
			this.ViewStart = len(this.Buffer) - 1
			w := this.Buffer[this.ViewStart].Moji.Width()
			for {
				if this.ViewStart <= 0 {
					break
				}
				_w := w + this.Buffer[this.ViewStart-1].Moji.Width()
				if _w >= this.ViewWidth() {
					break
				}
				w = _w
				this.ViewStart--
			}
			this.puts(this.Buffer[this.ViewStart:])
			this.Cursor = len(this.Buffer)
		}
		this.eraseline()
		return CONTINUE
	},
}

var CmdForwardChar = &Gommand{
	Name: F_FORWARD_CHAR,
	Func: func(ctx context.Context, this *Buffer) Result { // Ctrl-F
		if this.Cursor >= len(this.Buffer) {
			return CONTINUE
		}
		w := this.GetWidthBetween(this.ViewStart, this.Cursor+1)
		if w < this.ViewWidth() {
			// No Scroll
			this.puts(this.Buffer[this.Cursor : this.Cursor+1])
		} else {
			// Right Scroll
			this.GotoHead()
			if this.Buffer[this.Cursor].Moji.Width() > this.Buffer[this.ViewStart].Moji.Width() {
				this.ViewStart++
			}
			this.ViewStart++
			this.puts(this.Buffer[this.ViewStart : this.Cursor+1])
			this.eraseline()
		}
		this.Cursor++
		return CONTINUE
	},
}

var CmdBackwardDeleteChar = &Gommand{
	Name: F_BACKWARD_DELETE_CHAR,
	Func: func(ctx context.Context, this *Buffer) Result { // Backspace
		if this.Cursor > 0 {
			this.Cursor--
			this.Delete(this.Cursor, 1)
			if this.Cursor < this.ViewStart {
				this.ViewStart = this.Cursor
			}
			this.repaint()
		}
		return CONTINUE
	},
}

var CmdDeleteChar = &Gommand{
	Name: F_DELETE_CHAR,
	Func: func(ctx context.Context, this *Buffer) Result { // Del
		this.Delete(this.Cursor, 1)
		this.repaint()
		return CONTINUE
	},
}

var CmdDeleteOrAbort = &Gommand{
	Name: F_DELETE_OR_ABORT,
	Func: func(ctx context.Context, this *Buffer) Result { // Ctrl-D
		if len(this.Buffer) > 0 {
			return CmdDeleteChar.Func(ctx, this)
		}
		return ABORT
	},
}

func mojiAndStringToString(m Moji, s string) string {
	var buffer strings.Builder
	m.WriteTo(&buffer)
	buffer.WriteString(s)
	return buffer.String()
}

func keyFuncInsertSelf(ctx context.Context, this *Buffer, keys string) Result {
	if len(keys) == 2 && keys[0] == '\x1B' { // for AltGr-shift
		keys = keys[1:]
	}
	if moji.AreZeroWidthJoin(keys) && this.Cursor > 0 {
		this.pending = mojiAndStringToString(
			this.Buffer[this.Cursor-1].Moji,
			keys)
		return CmdBackwardDeleteChar.Func(ctx, this)
	} else if (moji.AreVariationSelectorLike(keys) || moji.AreEmojiModifier(keys)) && this.Cursor > 0 {
		baseMoji := this.Buffer[this.Cursor-1].Moji
		CmdBackwardDeleteChar.Func(ctx, this)
		keys = mojiAndStringToString(baseMoji, keys)
	} else if len(this.pending) > 0 {
		keys = this.pending + keys
		this.pending = ""
	}

	mojis := this.insertString(this.Cursor, keys)
	lenMoji := len(mojis)

	w := this.GetWidthBetween(this.ViewStart, this.Cursor)
	w1 := mojis.Width()
	this.Cursor += lenMoji
	if w+w1 >= this.ViewWidth() {
		// scroll left
		this.ResetViewStart()
	}
	this.repaint()
	return CONTINUE
}

var CmdKillLine = &Gommand{
	Name: F_KILL_LINE,
	Func: func(ctx context.Context, this *Buffer) Result {
		clipboard.WriteAll(this.SubString(this.Cursor, len(this.Buffer)))

		this.eraseline()
		u := &_Undo{
			pos:  this.Cursor,
			text: cell2string(this.Buffer[this.Cursor:]),
		}
		this.undoes = append(this.undoes, u)
		this.Buffer = this.Buffer[:this.Cursor]
		return CONTINUE
	},
}

var CmdKillWholeLine = &Gommand{
	Name: F_KILL_WHOLE_LINE,
	Func: func(ctx context.Context, this *Buffer) Result {
		u := &_Undo{
			pos:  0,
			text: cell2string(this.Buffer),
		}
		this.undoes = append(this.undoes, u)
		this.GotoHead()
		this.eraseline()
		this.Buffer = this.Buffer[:0]
		this.Cursor = 0
		this.ViewStart = 0
		return CONTINUE
	},
}

var CmdUnixWordRubout = &Gommand{
	Name: F_UNIX_WORD_RUBOUT,
	Func: func(ctx context.Context, this *Buffer) Result {
		orgCursorPos := this.Cursor
		for this.Cursor > 0 && moji.IsSpaceMoji(this.Buffer[this.Cursor-1].Moji) {
			this.Cursor--
		}
		newCursorPos := this.CurrentWordTop()
		clipboard.WriteAll(this.SubString(newCursorPos, orgCursorPos))
		this.Delete(newCursorPos, orgCursorPos-newCursorPos)
		this.Cursor = newCursorPos
		if newCursorPos-this.ViewStart < 2 {
			this.ResetViewStart()
		}
		this.repaint()
		return CONTINUE
	},
}

var CmdUnixLineDiscard = &Gommand{
	Name: F_UNIX_LINE_DISCARD,
	Func: func(ctx context.Context, this *Buffer) Result {
		clipboard.WriteAll(this.SubString(0, this.Cursor))
		this.Delete(0, this.Cursor)
		this.Cursor = 0
		this.ViewStart = 0
		this.repaint()
		return CONTINUE
	},
}

var CmdClearScreen = &Gommand{
	Name: F_CLEAR_SCREEN,
	Func: func(ctx context.Context, this *Buffer) Result {
		io.WriteString(this.Out, "\x1B[1;1H\x1B[2J")
		this.RepaintAll()
		return CONTINUE
	},
}

var CmdRepaintOnNewline = &Gommand{
	Name: F_REPAINT_ON_NEWLINE,
	Func: func(ctx context.Context, this *Buffer) Result {
		this.Out.WriteByte('\n')
		this.RepaintAll()
		return CONTINUE
	},
}

var CmdQuotedInsert = &Gommand{
	Name: F_QUOTED_INSERT,
	Func: func(ctx context.Context, this *Buffer) Result {
		io.WriteString(this.Out, ansiCursorOn)
		defer io.WriteString(this.Out, ansiCursorOff)

		this.Out.Flush()
		if key, err := this.GetKey(); err == nil {
			return keyFuncInsertSelf(ctx, this, key)
		}
		return CONTINUE
	},
}

var CmdYank = &Gommand{
	Name: F_YANK,
	Func: func(ctx context.Context, this *Buffer) Result {
		text, err := clipboard.ReadAll()
		if err != nil {
			return CONTINUE
		}
		text = strings.TrimRight(text, "\r\n\000")
		this.InsertAndRepaint(text)
		return CONTINUE
	},
}

var CmdYankWithQuote = &Gommand{
	Name: F_YANK_WITH_QUOTE,
	Func: func(ctx context.Context, this *Buffer) Result {
		text, err := clipboard.ReadAll()
		if err != nil {
			return CONTINUE
		}
		if strings.IndexRune(text, ' ') >= 0 &&
			!strings.HasPrefix(text, `"`) {

			text = `"` + strings.Replace(text, `"`, `""`, -1) + `"`
			text = strings.Replace(text, "\r\n", "\"\r\n\"", -1)
		}
		this.InsertAndRepaint(text)
		return CONTINUE
	},
}

var CmdSwapChar = &Gommand{
	Name: F_SWAPCHAR,
	Func: func(ctx context.Context, this *Buffer) Result {
		if len(this.Buffer) == this.Cursor {
			if this.Cursor < 2 {
				return CONTINUE
			}
			u := &_Undo{
				pos:  this.Cursor,
				del:  2,
				text: cell2string(this.Buffer[this.Cursor-2 : this.Cursor]),
			}
			this.undoes = append(this.undoes, u)
			this.Buffer[this.Cursor-2], this.Buffer[this.Cursor-1] = this.Buffer[this.Cursor-1], this.Buffer[this.Cursor-2]

			this.GotoHead()
			this.puts(this.Buffer[this.ViewStart:this.Cursor])
		} else {
			if this.Cursor < 1 {
				return CONTINUE
			}
			u := &_Undo{
				pos:  this.Cursor - 1,
				del:  2,
				text: cell2string(this.Buffer[this.Cursor-1 : this.Cursor+1]),
			}
			this.undoes = append(this.undoes, u)

			w := this.GetWidthBetween(this.ViewStart, this.Cursor+1)
			this.Buffer[this.Cursor-1], this.Buffer[this.Cursor] = this.Buffer[this.Cursor], this.Buffer[this.Cursor-1]
			this.GotoHead()
			if w >= this.ViewWidth() {
				this.ViewStart++
			}
			this.Cursor++
			this.puts(this.Buffer[this.ViewStart:this.Cursor])
		}
		return CONTINUE
	},
}

var CmdBackwardWord = &Gommand{
	Name: F_BACKWARD_WORD,
	Func: func(ctx context.Context, this *Buffer) Result {
		newPos := this.Cursor
		for newPos > 0 && moji.IsSpaceMoji(this.Buffer[newPos-1].Moji) {
			newPos--
		}
		for newPos > 0 && !moji.IsSpaceMoji(this.Buffer[newPos-1].Moji) {
			newPos--
		}
		if newPos < this.ViewStart {
			this.ViewStart = newPos
		}
		this.Cursor = newPos
		this.repaint()
		return CONTINUE
	},
}

var CmdForwardWord = &Gommand{
	Name: F_FORWARD_WORD,
	Func: func(ctx context.Context, this *Buffer) Result {
		newPos := this.Cursor
		for newPos < len(this.Buffer) && !moji.IsSpaceMoji(this.Buffer[newPos].Moji) {
			newPos++
		}
		for newPos < len(this.Buffer) && moji.IsSpaceMoji(this.Buffer[newPos].Moji) {
			newPos++
		}
		w := this.GetWidthBetween(this.ViewStart, newPos)
		if w < this.ViewWidth() {
			this.puts(this.Buffer[this.Cursor:newPos])
			this.Cursor = newPos
		} else {
			this.Cursor = newPos
			this.ResetViewStart()
			this.repaint()
		}
		return CONTINUE
	},
}

var CmdUndo = &Gommand{
	Name: F_UNDO,
	Func: func(ctx context.Context, this *Buffer) Result {
		if len(this.undoes) <= 0 {
			io.WriteString(this.Out, "\a")
			return CONTINUE
		}
		u := this.undoes[len(this.undoes)-1]
		this.undoes = this.undoes[:len(this.undoes)-1]

		this.GotoHead()
		if u.del > 0 {
			copy(this.Buffer[u.pos:], this.Buffer[u.pos+u.del:])
			this.Buffer = this.Buffer[:len(this.Buffer)-u.del]
		}
		if u.text != "" {
			t := mojis2cells(StringToMoji(u.text))
			// widen buffer
			this.Buffer = append(this.Buffer, t...)
			// make area
			copy(this.Buffer[u.pos+len(t):], this.Buffer[u.pos:])
			copy(this.Buffer[u.pos:], t)
			this.Cursor = u.pos + len(t)
		} else {
			this.Cursor = u.pos
		}
		this.ResetViewStart()
		this.repaint()
		return CONTINUE
	},
}