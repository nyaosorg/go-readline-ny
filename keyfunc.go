package readline

import (
	"context"
	"io"
	"strings"

	"github.com/atotto/clipboard"
)

func keyFuncEnter(ctx context.Context, this *Buffer) Result { // Ctrl-M
	return ENTER
}

func keyFuncIntr(ctx context.Context, this *Buffer) Result { // Ctrl-C
	this.Buffer = this.Buffer[:0]
	this.Cursor = 0
	this.ViewStart = 0
	this.undoes = nil
	return INTR
}

func keyFuncHead(ctx context.Context, this *Buffer) Result { // Ctrl-A
	this.GotoHead()
	this.Cursor = 0
	this.ViewStart = 0
	this.DrawFromHead()
	return CONTINUE
}

func keyFuncBackward(ctx context.Context, this *Buffer) Result { // Ctrl-B
	if this.Cursor <= 0 {
		return CONTINUE
	}
	this.Cursor--
	if this.Cursor < this.ViewStart {
		this.ViewStart--
		this.DrawFromHead()
	} else {
		this.backspace(this.Buffer[this.Cursor].Width())
	}
	return CONTINUE
}

func keyFuncTail(ctx context.Context, this *Buffer) Result { // Ctrl-E
	allength := this.GetWidthBetween(this.ViewStart, len(this.Buffer))
	if allength < this.ViewWidth() {
		this.puts(this.Buffer[this.Cursor:])
		this.Cursor = len(this.Buffer)
	} else {
		this.backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
		this.ViewStart = len(this.Buffer) - 1
		w := this.Buffer[this.ViewStart].Width()
		for {
			if this.ViewStart <= 0 {
				break
			}
			w_ := w + this.Buffer[this.ViewStart-1].Width()
			if w_ >= this.ViewWidth() {
				break
			}
			w = w_
			this.ViewStart--
		}
		this.puts(this.Buffer[this.ViewStart:])
		this.Cursor = len(this.Buffer)
	}
	return CONTINUE
}

func keyFuncForward(ctx context.Context, this *Buffer) Result { // Ctrl-F
	if this.Cursor >= len(this.Buffer) {
		return CONTINUE
	}
	w := this.GetWidthBetween(this.ViewStart, this.Cursor+1)
	if w < this.ViewWidth() {
		// No Scroll
		this.Buffer[this.Cursor].PrintTo(this.Out)
	} else {
		// Right Scroll
		this.GotoHead()
		if this.Buffer[this.Cursor].Width() > this.Buffer[this.ViewStart].Width() {
			this.ViewStart++
		}
		this.ViewStart++
		this.puts(this.Buffer[this.ViewStart : this.Cursor+1])
		this.Eraseline()
	}
	this.Cursor++
	return CONTINUE
}

func keyFuncBackSpace(ctx context.Context, this *Buffer) Result { // Backspace
	if this.Cursor > 0 {
		this.Cursor--
		delw := this.Delete(this.Cursor, 1)
		if this.Cursor >= this.ViewStart {
			this.backspace(delw)
		} else {
			this.ViewStart = this.Cursor
		}
		this.Repaint(this.Cursor, delw)
	}
	return CONTINUE
}

func keyFuncDelete(ctx context.Context, this *Buffer) Result { // Del
	delw := this.Delete(this.Cursor, 1)
	this.Repaint(this.Cursor, delw)
	return CONTINUE
}

func keyFuncDeleteOrAbort(ctx context.Context, this *Buffer) Result { // Ctrl-D
	if len(this.Buffer) > 0 {
		return keyFuncDelete(ctx, this)
	} else {
		return ABORT
	}
}

func keyFuncInsertSelf(ctx context.Context, this *Buffer, keys string) Result {
	if len(keys) == 2 && keys[0] == '\x1B' { // for AltGr-shift
		keys = keys[1:]
	}
	if ZeroWidthJoinSequenceOk && isZeroWidthJoinStr(keys) && this.Cursor > 0 {
		var tmp strings.Builder
		this.Buffer[this.Cursor-1].WriteTo(&tmp)
		tmp.WriteString(keys)
		this.pending = tmp.String()
		return keyFuncBackSpace(ctx, this)
	} else if VariationSequenceOk && isVariationSelectorLikeStr(keys) && this.Cursor > 0 {
		baseMoji := this.Buffer[this.Cursor-1]
		keyFuncBackSpace(ctx, this)

		var tmp strings.Builder
		baseMoji.WriteTo(&tmp)
		tmp.WriteString(keys)
		keys = tmp.String()
	} else if len(this.pending) > 0 {
		keys = this.pending + keys
		this.pending = ""
	}

	mojis := this.insertString(this.Cursor, keys)
	len_moji := len(mojis)

	w := this.GetWidthBetween(this.ViewStart, this.Cursor)
	w1 := mojis.Width()
	if w+w1 >= this.ViewWidth() {
		// scroll left
		this.GotoHead()
		this.Cursor += len_moji
		this.ResetViewStart()
		this.DrawFromHead()
	} else {
		this.Repaint(this.Cursor, -w1)
		this.Cursor += len_moji
	}
	return CONTINUE
}

func keyFuncClearAfter(ctx context.Context, this *Buffer) Result {
	clipboard.WriteAll(this.SubString(this.Cursor, len(this.Buffer)))

	this.Eraseline()
	u := &undoT{
		pos:  this.Cursor,
		text: moji2string(this.Buffer[this.Cursor:]),
	}
	this.undoes = append(this.undoes, u)
	this.Buffer = this.Buffer[:this.Cursor]
	return CONTINUE
}

func keyFuncClear(ctx context.Context, this *Buffer) Result {
	u := &undoT{
		pos:  0,
		text: moji2string(this.Buffer),
	}
	this.undoes = append(this.undoes, u)
	this.GotoHead()
	this.Eraseline()
	this.Buffer = this.Buffer[:0]
	this.Cursor = 0
	this.ViewStart = 0
	return CONTINUE
}

func keyFuncWordRubout(ctx context.Context, this *Buffer) Result {
	orgCursorPos := this.Cursor
	for this.Cursor > 0 && isSpaceMoji(this.Buffer[this.Cursor-1]) {
		this.Cursor--
	}
	newCursorPos := this.CurrentWordTop()
	clipboard.WriteAll(this.SubString(newCursorPos, orgCursorPos))
	keta := this.Delete(newCursorPos, orgCursorPos-newCursorPos)
	this.Cursor = newCursorPos
	if newCursorPos-this.ViewStart >= 2 {
		this.backspace(keta)
		this.Repaint(newCursorPos, keta)
	} else {
		this.backspace(this.GetWidthBetween(this.ViewStart, orgCursorPos))
		this.RepaintAfterPrompt()
	}
	return CONTINUE
}

func keyFuncClearBefore(ctx context.Context, this *Buffer) Result {
	this.GotoHead()
	clipboard.WriteAll(this.SubString(0, this.Cursor))
	this.Delete(0, this.Cursor)
	this.Cursor = 0
	this.ViewStart = 0
	this.DrawFromHead()
	return CONTINUE
}

func keyFuncCLS(ctx context.Context, this *Buffer) Result {
	io.WriteString(this.Out, "\x1B[1;1H\x1B[2J")
	this.RepaintAll()
	return CONTINUE
}

func keyFuncRepaintOnNewline(ctx context.Context, this *Buffer) Result {
	this.Out.WriteByte('\n')
	this.RepaintAll()
	return CONTINUE
}

func keyFuncQuotedInsert(ctx context.Context, this *Buffer) Result {
	io.WriteString(this.Out, ansiCursorOn)
	defer io.WriteString(this.Out, ansiCursorOff)

	this.Out.Flush()
	if key, err := this.GetKey(); err == nil {
		return keyFuncInsertSelf(ctx, this, key)
	} else {
		return CONTINUE
	}
}

func keyFuncPaste(ctx context.Context, this *Buffer) Result {
	text, err := clipboard.ReadAll()
	if err != nil {
		return CONTINUE
	}
	text = strings.TrimRight(text, "\r\n\000")
	this.InsertAndRepaint(text)
	return CONTINUE
}

func keyFuncPasteQuote(ctx context.Context, this *Buffer) Result {
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
}

func maxInt(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func keyFuncSwapChar(ctx context.Context, this *Buffer) Result {
	if len(this.Buffer) == this.Cursor {
		if this.Cursor < 2 {
			return CONTINUE
		}
		u := &undoT{
			pos:  this.Cursor,
			del:  2,
			text: moji2string(this.Buffer[this.Cursor-2 : this.Cursor]),
		}
		this.undoes = append(this.undoes, u)
		this.Buffer[this.Cursor-2], this.Buffer[this.Cursor-1] = this.Buffer[this.Cursor-1], this.Buffer[this.Cursor-2]

		redrawStart := maxInt(this.Cursor-2, this.ViewStart)
		this.backspace(this.GetWidthBetween(redrawStart, this.Cursor))
		this.puts(this.Buffer[redrawStart:this.Cursor])
	} else {
		if this.Cursor < 1 {
			return CONTINUE
		}
		u := &undoT{
			pos:  this.Cursor - 1,
			del:  2,
			text: moji2string(this.Buffer[this.Cursor-1 : this.Cursor+1]),
		}
		this.undoes = append(this.undoes, u)

		w := this.GetWidthBetween(this.ViewStart, this.Cursor+1)
		this.Buffer[this.Cursor-1], this.Buffer[this.Cursor] = this.Buffer[this.Cursor], this.Buffer[this.Cursor-1]
		if w >= this.ViewWidth() {
			// cursor move right and scroll
			w_1 := w - this.Buffer[this.Cursor].Width()
			this.backspace(w_1)
			this.ViewStart++
			this.puts(this.Buffer[this.ViewStart : this.Cursor+1])
		} else {
			// no necessary to scroll
			redrawStart := maxInt(this.Cursor-1, this.ViewStart)
			this.backspace(this.GetWidthBetween(redrawStart, this.Cursor))
			this.puts(this.Buffer[redrawStart : this.Cursor+1])
		}
		this.Cursor++
	}
	return CONTINUE
}

func keyFuncBackwardWord(ctx context.Context, this *Buffer) Result {
	newPos := this.Cursor
	for newPos > 0 && isSpaceMoji(this.Buffer[newPos-1]) {
		newPos--
	}
	for newPos > 0 && !isSpaceMoji(this.Buffer[newPos-1]) {
		newPos--
	}
	if newPos >= this.ViewStart {
		w := this.GetWidthBetween(newPos, this.Cursor)
		this.backspace(w)
		this.Cursor = newPos
	} else {
		this.GotoHead()
		this.Cursor = newPos
		this.ViewStart = newPos
		this.DrawFromHead()
	}
	return CONTINUE
}

func keyFuncForwardWord(ctx context.Context, this *Buffer) Result {
	newPos := this.Cursor
	for newPos < len(this.Buffer) && !isSpaceMoji(this.Buffer[newPos]) {
		newPos++
	}
	for newPos < len(this.Buffer) && isSpaceMoji(this.Buffer[newPos]) {
		newPos++
	}
	w := this.GetWidthBetween(this.ViewStart, newPos)
	if w < this.ViewWidth() {
		this.puts(this.Buffer[this.Cursor:newPos])
		this.Cursor = newPos
	} else {
		this.GotoHead()
		this.Cursor = newPos
		this.ResetViewStart()
		this.DrawFromHead()
	}
	return CONTINUE
}

func keyFuncUndo(ctx context.Context, this *Buffer) Result {
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
		t := string2moji(u.text)
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
	this.DrawFromHead()
	return CONTINUE
}
