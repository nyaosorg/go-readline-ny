package readline

import (
	"context"
)

// IHistory is the interface ReadLine can use as container for history.
// It can be set to Editor.History field
type IHistory interface {
	Len() int
	At(int) string
}

type _EmptyHistory struct{}

// Len always returns zero because the receiver is dummy.
func (_EmptyHistory) Len() int { return 0 }

// At always returns empty-string because the receiver is dummy.
func (_EmptyHistory) At(int) string { return "" }

// CmdPreviousHistory is the command that replaces the line to the previous entry in the history (Ctrl-P)
var CmdPreviousHistory = NewGoCommand("PREVIOUS_HISTORY", cmdPreviousHistory)

func (B *Buffer) saveModifiedHistory() {
	s := B.String()
	if B.historyPointer < B.History.Len() && B.History.At(B.historyPointer) == s {
		return
	}
	if B.modifiedHistory == nil {
		B.modifiedHistory = make(map[int]string)
	}
	B.modifiedHistory[B.historyPointer] = s
}

func (B *Buffer) getHistory() string {
	s, ok := B.modifiedHistory[B.historyPointer]
	if ok {
		return s
	}
	if B.historyPointer < 0 || B.historyPointer >= B.History.Len() {
		return ""
	}
	return B.History.At(B.historyPointer)
}

func cmdPreviousHistory(ctx context.Context, this *Buffer) Result {
	if this.History.Len() <= 0 {
		return CONTINUE
	}
	if this.historyPointer <= 0 {
		if !this.HistoryCycling {
			return CONTINUE
		}
		this.historyPointer = this.History.Len() + 1
	}
	this.saveModifiedHistory()
	this.historyPointer--
	CmdKillWholeLine.Func(ctx, this)
	if s := this.getHistory(); s != "" {
		this.InsertString(0, s)
		this.ViewStart = 0
		this.Cursor = 0
		CmdEndOfLine.Func(ctx, this)
	}
	return CONTINUE
}

// CmdNextHistory is the command that replaces the line to the next entry of the history. (Ctrl-N)
var CmdNextHistory = NewGoCommand("NEXT_HISTORY", cmdNextHistory)

func cmdNextHistory(ctx context.Context, this *Buffer) Result {
	if this.History.Len() <= 0 {
		return CONTINUE
	}
	if this.historyPointer+1 > this.History.Len() {
		return CONTINUE
	}
	this.saveModifiedHistory()
	this.historyPointer++
	CmdKillWholeLine.Func(ctx, this)
	if s := this.getHistory(); s != "" {
		this.InsertString(0, s)
		this.ViewStart = 0
		this.Cursor = 0
		CmdEndOfLine.Func(ctx, this)
	}
	return CONTINUE
}
