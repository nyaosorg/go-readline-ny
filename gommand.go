package readline

import (
	"context"
	"io"
	"strings"

	"github.com/nyaosorg/go-readline-ny/moji"
)

// GoCommand is the implement of Command which has a name and a function
type GoCommand struct {
	Name string
	Func func(ctx context.Context, buffer *Buffer) Result
}

// String returns GoCommand's name
func (K GoCommand) String() string {
	return K.Name
}

// Call calls the function the receiver contains
func (K *GoCommand) Call(ctx context.Context, buffer *Buffer) Result {
	if K.Func == nil {
		return CONTINUE
	}
	return K.Func(ctx, buffer)
}

var NameToFunc = map[string]Command{}

// NewGoCommand creates an instance of GoCommand with a name and an function.
func NewGoCommand(name string, f func(context.Context, *Buffer) Result) *GoCommand {
	instance := &GoCommand{
		Name: name,
		Func: f,
	}
	NameToFunc[name] = instance
	return instance
}

// CmdAcceptLine is the command that ends the editing successfully (for Ctrl-M)
// Then the method `(*Editor) ReadLine()` returns nil as error.
var CmdAcceptLine = NewGoCommand("ACCEPT_LINE", cmdAcceptLine)

func cmdAcceptLine(ctx context.Context, this *Buffer) Result {
	return ENTER
}

// CmdInterrupt is the command that aborts the editing (for Ctrl-C)
// Then the method `(*Editor) ReadLine()` returns `readline.CtrlC` as error.
var CmdInterrupt = NewGoCommand("INTR", cmdInterrupt)

func cmdInterrupt(ctx context.Context, this *Buffer) Result {
	this.Buffer = this.Buffer[:0]
	this.Cursor = 0
	this.ViewStart = 0
	this.undoes = nil
	return INTR
}

// CmdBeginningOfLine is the command that moves cursor to the top of line (for Ctrl-A)
var CmdBeginningOfLine = NewGoCommand("BEGINNING_OF_LINE", cmdBeginningOfLine)

func cmdBeginningOfLine(ctx context.Context, this *Buffer) Result {
	this.Cursor = 0
	this.ViewStart = 0
	this.repaint()
	return CONTINUE
}

// CmdBackwardChar is the command that moves cursor to the previous character (for Ctrl-B)
var CmdBackwardChar = NewGoCommand("BACKWARD_CHAR", cmdBackwardChar)

func cmdBackwardChar(ctx context.Context, this *Buffer) Result {
	if this.Cursor <= 0 {
		return CONTINUE
	}
	this.Cursor--
	if this.Cursor < this.ViewStart {
		this.ViewStart--
	}
	this.repaint()
	return CONTINUE
}

// CmdEndOfLine is the command that moves cursor to the end of the line (for Ctrl-E)
var CmdEndOfLine = NewGoCommand("END_OF_LINE", cmdEndOfLine)

func cmdEndOfLine(ctx context.Context, this *Buffer) Result {
	allength := this.GetWidthBetween(this.ViewStart, len(this.Buffer))
	if allength <= this.ViewWidth() {
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
			if _w > this.ViewWidth() {
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
}

// CmdForwardChar is the command that moves cursor to the next character (for Ctrl-F)
var CmdForwardChar = NewGoCommand("FORWARD_CHAR", cmdForwardChar)

func cmdForwardChar(ctx context.Context, this *Buffer) Result {
	if this.Cursor >= len(this.Buffer) {
		if len(this.suffix) > 0 {
			var s strings.Builder
			for _, c := range this.suffix {
				c.WriteTo(&s)
			}
			this.InsertAndRepaint(s.String())
		}
		return CONTINUE
	}
	w := this.GetWidthBetween(this.ViewStart, this.Cursor+1)
	if w <= this.ViewWidth() {
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
}

// CmdBackwardDeleteChar is the command that deletes a character on the leftside of cursor (for Backspace)
var CmdBackwardDeleteChar = NewGoCommand("BACKWARD_DELETE_CHAR", cmdBackwardDeleteChar)

func cmdBackwardDeleteChar(ctx context.Context, this *Buffer) Result {
	if this.Cursor > 0 {
		this.Cursor--
		this.Delete(this.Cursor, 1)
		if this.Cursor < this.ViewStart {
			this.ViewStart = this.Cursor
		}
		this.repaint()
	}
	return CONTINUE
}

// CmdDeleteChar is the command that delete a character on the cursor (for Del)
var CmdDeleteChar = NewGoCommand("DELETE_CHAR", cmdDeleteChar)

func cmdDeleteChar(ctx context.Context, this *Buffer) Result {
	this.Delete(this.Cursor, 1)
	this.repaint()
	return CONTINUE
}

// CmdDeleteOrAbort is the command that deletes a character on the cursor when the line has any characters. Otherwise quit the editing (for Ctrl-D)
var CmdDeleteOrAbort = NewGoCommand("DELETE_OR_ABORT", cmdDeleteOrAbort)

func cmdDeleteOrAbort(ctx context.Context, this *Buffer) Result {
	if len(this.Buffer) > 0 {
		return CmdDeleteChar.Func(ctx, this)
	}
	return ABORT
}

func mojiAndStringToString(m Moji, s string) string {
	var buffer strings.Builder
	m.WriteTo(&buffer)
	buffer.WriteString(s)
	return buffer.String()
}

// SelfInserter is a command object that inserts the string itself.
type SelfInserter string

func (s SelfInserter) String() string {
	return string(s)
}

func (s SelfInserter) Call(ctx context.Context, this *Buffer) Result {
	keys := string(s)
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
	if w+w1 > this.ViewWidth() {
		// scroll left
		this.ResetViewStart()
	}
	this.repaint()
	return CONTINUE
}

// CmdKillLine is the command that removes text from cursor to the end of the line (for Ctrl-K)
var CmdKillLine = NewGoCommand("KILL_LINE", cmdKillLine)

func cmdKillLine(ctx context.Context, this *Buffer) Result {
	this.Clipboard.Write(this.SubString(this.Cursor, len(this.Buffer)))

	this.eraseline()
	u := &_Undo{
		pos:  this.Cursor,
		text: cell2string(this.Buffer[this.Cursor:]),
	}
	this.undoes = append(this.undoes, u)
	this.Buffer = this.Buffer[:this.Cursor]
	return CONTINUE
}

// CmdKillWholeLine is the command that removes the whole characters of the line.
var CmdKillWholeLine = NewGoCommand("KILL_WHOLE_LINE", cmdKillWholeLine)

func cmdKillWholeLine(ctx context.Context, this *Buffer) Result {
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
}

// CmdUnixWordRubout is the command that removes the current word (for Ctrl-W)
var CmdUnixWordRubout = NewGoCommand("UNIX_WORD_RUBOUT", cmdUnixWordRubout)

func cmdUnixWordRubout(ctx context.Context, this *Buffer) Result {
	orgCursorPos := this.Cursor
	for this.Cursor > 0 && moji.IsSpaceMoji(this.Buffer[this.Cursor-1].Moji) {
		this.Cursor--
	}
	newCursorPos := this.CurrentWordTop()
	this.Clipboard.Write(this.SubString(newCursorPos, orgCursorPos))
	this.Delete(newCursorPos, orgCursorPos-newCursorPos)
	this.Cursor = newCursorPos
	if newCursorPos-this.ViewStart < 2 {
		this.ResetViewStart()
	}
	this.repaint()
	return CONTINUE
}

// CmdUnixLineDiscard is the command that removes from the top of the line until cursor. (for Ctrl-U)
var CmdUnixLineDiscard = NewGoCommand("UNIX_LINE_DISCARD", cmdUnixLineDiscard)

func cmdUnixLineDiscard(ctx context.Context, this *Buffer) Result {
	this.Clipboard.Write(this.SubString(0, this.Cursor))
	this.Delete(0, this.Cursor)
	this.Cursor = 0
	this.ViewStart = 0
	this.repaint()
	return CONTINUE
}

// CmdClearScreen is the command that clears the screen and repaints the line (for Ctrl-L)
var CmdClearScreen = NewGoCommand("CLEAR_SCREEN", cmdClearScreen)

func cmdClearScreen(ctx context.Context, this *Buffer) Result {
	io.WriteString(this.Out, "\x1B[1;1H\x1B[2J")
	this.RepaintAll()
	return CONTINUE
}

var CmdRepaintOnNewline = NewGoCommand("REPAINT_ON_NEWLINE", cmdRepaintOnNewline)

func cmdRepaintOnNewline(ctx context.Context, this *Buffer) Result {
	this.Out.WriteByte('\n')
	this.RepaintAll()
	return CONTINUE
}

// CmdQuotedInsert is the command that inserts the next typed character itself even if it is a control-character (for Ctrl-V or Ctrl-Q)
var CmdQuotedInsert = NewGoCommand("QUOTED_INSERT", cmdQuotedInsert)

func cmdQuotedInsert(ctx context.Context, this *Buffer) Result {
	io.WriteString(this.Out, ansiCursorOn)
	defer io.WriteString(this.Out, ansiCursorOff)

	if key, err := this.GetKey(); err == nil {
		return SelfInserter(key).Call(ctx, this)
	}
	return CONTINUE
}

// CmdYank is the command that inserts the string in the clipboard (for Ctrl-Y)
var CmdYank = NewGoCommand("YANK", cmdYank)

func cmdYank(ctx context.Context, this *Buffer) Result {
	text, err := this.Clipboard.Read()
	if err != nil {
		return CONTINUE
	}
	text = strings.TrimRight(text, "\r\n\000")
	this.InsertAndRepaint(text)
	return CONTINUE
}

// CmdYankWithQuote is the command that inserts the string in the clipboard and enclose it with double quoatations.
var CmdYankWithQuote = NewGoCommand("YANK_WITH_QUOTE", cmdYankWithQuote)

func cmdYankWithQuote(ctx context.Context, this *Buffer) Result {
	text, err := this.Clipboard.Read()
	if err != nil {
		return CONTINUE
	}
	if strings.ContainsRune(text, ' ') &&
		!strings.HasPrefix(text, `"`) {

		text = `"` + strings.Replace(text, `"`, `""`, -1) + `"`
		text = strings.Replace(text, "\r\n", "\"\r\n\"", -1)
	}
	this.InsertAndRepaint(text)
	return CONTINUE
}

// CmdSwapChar is the command that swaps the character on the cursor and on on the leftside of the cursor (for Ctrl-T)
var CmdSwapChar = NewGoCommand("SWAPCHAR", cmdSwapChar)

func cmdSwapChar(ctx context.Context, this *Buffer) Result {
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
		if w > this.ViewWidth() {
			this.ViewStart++
		}
		this.Cursor++
		this.puts(this.Buffer[this.ViewStart:this.Cursor])
	}
	return CONTINUE
}

// CmdBackwardWord is the command that moves cursor to the top of the previous word (for M-B)
var CmdBackwardWord = NewGoCommand("BACKWARD_WORD", cmdBackwardWord)

func cmdBackwardWord(ctx context.Context, this *Buffer) Result {
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
}

// CmdForwardWord is the command that move cursor to the top of the next word (for M-F)
var CmdForwardWord = NewGoCommand("FORWARD_WORD", cmdForwardWord)

func cmdForwardWord(ctx context.Context, this *Buffer) Result {
	newPos := this.Cursor
	for newPos < len(this.Buffer) && !moji.IsSpaceMoji(this.Buffer[newPos].Moji) {
		newPos++
	}
	for newPos < len(this.Buffer) && moji.IsSpaceMoji(this.Buffer[newPos].Moji) {
		newPos++
	}
	w := this.GetWidthBetween(this.ViewStart, newPos)
	if w <= this.ViewWidth() {
		this.puts(this.Buffer[this.Cursor:newPos])
		this.Cursor = newPos
	} else {
		this.Cursor = newPos
		this.ResetViewStart()
		this.repaint()
	}
	return CONTINUE
}

// CmdUndo is the command that executes undo (for Ctrl-U)
var CmdUndo = NewGoCommand("UNDO", cmdUndo)

func cmdUndo(ctx context.Context, this *Buffer) Result {
	if len(this.undoes) <= 0 {
		io.WriteString(this.Out, "\a")
		return CONTINUE
	}
	u := this.undoes[len(this.undoes)-1]
	this.undoes = this.undoes[:len(this.undoes)-1]

	this.GotoHead()
	if u.del > 0 {
		this.Buffer = sliceDelete(this.Buffer, u.pos, u.del)
	}
	if u.text != "" {
		cells := mojis2cells(StringToMoji(u.text))
		this.Buffer = sliceInsert(this.Buffer, u.pos, cells...)
		this.Cursor = u.pos + len(cells)
	} else {
		this.Cursor = u.pos
	}
	this.ResetViewStart()
	this.repaint()
	return CONTINUE
}
